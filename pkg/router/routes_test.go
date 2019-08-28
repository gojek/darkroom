package router

import (
	"context"
	"net/http"
	"testing"

	"github.com/gojek/darkroom/pkg/config"
	"github.com/gojek/darkroom/pkg/service"
	"github.com/gojek/darkroom/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter(&service.Dependencies{Storage: &mockStorage{}, Manipulator: &service.MockManipulator{}})
	assert.NotNil(t, router)
}

func TestNewRouterInDebugMode(t *testing.T) {
	v := config.Viper()
	v.Set("debug", "true")
	config.Update()

	router := NewRouter(&service.Dependencies{Storage: &mockStorage{}, Manipulator: &service.MockManipulator{}})
	assert.NotNil(t, router)
}

func TestNewRouterWithPathPrefix(t *testing.T) {
	v := config.Viper()
	v.Set("source.kind", "s3")
	v.Set("source.pathPrefix", "/path/to/folder")
	config.Update()

	router := NewRouter(&service.Dependencies{Storage: &mockStorage{}, Manipulator: &service.MockManipulator{}})
	assert.NotNil(t, router)
}

type mockStorage struct {
}

func (m *mockStorage) Get(ctx context.Context, path string) storage.IResponse {
	return storage.NewResponse([]byte(nil), http.StatusOK, nil)
}
