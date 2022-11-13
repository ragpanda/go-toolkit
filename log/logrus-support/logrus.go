package logrus_support

import (
	"context"
	"io"
	"os"

	"github.com/ragpanda/go-toolkit/log/consts"
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	root  *logrus.Logger
	entry *logrus.Entry
}

type LogrusConfig struct {
	Level        logrus.Level
	Output       io.Writer
	DisableColor bool
	SkipPkg      []string
}

func NewLogrusLogger(ctx context.Context, config *LogrusConfig) *LogrusLogger {
	l := logrus.New()
	logger := &LogrusLogger{
		root:  l,
		entry: l.WithContext(ctx),
	}

	if config == nil {
		config = &LogrusConfig{}
	}

	if config.Level == 0 {
		config.Level = logrus.DebugLevel
	}

	logger.root.SetLevel(config.Level)

	if config.Output != nil {
		logger.root.SetOutput(config.Output)
	} else {
		logger.root.SetOutput(os.Stdout)
	}

	logger.root.SetFormatter(&ConsoleFormatter{
		DisableTimestamp: false,
		EnableColors:     !config.DisableColor,
		EnableEntryOrder: true,
	})

	logger.root.AddHook(NewCallerHook(config.Level, config.SkipPkg))
	return logger
}

func (self *LogrusLogger) Log(ctx context.Context, level consts.LogLevel, format string, args ...interface{}) {
	self.entry.WithContext(ctx).Logf(levelMapping[level], format, args...)
}

func (self *LogrusLogger) WithFields(f consts.Fields) consts.Logger {
	newLogger := LogrusLogger{
		root:  self.root,
		entry: self.entry.WithFields((logrus.Fields)(f)),
	}
	return &newLogger
}

var levelMapping = map[consts.LogLevel]logrus.Level{
	consts.LogLevelDebug: logrus.DebugLevel,
	consts.LogLevelInfo:  logrus.InfoLevel,
	consts.LogLevelWarn:  logrus.WarnLevel,
	consts.LogLevelError: logrus.ErrorLevel,
	consts.LogLevelFatal: logrus.FatalLevel,
	consts.LogLevelPanic: logrus.PanicLevel,
}
