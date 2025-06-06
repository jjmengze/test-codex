package logger

import "context"

var (
	Debug = LogLevel{level: 0, name: "debug"}
	Info  = LogLevel{level: 1, name: "info"}
	Warn  = LogLevel{level: 2, name: "warn"}
	Error = LogLevel{level: 3, name: "error"}
	Fatal = LogLevel{level: 4, name: "fatal"}
	Force = LogLevel{level: 5, name: "force"}
)

type LogLevel struct {
	level int
	name  string
}

type InfLogKeyValue interface {
	Key() any
	Value() any
}

type ClosureFunc func() []InfLogKeyValue

// Logger -
type Logger interface {
	// DebugF Print format like fmt.Printf of Debug Level
	DebugF(msg string, args ...interface{})
	// DebugW Print wrap of Debug Level
	// ex: DebugW("test","uuid",1234)
	DebugW(msg string, args ...interface{})
	// DebugM Print wrap of Debug Level
	// ex: DebugM("test",map[platform_interface{}]platform_interface{}{"uuid":1234})
	DebugM(msg string, valueMap map[interface{}]interface{})
	// DebugK Print msg & data of Debug Level
	//
	// Deprecated: should not use this function, it will be removed in the future.
	// Recommended to use DebugM instead.
	// ex: DebugK("test",struct {
	//				UUID string `json:"uuid"`
	//			}{
	//				UUID: "1234",
	//			})
	// output: test   {Logger.GetKFormatKey(): {"uuid":"1234"}}
	DebugK(msg string, value interface{})

	// InfoF Print format like fmt.Printf of Info Level
	InfoF(msg string, args ...interface{})
	// InfoW Print wrap of Info Level
	//
	// ex: InfoW("test","uuid",1234)
	InfoW(msg string, args ...interface{})
	// InfoM Print wrap of Info Level
	//
	// ex: InfoM("test",map[any]any{"uuid":1234})
	InfoM(msg string, valueMap map[interface{}]interface{})
	// InfoK Print msg & data of Info Level
	//
	// Deprecated: should not use this function, it will be removed in the future.
	// Recommended to use InfoM instead.
	// ex: InfoK("test",struct {
	//				UUID string `json:"uuid"`
	//			}{
	//				UUID: "1234",
	//			})
	// output: test   {kFormatKey: {"uuid":"1234"}}
	InfoK(msg string, value interface{})

	// WarnF Print format like fmt.Printf of Warn Level
	WarnF(msg string, args ...interface{})
	// WarnW Print wrap of Warn Level
	// ex: WarnW("test","uuid",1234)
	WarnW(msg string, args ...interface{})
	// WarnM Print wrap of Warn Level
	// ex: WarnM("test",map[platform_interface{}]platform_interface{}{"uuid":1234})
	WarnM(msg string, valueMap map[interface{}]interface{})
	// WarnK Print msg & data of Warn Level
	//
	// Deprecated: should not use this function, it will be removed in the future.
	// Recommended to use WarnM instead.
	// ex: WarnK("test",struct {
	//				UUID string `json:"uuid"`
	//			}{
	//				UUID: "1234",
	//			})
	// output: test   {Logger.GetKFormatKey(): {"uuid":"1234"}}
	WarnK(msg string, value interface{})

	// ErrorF Print format like fmt.Printf of Error Level
	ErrorF(msg string, args ...interface{})
	// ErrorW Print wrap of Error Level
	// ex: ErrorW("test","uuid",1234)
	ErrorW(msg string, args ...interface{})
	// ErrorM Print wrap of Error Level
	// ex: ErrorM("test",map[platform_interface{}]platform_interface{}{"uuid":1234})
	ErrorM(msg string, valueMap map[interface{}]interface{})
	// ErrorK Print msg & data of Error Level
	//
	// Deprecated: should not use this function, it will be removed in the future.
	// Recommended to use ErrorM instead.
	// ex: ErrorK("test",struct {
	//				UUID string `json:"uuid"`
	//			}{
	//				UUID: "1234",
	//			})
	// output: test   {Logger.GetKFormatKey(): {"uuid":"1234"}}
	ErrorK(msg string, value interface{})

	// FatalF Print format like fmt.Printf of Fatal Level
	FatalF(msg string, args ...interface{})
	// FatalW Print wrap of Fatal Level
	// ex: FatalW("test","uuid",1234)
	FatalW(msg string, args ...interface{})
	// FatalM Print wrap of Fatal Level
	// ex: FatalM("test",map[platform_interface{}]platform_interface{}{"uuid":1234})
	FatalM(msg string, valueMap map[interface{}]interface{})
	// FatalK Print msg & data of Fatal Level
	//
	// Deprecated: should not use this function, it will be removed in the future.
	// Recommended to use FatalM instead.
	// ex: FatalK("test",struct {
	//				UUID string `json:"uuid"`
	//			}{
	//				UUID: "1234",
	//			})
	// output: test   {Logger.GetKFormatKey(): {"uuid":"1234"}}
	FatalK(msg string, value interface{})

	Panic(format string, args ...interface{})
	// TODO
	// WithError, WithContext

	Output(callDepth int, s string) error

	// WithAnyMap add map to logger
	// ex: WithAnyMap(map[any]any{"uuid":1234})
	WithAnyMap(valueMap map[any]any) Logger

	// With add key-value to logger
	// ex:
	//
	// With(
	//		Attr("uuid",1234),
	//		Attr("content",struct {
	//			UUID string `json:"uuid"`
	//		}{
	//			UUID: "1234",
	//		}),
	// )
	With(...InfLogKeyValue) Logger

	// WithContext add context to logger
	WithContext(ctx context.Context) Logger

	// WithLogger add logger to context
	WithLogger(ctx context.Context, log Logger) context.Context

	// WithClosures add closure function to logger
	// When you want to add some dynamic fields to log, you can use this function.
	// Experiment feature
	// Warn: plz use with caution, this feature is not stable.
	// It may be removed in the future.
	WithClosures(closures ...ClosureFunc) Logger
}
