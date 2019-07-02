package config

import (
	"github.com/spf13/viper"
	"strings"
	"sync"
)

const (
	configFileName = "config"
	configFileExt  = "yaml"
)

var viperInstance *viper.Viper
var viperInit sync.Once

// Viper creates a new singleton instance of viper.Viper
func Viper() *viper.Viper {
	viperInit.Do(func() {
		viperInstance = viper.New()
		viperInstance.AutomaticEnv()
		viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viperInstance.SetConfigName(configFileName)
		viperInstance.AddConfigPath(".")
		viperInstance.SetConfigType(configFileExt)
		_ = viperInstance.ReadInConfig()
	})
	return viperInstance
}
