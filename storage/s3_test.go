package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestNewS3Storage(t *testing.T) {
	s3s := NewS3Storage()
	assert.NotNil(t, s3s)
}

func TestS3Storage_Get(t *testing.T) {
	mockStorage := &MockBaseStorage{}
	s3s := S3Storage{
		base: mockStorage,
		pathPrefix: "/path/to",
	}
	mockStorage.On("Get", mock.Anything, mock.Anything).Return([]byte(nil), http.StatusOK, nil)
	s3s.Get(context.Background(), "/inner-path")
}

func TestS3Storage_GetWithoutPathPrefix(t *testing.T) {
	mockStorage := &MockBaseStorage{}
	s3s := S3Storage{
		base: mockStorage,
	}
	mockStorage.On("Get", mock.Anything, mock.Anything).Return([]byte(nil), http.StatusOK, nil)
	s3s.Get(context.Background(), "/inner-path")
}

