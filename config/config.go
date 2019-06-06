package config

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/spf13/viper"
	"***REMOVED***/darkroom/storage"
	"sync"
)

type config struct {
	logger     loggerConfig
	appInfo    appInfo
	bucketInfo bucketInfo
	debugMode  bool
	port       int
	cacheTime  int
	hystrixCmd storage.HystrixCommand
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
	port := viper.GetInt("port")
	if port == 0 {
		port = 3000 // Fallback to default port
	}
	return &config{
		logger: loggerConfig{
			level:  viper.GetString("log.level"),
			format: viper.GetString("log.format"),
		},
		appInfo: appInfo{
			name:        viper.GetString("app.name"),
			version:     viper.GetString("app.version"),
			description: viper.GetString("app.description"),
		},
		bucketInfo: bucketInfo{
			name:       viper.GetString("bucket.name"),
			region:     viper.GetString("bucket.region"),
			accessKey:  viper.GetString("bucket.accessKey"),
			secretKey:  viper.GetString("bucket.secretKey"),
			pathPrefix: viper.GetString("bucket.pathPrefix"),
		},
		debugMode: viper.GetBool("debug"),
		port:      port,
		cacheTime: viper.GetInt("cache.time"),
		hystrixCmd: storage.HystrixCommand{
			Name: viper.GetString("hystrix.command.name"),
			Config: hystrix.CommandConfig{
				Timeout:                viper.GetInt("hystrix.config.timeout"),
				MaxConcurrentRequests:  viper.GetInt("hystrix.config.maxConcurrentRequests"),
				RequestVolumeThreshold: viper.GetInt("hystrix.config.requestVolumeThreshold"),
				SleepWindow:            viper.GetInt("hystrix.config.sleepWindow"),
				ErrorPercentThreshold:  viper.GetInt("hystrix.config.errorPercentThreshold"),
			},
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

func DebugModeEnabled() bool {
	return getConfig().debugMode
}

func Port() int {
	return getConfig().port
}

func BucketName() string {
	return getConfig().bucketInfo.name
}

func BucketRegion() string {
	return getConfig().bucketInfo.region
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

func HystrixCommand() storage.HystrixCommand {
	return getConfig().hystrixCmd
}
