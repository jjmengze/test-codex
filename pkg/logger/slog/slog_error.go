package slog

import (
	"context"
	"log/slog"
)

func (l *slogLogger) ErrorF(msg string, args ...interface{}) {
	l.outputF(context.Background(), slog.LevelError, msg, args...)
}

func (l *slogLogger) ErrorW(msg string, args ...interface{}) {
	l.outputW(context.Background(), slog.LevelError, msg, args...)
}

func (l *slogLogger) ErrorM(msg string, valueMap map[interface{}]interface{}) {
	l.outputM(context.Background(), slog.LevelError, msg, valueMap)
}

func (l *slogLogger) ErrorK(msg string, value interface{}) {
	l.outputK(context.Background(), slog.LevelError, msg, value)
}
