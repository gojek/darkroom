package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

func initViper() {
	viper.AutomaticEnv()
	viper.SetConfigName("application")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}

var replacer = strings.NewReplacer(".", "_")

func getString(key string) string {
	envKey := strings.ToUpper(replacer.Replace(key))
	value := os.Getenv(envKey)
	if value == "" {
		value = viper.GetString(key)
	}
	return value
}
