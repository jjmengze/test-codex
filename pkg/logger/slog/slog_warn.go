package slog

import (
	"context"
	"log/slog"
)

func (l *slogLogger) WarnF(msg string, args ...interface{}) {
	l.outputF(context.Background(), slog.LevelWarn, msg, args...)
}

func (l *slogLogger) WarnW(msg string, args ...interface{}) {
	l.outputW(context.Background(), slog.LevelWarn, msg, args...)
}

func (l *slogLogger) WarnM(msg string, valueMap map[interface{}]interface{}) {
	l.outputM(context.Background(), slog.LevelWarn, msg, valueMap)
}

func (l *slogLogger) WarnK(msg string, value interface{}) {
	l.outputK(context.Background(), slog.LevelWarn, msg, value)
}
