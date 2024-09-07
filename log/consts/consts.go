package consts

import "context"

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
	LogLevelPanic
)

type Fields map[string]string

type Logger interface {
	Log(ctx context.Context, level LogLevel, format string, args ...interface{})
	WithFields(Fields) Logger
}
