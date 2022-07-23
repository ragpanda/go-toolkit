package logrus_support

import (
	"io"
	"os"

	"github.com/ragpanda/go-toolkit/log"
	"github.com/sirupsen/logrus"
)

type LogrusCommonConfig struct {
	Level        logrus.Level
	Output       io.Writer
	DisableColor bool
}

func Init(l *LogrusLogger, config *LogrusCommonConfig) {
	if config == nil {
		config = &LogrusCommonConfig{}
	}
	if config.Level == 0 {
		config.Level = logrus.DebugLevel
	}

	l.root.SetLevel(config.Level)

	if config.Output != nil {
		l.root.SetOutput(config.Output)
	} else {
		l.root.SetOutput(os.Stdout)
	}

	l.root.SetFormatter(&ConsoleFormatter{
		DisableTimestamp: false,
		EnableColors:     !config.DisableColor,
		EnableEntryOrder: true,
	})

	l.root.AddHook(NewCallerHook(config.Level))
	log.SetGlobal(l)
}
