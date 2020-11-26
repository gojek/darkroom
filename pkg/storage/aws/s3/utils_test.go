package s3

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getStatusCodeFromError(t *testing.T) {
	var statusCode int
	assert.Equal(t, http.StatusForbidden, getStatusCodeFromError(errors.New("status code: 403"), &statusCode))
	assert.Equal(t, http.StatusNotFound, getStatusCodeFromError(errors.New("status code: 404"), &statusCode))
	assert.Equal(t, http.StatusUnauthorized, getStatusCodeFromError(errors.New("status code: 401"), &statusCode))
	assert.Equal(t, http.StatusUnprocessableEntity, getStatusCodeFromError(errors.New("status code: 4xx"), &statusCode))
	assert.Equal(t, http.StatusUnprocessableEntity, getStatusCodeFromError(errors.New("status code: 422"), &statusCode))
	assert.Equal(t, http.StatusOK, getStatusCodeFromError(nil, nil))

	statusCode = http.StatusPartialContent
	assert.Equal(t, http.StatusPartialContent, getStatusCodeFromError(nil, &statusCode))
}
