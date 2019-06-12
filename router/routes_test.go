package router

import (
	"github.com/stretchr/testify/assert"
	"***REMOVED***/darkroom/core/service"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter(&service.Dependencies{})
	assert.NotNil(t, router)
}
