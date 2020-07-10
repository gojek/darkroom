package server

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServer(t *testing.T) {
	s := NewServer(Options{
		Handler:       mux.NewRouter(),
		Port:          3000,
		LifeCycleHook: nil,
	})
	assert.NotNil(t, s)
}
