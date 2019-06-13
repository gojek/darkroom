package router

import (
	"github.com/stretchr/testify/assert"
	"***REMOVED***/darkroom/core/config"
	"***REMOVED***/darkroom/core/service"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter(&service.Dependencies{})
	assert.NotNil(t, router)
}

func TestNewRouterInDebugMode(t *testing.T) {
	v := config.Viper()
	v.Set("debug", "true")
	config.Update()

	router := NewRouter(&service.Dependencies{})
	assert.NotNil(t, router)
}

func TestNewRouterWithPathPrefix(t *testing.T) {
	v := config.Viper()
	v.Set("source.kind", "s3")
	v.Set("source.pathPrefix", "/path/to/folder")
	config.Update()

	router := NewRouter(&service.Dependencies{})
	assert.NotNil(t, router)
}