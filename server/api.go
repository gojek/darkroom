package server

import (
	"***REMOVED***/darkroom/server/config"
	"***REMOVED***/darkroom/server/logger"
)

func Start() {
	logger.Infof("Starting %s", config.AppName())
}