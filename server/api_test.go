package server

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServer(t *testing.T) {
	s := NewServer(WithHandler(mux.NewRouter()))
	assert.NotNil(t, s)
}
