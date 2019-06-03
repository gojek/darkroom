package config

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigCases(t *testing.T) {
	initViper()
	cases := []struct {
		key      string
		callFunc func() string
	}{
		{
			key:      "app.name",
			callFunc: AppName,
		},
		{
			key:      "app.version",
			callFunc: AppVersion,
		},
		{
			key:      "app.description",
			callFunc: AppDescription,
		},
		{
			key:      "log.level",
			callFunc: LogLevel,
		},
		{
			key:      "log.format",
			callFunc: LogFormat,
		},
	}

	for _, c := range cases {
		assert.Equal(t, viper.GetString(c.key), c.callFunc())
	}
}
