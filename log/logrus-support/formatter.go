package logrus_support

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"
)

type ConsoleFormatter struct {
	DisableTimestamp bool
	EnableColors     bool
	EnableEntryOrder bool

	Host string
	IP   string
}

// Formatter implements logrus.Formatter
func (f *ConsoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	if f.EnableColors {
		colorStr := logTypeToColor(entry.Level)
		fmt.Fprintf(b, colorStr)
	}

	if !f.DisableTimestamp {
		fmt.Fprintf(b, "[%s] ", entry.Time.Format(defaultLogTimeFormat))
	}

	fmt.Fprintf(b, "[%s] ", entry.Level.String())

	if f.IP != "" {
		fmt.Fprintf(b, "[%s] ", f.IP)
	}
	if f.Host != "" {
		fmt.Fprintf(b, "[%s] ", f.Host)
	}

	if _, ok := entry.Data["file"]; ok {
		fmt.Fprintf(b, "[%s:%d %s] ", entry.Data["file"], entry.Data["line"], entry.Data["func"])
	}

	fmt.Fprintf(b, "%s", entry.Message)

	if f.EnableEntryOrder {
		keys := make([]string, 0, len(entry.Data))
		for k := range entry.Data {
			if k != "file" && k != "line" && k != "func" {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(b, " %v=%v", k, entry.Data[k])
		}
	} else {
		for k, v := range entry.Data {
			if k != "file" && k != "line" && k != "func" {
				fmt.Fprintf(b, " %v=%v", k, v)
			}
		}
	}

	b.WriteByte('\n')

	if f.EnableColors {
		b.WriteString(Reset)
	}
	return b.Bytes(), nil
}

const (
	defaultLogTimeFormat = "2006-01-02T15:04:05.000"
)

// logTypeToColor converts the Level to a color string.
func logTypeToColor(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return White
	case logrus.InfoLevel:
		return Gray
	case logrus.WarnLevel:
		return Yellow
	case logrus.ErrorLevel:
		return Red
	case logrus.FatalLevel:
		return RedBold
	case logrus.PanicLevel:
		return RedBold
	}

	return Gray
}

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Gray    = "\033[38m"

	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)
