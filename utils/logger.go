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
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

const (
	loggerContextKey = "logger-context-key"
)

// NewLogger creates and returns a new instance of the Logger type. This
// initializes the logger with a tagged logrus Entry initialized with the
// version string, hostname and two logging options (syslog & logpath)
// If syslog is true the logger will use local syslog method
// If syslog is false logpath will be evaluated as log file with rotation
// STDOUT is the logging fallback method
func NewLogger(version string, syslog bool, logpath string) Logger {
	log := logrus.New()
	if syslog {
		hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_INFO, "device-hub")
		if err == nil {
			log.Hooks.Add(hook)
			return &l{entry: log}
		}
	}
	log.Formatter = new(logrus.JSONFormatter)
	log.Level = logrus.InfoLevel
	// Test values for rotation, these should be parametrized
	if len(logpath) != 0 {
		log.Out = &lumberjack.Logger{
			Filename:   logpath,
			MaxSize:    1, // Mb
			MaxBackups: 3,
			MaxAge:     28, // days
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
		"name":     "device-hub",
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
