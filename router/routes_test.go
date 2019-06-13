package router

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"***REMOVED***/darkroom/core/config"
	"***REMOVED***/darkroom/core/service"
	"***REMOVED***/darkroom/storage"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter(&service.Dependencies{Storage: &mockStorage{}, Manipulator: &mockManipulator{}})
	assert.NotNil(t, router)
}

func TestNewRouterInDebugMode(t *testing.T) {
	v := config.Viper()
	v.Set("debug", "true")
	config.Update()

	router := NewRouter(&service.Dependencies{Storage: &mockStorage{}, Manipulator: &mockManipulator{}})
	assert.NotNil(t, router)
}

func TestNewRouterWithPathPrefix(t *testing.T) {
	v := config.Viper()
	v.Set("source.kind", "s3")
	v.Set("source.pathPrefix", "/path/to/folder")
	config.Update()

	router := NewRouter(&service.Dependencies{Storage: &mockStorage{}, Manipulator: &mockManipulator{}})
	assert.NotNil(t, router)
}

type mockStorage struct {
}

func (m *mockStorage) Get(ctx context.Context, path string) storage.IResponse {
	return storage.NewResponse([]byte(nil), http.StatusOK, nil)
}

type mockManipulator struct {
}

func (m *mockManipulator) Process(ctx context.Context, data []byte, params map[string]string) ([]byte, error) {
	return nil, nil
}
