package slog

import (
	"context"
	"log/slog"
)

func (l *slogLogger) InfoF(msg string, args ...interface{}) {
	l.outputF(context.Background(), slog.LevelInfo, msg, args...)
}

func (l *slogLogger) InfoW(msg string, args ...interface{}) {
	l.outputW(context.Background(), slog.LevelInfo, msg, args...)
}

func (l *slogLogger) InfoM(msg string, valueMap map[interface{}]interface{}) {
	l.outputM(context.Background(), slog.LevelInfo, msg, valueMap)
}

func (l *slogLogger) InfoK(msg string, value interface{}) {
	l.outputK(context.Background(), slog.LevelInfo, msg, value)
}
