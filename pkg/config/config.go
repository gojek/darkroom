package config

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gojek/darkroom/pkg/storage"
	"sync"
)

type config struct {
	logLevel                        string
	debugMode                       bool
	port                            int
	cacheTime                       int
	source                          CloudSource
	enableConcurrentOpacityChecking bool
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
	v := Viper()
	port := v.GetInt("port")
	if port == 0 {
		port = 3000 // Fallback to default port
	}

	s := CloudSource{
		Kind: v.GetString("source.kind"),
		HystrixCommand: storage.HystrixCommand{
			Name: v.GetString("source.hystrix.commandName"),
			Config: hystrix.CommandConfig{
				Timeout:                v.GetInt("source.hystrix.timeout"),
				MaxConcurrentRequests:  v.GetInt("source.hystrix.maxConcurrentRequests"),
				RequestVolumeThreshold: v.GetInt("source.hystrix.requestVolumeThreshold"),
				SleepWindow:            v.GetInt("source.hystrix.sleepWindow"),
				ErrorPercentThreshold:  v.GetInt("source.hystrix.errorPercentThreshold")},
		},
		PathPrefix: v.GetString("source.pathPrefix"),
	}
	s.readValue()

	return &config{
		logLevel:                        v.GetString("log.level"),
		debugMode:                       v.GetBool("debug"),
		port:                            port,
		cacheTime:                       v.GetInt("cache.time"),
		source:                          s,
		enableConcurrentOpacityChecking: v.GetBool("enableConcurrentOpacityChecking"),
	}
}

// Update creates a new instance of the configuration and reads all values again
func Update() {
	instance = newConfig()
}

// LogLevel returns the log level for logger from the environment
func LogLevel() string {
	return getConfig().logLevel
}

// DebugModeEnabled returns the debug mode bool from the environment
func DebugModeEnabled() bool {
	return getConfig().debugMode
}

// Port returns the application port to be used from the environment
func Port() int {
	return getConfig().port
}

// CacheTime returns the time to set the cache-time in image handler from the environment
func CacheTime() int {
	return getConfig().cacheTime
}

// Source returns the source struct after it is initialised from the environment values
func Source() *CloudSource {
	return &getConfig().source
}

// ConcurrentOpacityCheckingEnabled returns true if we want to process image using multiple cores (checking isOpaque)
func ConcurrentOpacityCheckingEnabled() bool {
	return getConfig().enableConcurrentOpacityChecking
}
