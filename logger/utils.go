package logger

import (
	"github.com/sirupsen/logrus"
	"***REMOVED***/darkroom/core/config"
)

func getLogLevel() logrus.Level {
	if config.LogLevel() == "" {
		return logrus.DebugLevel
	}
	level, err := logrus.ParseLevel(config.LogLevel())
	if err != nil {
		panic(err)
	}
	return level
}
