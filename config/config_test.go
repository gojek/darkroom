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
		{
			key:      "bucket.name",
			callFunc: BucketName,
		},
		{
			key:      "bucket.accessKey",
			callFunc: BucketAccessKey,
		},
		{
			key:      "bucket.secretKey",
			callFunc: BucketSecretKey,
		},
		{
			key:      "bucket.pathPrefix",
			callFunc: BucketPathPrefix,
		},
	}

	for _, c := range cases {
		assert.Equal(t, viper.GetString(c.key), c.callFunc())
	}
}

func TestConfigCasesWithFeatureToggle(t *testing.T) {
	initViper()
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
		assert.Equal(t, viper.GetBool(c.key), c.callFunc())
	}

	assert.Equal(t, false, getFeature("nonexistingkey"))
}

func TestConfigCasesWithIntValues(t *testing.T) {
	initViper()
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
		assert.Equal(t, viper.GetInt(c.key), c.callFunc())
	}

	assert.Equal(t, 0, getInt("nonexistingkey"))
}
