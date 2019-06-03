package logger

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"***REMOVED***/darkroom/server/config"
	"sync"
)

var instance *logrus.Logger
var once sync.Once

const jsonLoggerType = "json"

func getLogger() *logrus.Logger {
	once.Do(func() {
		instance = newLogger()
	})
	return instance
}

func newLogger() *logrus.Logger {
	logger := &logrus.Logger{
		Out:   os.Stdout,
		Hooks: make(logrus.LevelHooks),
		Level: getLogLevel(),
	}
	if config.LogFormat() == jsonLoggerType {
		logger.Formatter = &logrus.JSONFormatter{}
	} else {
		logger.Formatter = &logrus.TextFormatter{}
	}
	return logger
}

func AddHook(hook logrus.Hook) {
	getLogger().Hooks.Add(hook)
}

func Debug(args ...interface{}) {
	getLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	getLogger().Debugf(format, args...)
}

func Debugln(args ...interface{}) {
	getLogger().Debugln(args...)
}

func Error(args ...interface{}) {
	getLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	getLogger().Errorf(format, args...)
}

func Errorln(args ...interface{}) {
	getLogger().Errorln(args...)
}

func Fatal(args ...interface{}) {
	getLogger().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	getLogger().Fatalf(format, args...)
}

func Fatalln(args ...interface{}) {
	getLogger().Fatalln(args...)
}

func Info(args ...interface{}) {
	getLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	getLogger().Infof(format, args...)
}

func Infoln(args ...interface{}) {
	getLogger().Infoln(args...)
}

func Warn(args ...interface{}) {
	getLogger().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	getLogger().Warnf(format, args...)
}

func Warnln(args ...interface{}) {
	getLogger().Warnln(args...)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return getLogger().WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return getLogger().WithFields(fields)
}

func WithRequest(r *http.Request) *logrus.Entry {
	return getLogger().WithFields(logrus.Fields{
		"Method": r.Method,
		"Host":   r.Host,
		"Path":   r.URL.Path,
	})
}
