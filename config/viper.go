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

/* getString() first tries to get key from environment, if not found then checks inside
 * the YAML config file.
 * Note: A key of the format `app.env.var` will result in `APP_ENV_VAR` variable for the
 * OS context while checking.
 */
func getString(key string) string {
	envKey := strings.ToUpper(replacer.Replace(key))
	value := os.Getenv(envKey)
	if value == "" {
		value = viper.GetString(key)
	}
	return value
}
