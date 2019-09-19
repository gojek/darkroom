package service

import (
	"testing"

	"github.com/gojek/darkroom/pkg/config"
	"github.com/gojek/darkroom/pkg/storage/aws/cloudfront"
	"github.com/gojek/darkroom/pkg/storage/aws/s3"
	"github.com/gojek/darkroom/pkg/storage/webfolder"
	"github.com/stretchr/testify/assert"
)

func TestNewDependencies(t *testing.T) {
	deps := NewDependencies()
	assert.NotNil(t, deps)
	assert.Nil(t, deps.Storage)
}

func TestGetDefaultParams(t *testing.T) {
	cases := []struct {
		defaultParams string
		expectedRes   map[string]string
	}{
		{
			defaultParams: "foo=bar",
			expectedRes:   map[string]string{"foo": "bar"},
		},
		{
			defaultParams: "foo=foo,bar",
			expectedRes:   map[string]string{"foo": "foo,bar"},
		},
		{
			defaultParams: "invalid",
			expectedRes:   map[string]string{},
		},
	}
	for _, c := range cases {
		v := config.Viper()
		v.Set("defaultParams", c.defaultParams)
		config.Update()

		assert.Equal(t, c.expectedRes, getDefaultParams())
	}
}

func TestNewDependenciesWithWebFolderStorage(t *testing.T) {
	v := config.Viper()
	v.Set("source.kind", "WebFolder")
	v.Set("source.baseURL", "https://example.com/path/to/folder")
	config.Update()

	deps := NewDependencies()
	assert.NotNil(t, deps)
	assert.IsType(t, &webfolder.Storage{}, deps.Storage)
}

func TestNewDependenciesWithS3Storage(t *testing.T) {
	v := config.Viper()
	v.Set("source.kind", "S3")
	config.Update()

	deps := NewDependencies()
	assert.NotNil(t, deps)
	assert.IsType(t, &s3.Storage{}, deps.Storage)
}

func TestNewDependenciesWithCloudfrontStorage(t *testing.T) {
	v := config.Viper()
	v.Set("source.kind", "Cloudfront")
	v.Set("source.secureProtocol", "true")
	config.Update()

	deps := NewDependencies()
	assert.NotNil(t, deps)
	assert.IsType(t, &cloudfront.Storage{}, deps.Storage)
}
