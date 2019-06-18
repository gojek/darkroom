package s3

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_getStatusCodeFromError(t *testing.T) {
	assert.Equal(t, http.StatusForbidden, getStatusCodeFromError(errors.New("status code: 403")))
	assert.Equal(t, http.StatusNotFound, getStatusCodeFromError(errors.New("status code: 404")))
	assert.Equal(t, http.StatusUnauthorized, getStatusCodeFromError(errors.New("status code: 401")))
	assert.Equal(t, http.StatusUnprocessableEntity, getStatusCodeFromError(errors.New("status code: 4xx")))
	assert.Equal(t, http.StatusUnprocessableEntity, getStatusCodeFromError(errors.New("status code: 422")))
	assert.Equal(t, http.StatusOK, getStatusCodeFromError(nil))
}
