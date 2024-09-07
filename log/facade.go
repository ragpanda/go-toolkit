package log

import (
	"context"

	"github.com/ragpanda/go-toolkit/log/consts"
	logrus_support "github.com/ragpanda/go-toolkit/log/logrus-support"
)

func Debug(ctx context.Context, format string, args ...interface{}) {
	Log(ctx, consts.LogLevelDebug, format, args...)
}

func Info(ctx context.Context, format string, args ...interface{}) {
	Log(ctx, consts.LogLevelInfo, format, args...)
}
func Warn(ctx context.Context, format string, args ...interface{}) {
	Log(ctx, consts.LogLevelWarn, format, args...)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	Log(ctx, consts.LogLevelError, format, args...)
}

func Fatal(ctx context.Context, format string, args ...interface{}) {
	Log(ctx, consts.LogLevelFatal, format, args...)
}

func Panic(ctx context.Context, format string, args ...interface{}) {
	Log(ctx, consts.LogLevelPanic, format, args...)
}

func WithField(ctx context.Context, kv consts.Fields) context.Context {
	l := globalLogger.WithFields(kv)
	return context.WithValue(ctx, LoggerCtxKey, l)
}

func Log(ctx context.Context, level consts.LogLevel, format string, args ...interface{}) {
	if ctx == nil {
		ctx = context.Background()
	}

	value := ctx.Value(LoggerCtxKey)
	if value != nil {
		logger, ok := value.(consts.Logger)
		if ok {
			logger.Log(ctx, level, format, args...)
			return
		}

	}

	globalLogger.Log(ctx, level, format, args...)
}

const LoggerCtxKey = "__logger"

var globalLogger consts.Logger

func SetGlobal(l consts.Logger) {
	globalLogger = l
}

func GetGlobal() consts.Logger {
	return globalLogger
}

func init() {
	l := logrus_support.NewLogrusLogger(context.Background(), nil)
	SetGlobal(l)
}
