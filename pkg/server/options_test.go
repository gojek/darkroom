package server

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithHandler(t *testing.T) {
	r := mux.NewRouter()
	s := NewServer(WithHandler(r))
	assert.Equal(t, r, s.handler)
}
