package logrus_support

import (
	"context"

	"github.com/ragpanda/go-toolkit/log"
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	root  *logrus.Logger
	entry *logrus.Entry
}

func NewLogrusLogger(ctx context.Context, l *logrus.Logger) *LogrusLogger {
	if l == nil {
		l = logrus.StandardLogger()
	}

	logger := &LogrusLogger{
		root:  l,
		entry: l.WithContext(ctx),
	}
	return logger
}

func (self *LogrusLogger) Log(ctx context.Context, level log.LogLevel, format string, args ...interface{}) {
	self.entry.WithContext(ctx).Logf(levelMapping[level], format, args...)
}

func (self *LogrusLogger) WithFields(f log.Fields) log.Logger {
	newLogger := LogrusLogger{
		root:  self.root,
		entry: self.entry.WithFields((logrus.Fields)(f)),
	}
	return &newLogger
}

var levelMapping = map[log.LogLevel]logrus.Level{
	log.LogLevelDebug: logrus.DebugLevel,
	log.LogLevelInfo:  logrus.InfoLevel,
	log.LogLevelWarn:  logrus.WarnLevel,
	log.LogLevelError: logrus.ErrorLevel,
	log.LogLevelFatal: logrus.FatalLevel,
	log.LogLevelPanic: logrus.PanicLevel,
}



