package logrus_support

import (
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type CallerHook struct {
	level       logrus.Level
	skipPkg     []string
	maxFilePath int
}

func NewCallerHook(level logrus.Level, skipPath []string) *CallerHook {
	c := &CallerHook{
		level: level,
		skipPkg: []string{
			"github.com/sirupsen/logrus",
			"github.com/ragpanda/go-toolkit/log/logrus-support",
			"github.com/ragpanda/go-toolkit/log",
		},
		maxFilePath: 3,
	}

	c.skipPkg = append(c.skipPkg, skipPath...)

	return c
}

func (self *CallerHook) Levels() (r []logrus.Level) {
	for _, v := range logrus.AllLevels {
		if self.level <= v {
			r = append(r, v)
		}
	}
	return logrus.AllLevels

}

func (hook *CallerHook) Fire(entry *logrus.Entry) error {
	pc := make([]uintptr, 10)
	cnt := runtime.Callers(6, pc)

	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		if !hook.isSkippedPackageName(name) {
			file, line := fu.FileLine(pc[i] - 1)
			pathSp := strings.Split(file, "/")
			if len(pathSp) > hook.maxFilePath {
				file = path.Join(pathSp[len(pathSp)-hook.maxFilePath:]...)
			}

			funcSp := strings.Split(name, "/")
			if len(funcSp) > 2 {
				name = path.Join(funcSp[len(funcSp)-2:]...)
			}

			entry.Data["file"] = file
			entry.Data["line"] = line
			entry.Data["func"] = name
			break
		}
	}

	return nil
}

func (hook *CallerHook) isSkippedPackageName(name string) bool {
	for _, pkgName := range hook.skipPkg {
		if strings.Contains(name, pkgName) {
			return true
		}
	}
	return false

}
