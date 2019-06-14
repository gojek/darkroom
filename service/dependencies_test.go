package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"***REMOVED***/darkroom/core/pkg/config"
	"***REMOVED***/darkroom/storage/cloudfront"
	"***REMOVED***/darkroom/storage/s3"
	"***REMOVED***/darkroom/storage/webfolder"
)

func TestNewDependencies(t *testing.T) {
	deps := NewDependencies()
	assert.NotNil(t, deps)
	assert.Nil(t, deps.Storage)
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
