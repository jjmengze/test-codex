package slog

import (
	"context"
)

// FatalF since slog does not provide a 'fatal' level.
// See also https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md?pli=1#levels
func (l *slogLogger) FatalF(msg string, args ...interface{}) {
	l.outputF(context.Background(), LevelFatal, msg, args...)
	panic(msg)
}

// FatalW since slog does not provide a 'fatal' level.
// See also https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md?pli=1#levels
func (l *slogLogger) FatalW(msg string, args ...interface{}) {
	l.outputW(context.Background(), LevelFatal, msg, args...)
	panic(msg)
}

// FatalM since slog does not provide a 'fatal' level.
// See also https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md?pli=1#levels
func (l *slogLogger) FatalM(msg string, valueMap map[interface{}]interface{}) {
	l.outputM(context.Background(), LevelFatal, msg, valueMap)
	panic(msg)
}

// FatalK since slog does not provide a 'fatal' level.
// See also https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md?pli=1#levels
func (l *slogLogger) FatalK(msg string, value interface{}) {
	l.outputK(context.Background(), LevelFatal, msg, value)
	panic(msg)
}
