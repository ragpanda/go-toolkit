package log

import "context"

func Debug(ctx context.Context, format string, args ...interface{}) {
	Log(ctx, LogLevelDebug, format, args...)
}

func Info(ctx context.Context, format string, args ...interface{}) {
	Log(ctx, LogLevelInfo, format, args...)
}
func Warn(ctx context.Context, format string, args ...interface{}) {
	Log(ctx, LogLevelWarn, format, args...)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	Log(ctx, LogLevelError, format, args...)
}

func Fatal(ctx context.Context, format string, args ...interface{}) {
	Log(ctx, LogLevelFatal, format, args...)
}

func Panic(ctx context.Context, format string, args ...interface{}) {
	Log(ctx, LogLevelPanic, format, args...)
}

func Log(ctx context.Context, level LogLevel, format string, args ...interface{}) {
	logger, ok := ctx.Value(LoggerCtxKey).(Logger)
	if ok {
		logger.Log(ctx, level, format, args...)
		return
	}

	globalLogger.Log(ctx, level, format, args...)
}

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
	LogLevelPanic
)

type Fields map[string]interface{}

type Logger interface {
	Log(ctx context.Context, level LogLevel, format string, args ...interface{})
	WithFields(Fields) Logger
}

const LoggerCtxKey = "__logger"

var globalLogger Logger

func SetGlobal(l Logger) {
	globalLogger = l
}

func GetGlobal() Logger {
	return globalLogger
}
