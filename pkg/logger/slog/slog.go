package slog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"go.opentelemetry.io/otel/trace"
	"log-receiver/pkg/logger"
	zap "log-receiver/pkg/logger/zap"
)

const (
	LevelDebug slog.Level = slog.LevelDebug
	LevelInfo  slog.Level = slog.LevelInfo
	LevelWarn  slog.Level = slog.LevelWarn
	LevelError slog.Level = slog.LevelError
	LevelFatal slog.Level = 12
	LevelPanic slog.Level = 16
)

// slogLogger is a logger implementation that wraps slog.Logger.
// Note: It's no need to write this code, it's doubled effort.
// It's the same as NewSlogLogger func return value.
var _ logger.Logger = (*slogLogger)(nil)

var (
	defaultServiceKey = "service"
	defaultKFormatKey = "content"

	globalLoggerMutex sync.RWMutex
	globalLogger      = NewSlogLogger(Config{SlogHandler: zap.NewZapHandler(os.Stderr, &zap.HandlerOptions{
		AddSource: true,
		Level:     logger.Info,
	})})
)

type loggerContextKey struct{}

type Config struct {
	SlogHandler slog.Handler
	ServiceName string
	ServiceKey  string
	KFormatKey  string
}

type slogLogger struct {
	innerLogger   *slog.Logger
	afterFuncList []logger.ClosureFunc

	enabledKJsonContent bool
	kFormatKey          string
}

func (l *slogLogger) WithClosures(closures ...logger.ClosureFunc) logger.Logger {
	c := *l
	c.afterFuncList = append(c.afterFuncList, closures...)
	return &c
}

func (l *slogLogger) With(value ...logger.InfLogKeyValue) logger.Logger {
	c := *l
	var args []any
	for _, kv := range value {
		args = append(args, kv.Key(), kv.Value())
	}
	c.innerLogger = c.innerLogger.With(args...)
	return &c
}

func (l *slogLogger) WithAnyMap(valueMap map[any]any) logger.Logger {
	c := *l
	args := l.mFormat(valueMap)
	c.innerLogger = c.innerLogger.With(args...)
	return &c
}

// NewSlogLogger creates a new slog logger.
func NewSlogLogger(config Config) logger.Logger {
	if config.SlogHandler == nil {
		config.SlogHandler = zap.NewZapHandler(os.Stderr, &zap.HandlerOptions{
			AddSource: true,
			Level:     logger.Info,
		})
	}
	if config.ServiceKey == "" {
		config.ServiceKey = defaultServiceKey
	}
	if config.ServiceName == "" && len(os.Args) >= 2 {
		config.ServiceName = os.Args[1]
	}

	innerLogger := slog.New(config.SlogHandler)

	if config.ServiceKey != "" && config.ServiceName != "" {
		innerLogger = slog.New(config.SlogHandler).With(config.ServiceKey, config.ServiceName)
	}

	var enabledKJsonContent bool
	kFormatKey := defaultKFormatKey
	if config.KFormatKey != "" {
		enabledKJsonContent = true
		kFormatKey = config.KFormatKey
	}

	return &slogLogger{
		innerLogger:         innerLogger,
		enabledKJsonContent: enabledKJsonContent,
		kFormatKey:          kFormatKey,
	}
}

func SetGlobalLogger(logger logger.Logger) {
	globalLoggerMutex.Lock()
	defer globalLoggerMutex.Unlock()

	globalLogger = logger
}

func GetGlobalLogger() logger.Logger {
	globalLoggerMutex.RLock()
	defer globalLoggerMutex.RUnlock()

	return globalLogger
}

func (l *slogLogger) mFormat(valueMap map[any]any) []any {
	var args []any

	for k, v := range valueMap {
		args = append(args, k, v)
	}

	return args
}

func (l *slogLogger) kFormat(value any) []any {
	var args []any

	args = append(args, l.kFormatKey, value)

	return args
}

func (l *slogLogger) Output(callDepth int, s string) error {
	l.InfoW(s)
	return nil
}

// Panic since slog does not provide a 'Panic' method.
// See also https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md?pli=1#levels
func (l *slogLogger) Panic(msg string, args ...any) {
	msg = fmt.Sprintf(msg, args...)
	l.innerLogger.Log(context.Background(), LevelPanic, msg)
	panic(msg)
}

func (l *slogLogger) WithContext(ctx context.Context) logger.Logger {
	// if context has logger, return it immediately
	ctxLogger := ctx.Value(loggerContextKey{})
	if ctxLogger != nil {
		return ctxLogger.(logger.Logger)
	}

	span := trace.SpanContextFromContext(ctx)
	if span.IsValid() {
		// for concurrency safety
		c := *l
		c.innerLogger = c.innerLogger.With(
			"trace-id", span.TraceID().String(),
			"span-id", span.SpanID().String(),
		)

		return &c
	}

	return l
}

func (l *slogLogger) WithLogger(ctx context.Context, logger logger.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

func (l *slogLogger) after() []any {
	var afterAttrs []any

	for _, f := range l.afterFuncList {
		data := f()
		for _, d := range data {
			afterAttrs = append(afterAttrs, d.Key(), d.Value())
		}
	}

	return afterAttrs
}
