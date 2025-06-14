// Copyright (c) 2023 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

//go:build go1.21

package handler

import (
	"context"
	"log/slog"
	"runtime"

	"log-receiver/pkg/logger/zap/internal/stacktrace"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// Skipping 4:
	// slog_logger/slog_xxx slog_logger.(*slogLogger).XXX
	// zap/handler handler/zap.(*handler).Handle
	// slog/logger log/slog.(*Logger).log
	// slog/logger log/slog.(*Logger).<level>
	defaultCallerSkip = 5
)

var _ slog.Handler = (*handler)(nil)

// handler implements the slog.handler by writing to a zap Core.
type handler struct {
	core       zapcore.Core
	addCaller  bool
	addStackAt slog.Level

	// List of unapplied groups.
	//
	// These are applied only if we encounter a real field
	// to avoid creating empty namespaces -- which is disallowed by slog's
	// usage contract.
	groups []string
}

// newHandler builds a [handler] that writes to the supplied [zapcore.Core]
// with options.
func newHandler(core zapcore.Core, opts ...handlerOption) *handler {
	h := &handler{
		core:       core,
		addStackAt: slog.LevelError,
	}
	for _, v := range opts {
		v.apply(h)
	}
	return h
}

// groupObject holds all the Attrs saved in a slog.GroupValue.
type groupObject []slog.Attr

func (gs groupObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, attr := range gs {
		convertAttrToField(attr).AddTo(enc)
	}
	return nil
}

func convertAttrToField(attr slog.Attr) zapcore.Field {
	if attr.Equal(slog.Attr{}) {
		// Ignore empty attrs.
		return zap.Skip()
	}

	switch attr.Value.Kind() {
	case slog.KindBool:
		return zap.Bool(attr.Key, attr.Value.Bool())
	case slog.KindDuration:
		return zap.Duration(attr.Key, attr.Value.Duration())
	case slog.KindFloat64:
		return zap.Float64(attr.Key, attr.Value.Float64())
	case slog.KindInt64:
		return zap.Int64(attr.Key, attr.Value.Int64())
	case slog.KindString:
		return zap.String(attr.Key, attr.Value.String())
	case slog.KindTime:
		return zap.Time(attr.Key, attr.Value.Time())
	case slog.KindUint64:
		return zap.Uint64(attr.Key, attr.Value.Uint64())
	case slog.KindGroup:
		if attr.Key == "" {
			// Inlines recursively.
			return zap.Inline(groupObject(attr.Value.Group()))
		}
		return zap.Object(attr.Key, groupObject(attr.Value.Group()))
	case slog.KindLogValuer:
		return convertAttrToField(slog.Attr{
			Key: attr.Key,
			// TODO: resolve the value in a lazy way.
			// This probably needs a new Zap field type
			// that can be resolved lazily.
			Value: attr.Value.Resolve(),
		})
	default:
		return zap.Any(attr.Key, attr.Value.Any())
	}
}

// convertSlogLevel maps slog Levels to zap Levels.
// Note that there is some room between slog levels while zap levels are continuous, so we can't 1:1 map them.
// See also https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md?pli=1#levels
func convertSlogLevel(l slog.Level) zapcore.Level {
	switch {
	case l >= LevelPanic:
		return zapcore.PanicLevel
	case l >= LevelFatal:
		return zapcore.FatalLevel
	case l >= LevelError:
		return zapcore.ErrorLevel
	case l >= LevelWarn:
		return zapcore.WarnLevel
	case l >= LevelInfo:
		return zapcore.InfoLevel
	default:
		return zapcore.DebugLevel
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.core.Enabled(convertSlogLevel(level))
}

// Handle handles the Record.
func (h *handler) Handle(ctx context.Context, record slog.Record) error {
	ent := zapcore.Entry{
		Level:   convertSlogLevel(record.Level),
		Time:    record.Time,
		Message: record.Message,
	}
	ce := h.core.Check(ent, nil)
	if ce == nil {
		return nil
	}

	if h.addCaller && record.PC != 0 {
		pc, file, line, ok := runtime.Caller(defaultCallerSkip)
		funcDetail := runtime.FuncForPC(pc)
		if ok && funcDetail != nil {
			ce.Caller = zapcore.EntryCaller{
				Defined:  true,
				PC:       pc,
				File:     file,
				Line:     line,
				Function: funcDetail.Name(),
			}
		}
	}

	if record.Level >= h.addStackAt {
		ce.Stack = stacktrace.Take(defaultCallerSkip)
	}

	fields := make([]zapcore.Field, 0, record.NumAttrs()+len(h.groups))

	var addedNamespace bool
	record.Attrs(func(attr slog.Attr) bool {
		f := convertAttrToField(attr)
		if !addedNamespace && len(h.groups) > 0 && f != zap.Skip() {
			// Namespaces are added only if at least one field is present
			// to avoid creating empty groups.
			fields = h.appendGroups(fields)
			addedNamespace = true
		}
		fields = append(fields, f)
		return true
	})

	ce.Write(fields...)
	return nil
}

func (h *handler) appendGroups(fields []zapcore.Field) []zapcore.Field {
	for _, g := range h.groups {
		fields = append(fields, zap.Namespace(g))
	}
	return fields
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver.go's attributes and the arguments.
func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	fields := make([]zapcore.Field, 0, len(attrs)+len(h.groups))
	var addedNamespace bool
	for _, attr := range attrs {
		f := convertAttrToField(attr)
		if !addedNamespace && len(h.groups) > 0 && f != zap.Skip() {
			// Namespaces are added only if at least one field is present
			// to avoid creating empty groups.
			fields = h.appendGroups(fields)
			addedNamespace = true
		}
		fields = append(fields, f)
	}

	cloned := *h
	cloned.core = h.core.With(fields)
	if addedNamespace {
		// These groups have been applied so we can clear them.
		cloned.groups = nil
	}
	return &cloned
}

// WithGroup returns a new Handler with the given group appended to
// the receiver.go's existing groups.
func (h *handler) WithGroup(group string) slog.Handler {
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = group

	cloned := *h
	cloned.groups = newGroups
	return &cloned
}
