package server

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServer(t *testing.T) {
	s := NewServer(WithHandler(mux.NewRouter()))
	s.AddLifeCycleHook(NewLifeCycleHook(func() {}, func() {}))
	assert.NotNil(t, s)
}
