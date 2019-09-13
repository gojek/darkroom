package config

import (
	"strings"
	"sync"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gojek/darkroom/pkg/storage"
)

type config struct {
	logLevel                        string
	debugMode                       bool
	port                            int
	cacheTime                       int
	dataSource                      Source
	enableConcurrentOpacityChecking bool
	defaultParams                   string
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

	s := Source{
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
		dataSource:                      s,
		enableConcurrentOpacityChecking: v.GetBool("enableConcurrentOpacityChecking"),
		defaultParams:                   v.GetString("defaultParams"),
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

// DataSource returns the source struct after it is initialised from the environment values
func DataSource() *Source {
	return &getConfig().dataSource
}

// ConcurrentOpacityCheckingEnabled returns true if we want to process image using multiple cores (checking isOpaque)
func ConcurrentOpacityCheckingEnabled() bool {
	return getConfig().enableConcurrentOpacityChecking
}

// DefaultParams returns []string of default parameters (separated by semicolon) which will be applied to all image request, following the existing contract
func DefaultParams() []string {
	return strings.Split(getConfig().defaultParams, ";")
}
