package config

import (
	"sync"
)

type config struct {
	logger     loggerConfig
	appInfo    appInfo
	bucketInfo bucketInfo
	debugMode  bool
	port       int
	cacheTime  int
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
	port := getInt("port")
	if port == 0 {
		port = 3000 // Fallback to default port
	}
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
		bucketInfo: bucketInfo{
			name:       getString("bucket.name"),
			accessKey:  getString("bucket.accessKey"),
			secretKey:  getString("bucket.secretKey"),
			pathPrefix: getString("bucket.pathPrefix"),
		},
		debugMode: getFeature("debug"),
		port:      port,
		cacheTime: getInt("cache.time"),
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

func DebugModeEnabled() bool {
	return getConfig().debugMode
}

func Port() int {
	return getConfig().port
}

func BucketName() string {
	return getConfig().bucketInfo.name
}

func BucketAccessKey() string {
	return getConfig().bucketInfo.accessKey
}

func BucketSecretKey() string {
	return getConfig().bucketInfo.secretKey
}

func BucketPathPrefix() string {
	return getConfig().bucketInfo.pathPrefix
}

func CacheTime() int {
	return getConfig().cacheTime
}
