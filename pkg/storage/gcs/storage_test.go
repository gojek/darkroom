package gcs

import (
	"context"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

const (
	validPath    = "path/to/valid-file"
)

type StorageTestSuite struct {
	suite.Suite
	storage Storage
}

func (s *StorageTestSuite) SetupTest() {
	s.storage = *NewStorage()
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (s *StorageTestSuite) TestNewStorage() {
	s.NotNil(s.storage)
}

func (s *StorageTestSuite) TestStorage_Get() {
	res := s.storage.Get(context.Background(), validPath)

	s.Nil(res.Error())
	s.Equal([]byte("someData"), res.Data())
	s.Equal(http.StatusOK, res.Status())
}
