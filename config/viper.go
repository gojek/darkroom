package config

import (
	"github.com/spf13/viper"
	"strings"
)

const (
	configFileName = "config"
	configFileExt  = "yaml"
)

func initViper() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName(configFileName)
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.SetConfigType(configFileExt)
	_ = viper.ReadInConfig()
}
