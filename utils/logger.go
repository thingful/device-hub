// Copyright © 2017 thingful
package utils

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

const (
	loggerContextKey = "logger-context-key"
)

// NewLogger creates and returns a new instance of the Logger type. This
// initializes the logger with a tagged logrus Entry initialized with the
// version string and hostname.
func NewLogger(version, logpath string) Logger {
	log := logrus.New()
	log.Formatter = new(logrus.JSONFormatter)
	log.Level = logrus.InfoLevel
	// log.Out = os.Stdout
	if len(logpath) != 0 {
		log.Out = &lumberjack.Logger{
			Filename:   logpath,
			MaxSize:    1,
			MaxBackups: 3,
			MaxAge:     28,
		}
	} else {
		log.Out = os.Stdout
	}
	logger := log.WithFields(defaultFields(version))
	return &l{entry: logger}
}

// NewNoOpLogger returns a noop Logger
func NewNoOpLogger() Logger {
	lgr := logrus.New()
	lgr.Out = ioutil.Discard
	logger := lgr.WithFields(logrus.Fields{})
	return &l{entry: logger}
}

func defaultFields(version string) logrus.Fields {
	hostname, err := os.Hostname()

	if err != nil {
		hostname = "UNKNOWN"
	}

	pid := os.Getpid()

	return logrus.Fields{
		"name":     "pomelo",
		"version":  version,
		"hostname": hostname,
		"pid":      pid,
	}
}

type Logger interface {
	Error(args ...interface{})
	Info(args ...interface{})
	TimeAsInfo(start time.Time, message ...string)
}

type l struct {
	entry *logrus.Entry
}

// TimeAsInfo can be used in a defer to time execution
func (l *l) TimeAsInfo(start time.Time, message ...string) {
	elapsed := time.Since(start)

	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])

	l.entry.WithFields(logrus.Fields{
		"caller":   f.Name(),
		"duration": elapsed.String(),
	}).Info(strings.Join(message, ","))
}

func (l *l) Info(args ...interface{}) {
	l.entry.Info(args)
}

func (l *l) Error(args ...interface{}) {
	l.entry.Error(args)
}
