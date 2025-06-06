package slog

import (
	"context"
	"fmt"
	"log/slog"
)

func (l *slogLogger) outputF(ctx context.Context, level slog.Level, msg string, args ...interface{}) {
	if !l.innerLogger.Enabled(ctx, level) {
		return
	}

	var afterAttrs []any

	if len(l.afterFuncList) > 0 {
		afterAttrs = append(afterAttrs, l.after()...)
	}

	preMsg := fmt.Sprintf(msg, args...)
	l.innerLogger.Log(ctx, level, preMsg, afterAttrs...)
}

func (l *slogLogger) outputW(ctx context.Context, level slog.Level, msg string, args ...interface{}) {
	if !l.innerLogger.Enabled(ctx, level) {
		return
	}

	if len(l.afterFuncList) > 0 {
		args = append(args, l.after()...)
	}

	l.innerLogger.Log(ctx, level, msg, args...)
}

func (l *slogLogger) outputM(ctx context.Context, level slog.Level, msg string, valueMap map[interface{}]interface{}) {
	if !l.innerLogger.Enabled(ctx, level) {
		return
	}

	args := l.mFormat(valueMap)
	if len(l.afterFuncList) > 0 {
		args = append(args, l.after()...)
	}
	l.innerLogger.Log(ctx, level, msg, args...)
}

func (l *slogLogger) outputK(ctx context.Context, level slog.Level, msg string, value interface{}) {
	if !l.innerLogger.Enabled(ctx, level) {
		return
	}

	args := l.kFormat(value)
	if len(l.afterFuncList) > 0 {
		args = append(args, l.after()...)
	}
	l.innerLogger.Log(ctx, level, msg, args...)
}
