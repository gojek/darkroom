package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
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

func getFeature(key string) bool {
	value := getString(key)
	if value == "" {
		return false
	}

	val, err := strconv.ParseBool(value)
	handleNoKey(key, err)
	return val
}

func getInt(key string) int {
	value := getString(key)
	if value == "" {
		return viper.GetInt(key)
	}
	val, err := strconv.Atoi(value)
	handleNoKey(key, err)
	return val
}

func handleNoKey(key string, err error) {
	if err != nil {
		log.Fatalf("couldn't parse key %s, error: %s", key, err)
	}
}
