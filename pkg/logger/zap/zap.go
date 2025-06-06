package handler

import (
	"io"
	"log/slog"

	"log-receiver/pkg/logger"

	"go.uber.org/zap/zapcore"
)

const (
	LevelDebug slog.Level = slog.LevelDebug
	LevelInfo  slog.Level = slog.LevelInfo
	LevelWarn  slog.Level = slog.LevelWarn
	LevelError slog.Level = slog.LevelError
	LevelFatal slog.Level = 12
	LevelPanic slog.Level = 16
)

type HandlerOptions struct {
	AddSource       bool
	AddStacktraceAt slog.Level
	Level           logger.LogLevel
	EncoderProvider func(cfg zapcore.EncoderConfig) zapcore.Encoder
	EncoderConfig   *zapcore.EncoderConfig
}

func NewZapHandler(w io.Writer, opts *HandlerOptions) slog.Handler {
	if opts == nil {
		opts = &HandlerOptions{}
	}
	if opts.EncoderProvider == nil {
		opts.EncoderProvider = zapcore.NewJSONEncoder
	}
	if opts.EncoderConfig == nil {
		opts.EncoderConfig = NewRedisEncoderConfig()
	}
	level := ikitToZapLevelMap[opts.Level]
	addStacktraceAt := LevelError
	if opts.AddStacktraceAt != 0 {
		addStacktraceAt = opts.AddStacktraceAt
	}

	encoder := opts.EncoderProvider(*opts.EncoderConfig)
	output := zapcore.Lock(zapcore.AddSync(w))
	core := zapcore.NewCore(encoder, output, level)

	return newHandler(core, withCaller(opts.AddSource), withStacktrace(addStacktraceAt))
}

func NewEncoderConfig() *zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "T",
		NameKey:        "name",
		CallerKey:      "caller",
		FunctionKey:    "fun",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder

	return &encoderConfig
}

func NewRedisEncoderConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
