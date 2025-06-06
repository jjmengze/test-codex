package slog

import (
	"context"
	"log/slog"
)

func (l *slogLogger) DebugF(msg string, args ...any) {
	l.outputF(context.Background(), slog.LevelDebug, msg, args...)
}

func (l *slogLogger) DebugW(msg string, args ...any) {
	l.outputW(context.Background(), slog.LevelDebug, msg, args...)
}

func (l *slogLogger) DebugM(msg string, valueMap map[any]any) {
	l.outputM(context.Background(), slog.LevelDebug, msg, valueMap)
}

func (l *slogLogger) DebugK(msg string, value any) {
	l.outputK(context.Background(), slog.LevelDebug, msg, value)
}
