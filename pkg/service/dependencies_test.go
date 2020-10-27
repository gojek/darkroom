package service

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gojek/darkroom/pkg/config"
	"github.com/gojek/darkroom/pkg/storage/aws/cloudfront"
	"github.com/gojek/darkroom/pkg/storage/aws/s3"
	"github.com/gojek/darkroom/pkg/storage/webfolder"
	"github.com/stretchr/testify/assert"
)

func TestNewDependencies(t *testing.T) {
	deps, err := NewDependencies(prometheus.NewRegistry())
	assert.Error(t, err)
	assert.Nil(t, deps)
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
	v.Set("metrics.system", "prometheus")
	config.Update()

	deps, err := NewDependencies(prometheus.NewRegistry())
	assert.NoError(t, err)
	assert.NotNil(t, deps)
	assert.IsType(t, &webfolder.Storage{}, deps.Storage)
}

func TestNewDependenciesWithS3Storage(t *testing.T) {
	v := config.Viper()
	v.Set("source.kind", "S3")
	v.Set("metrics.system", "prometheus")
	config.Update()

	deps, err := NewDependencies(prometheus.NewRegistry())
	assert.NoError(t, err)
	assert.NotNil(t, deps)
	assert.IsType(t, &s3.Storage{}, deps.Storage)
}

func TestNewDependenciesWithCloudfrontStorage(t *testing.T) {
	v := config.Viper()
	v.Set("source.kind", "Cloudfront")
	v.Set("source.secureProtocol", "true")
	v.Set("metrics.system", "prometheus")
	config.Update()

	deps, err := NewDependencies(prometheus.NewRegistry())
	assert.NoError(t, err)
	assert.NotNil(t, deps)
	assert.IsType(t, &cloudfront.Storage{}, deps.Storage)
}
