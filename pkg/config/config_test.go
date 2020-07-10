package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigCasesWithStringValues(t *testing.T) {
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
			key:      "cache.time",
			callFunc: CacheTime,
		},
	}
	for _, c := range cases {
		assert.Equal(t, v.GetInt(c.key), c.callFunc())
	}

	assert.Equal(t, 0, v.GetInt("nonexistingkey"))
}

func TestConfigCasesWithStringSliceValues(t *testing.T) {
	v := Viper()
	v.Set("defaultParams", "auto=compress")
	Update()
	cases := []struct {
		key      string
		callFunc func() []string
	}{
		{
			key:      "defaultParams",
			callFunc: DefaultParams,
		},
	}

	for _, c := range cases {
		assert.Equal(t, v.GetStringSlice(c.key), c.callFunc())
	}
}
