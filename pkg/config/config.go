package config

import (
	"github.com/afex/hystrix-go/hystrix"
	"***REMOVED***/darkroom/core/pkg/storage"
	"sync"
)

type config struct {
	logLevel  string
	app       app
	debugMode bool
	port      int
	cacheTime int
	source    source
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

	s := source{
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
		logLevel: v.GetString("log.level"),
		app: app{
			name:        v.GetString("app.Name"),
			version:     v.GetString("app.version"),
			description: v.GetString("app.description"),
		},
		debugMode: v.GetBool("debug"),
		port:      port,
		cacheTime: v.GetInt("cache.time"),
		source:    s,
	}
}

func Update() {
	instance = newConfig()
}

func LogLevel() string {
	return getConfig().logLevel
}

func AppName() string {
	return getConfig().app.name
}

func AppVersion() string {
	return getConfig().app.version
}

func AppDescription() string {
	return getConfig().app.description
}

func DebugModeEnabled() bool {
	return getConfig().debugMode
}

func Port() int {
	return getConfig().port
}

func CacheTime() int {
	return getConfig().cacheTime
}

func Source() *source {
	return &getConfig().source
}
