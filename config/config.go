package config

import (
	"sync"
)

type config struct {
	logger  loggerConfig
	appInfo appInfo
}

var instance *config
var once sync.Once

func getConfig() *config {
	once.Do(func() {
		instance = newConfig()
	})
	return instance
}

func newConfig() *config {
	initViper()
	return &config{
		logger: loggerConfig{
			level:  getString("log.level"),
			format: getString("log.format"),
		},
		appInfo: appInfo{
			name:        getString("app.name"),
			version:     getString("app.version"),
			description: getString("app.description"),
		},
	}
}

func LogLevel() string {
	return getConfig().logger.level
}

func LogFormat() string {
	return getConfig().logger.format
}

func AppName() string {
	return getConfig().appInfo.name
}

func AppVersion() string {
	return getConfig().appInfo.version
}

func AppDescription() string {
	return getConfig().appInfo.description
}
