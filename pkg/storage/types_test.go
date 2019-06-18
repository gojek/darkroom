package storage

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewResponse(t *testing.T) {
	err := errors.New("randomError")
	r := NewResponse([]byte("randomBytes"), http.StatusBadRequest, err)

	assert.Equal(t, []byte("randomBytes"), r.Data())
	assert.Equal(t, http.StatusBadRequest, r.Status())
	assert.Equal(t, err, r.Error())
}
