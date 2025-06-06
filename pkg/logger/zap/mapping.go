package handler

import (
	"go.uber.org/zap/zapcore"
	"log-receiver/pkg/logger"
)

var ikitToZapLevelMap = map[logger.LogLevel]zapcore.Level{
	logger.Debug: zapcore.DebugLevel,
	logger.Info:  zapcore.InfoLevel,
	logger.Warn:  zapcore.WarnLevel,
	logger.Error: zapcore.ErrorLevel,
	logger.Fatal: zapcore.FatalLevel,
}
