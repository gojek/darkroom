package logger

import (
	"github.com/sirupsen/logrus"
	"***REMOVED***/darkroom/server/config"
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
