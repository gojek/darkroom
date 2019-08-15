package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigCases(t *testing.T) {
	v := Viper()
	cases := []struct {
		key      string
		callFunc func() string
	}{
		{
			key:      "log.level",
			callFunc: LogLevel,
		},
	}

	for _, c := range cases {
		assert.Equal(t, v.GetString(c.key), c.callFunc())
	}
}

func TestConfigCasesWithFeatureToggle(t *testing.T) {
	v := Viper()
	cases := []struct {
		key      string
		callFunc func() bool
	}{
		{
			key:      "debug",
			callFunc: DebugModeEnabled,
		},
	}
	for _, c := range cases {
		assert.Equal(t, v.GetBool(c.key), c.callFunc())
	}

	assert.Equal(t, false, v.GetBool("nonexistingkey"))
}

func TestConfigCasesWithIntValues(t *testing.T) {
	v := Viper()
	v.Set("port", 3000)
	cases := []struct {
		key      string
		callFunc func() int
	}{
		{
			key:      "port",
			callFunc: Port,
		},
		{
			key:      "cache.time",
			callFunc: CacheTime,
		},
	}
	for _, c := range cases {
		assert.Equal(t, v.GetInt(c.key), c.callFunc())
	}

	assert.Equal(t, 0, v.GetInt("nonexistingkey"))
}
