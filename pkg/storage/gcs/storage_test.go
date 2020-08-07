package gcs

import (
	"context"
	"github.com/stretchr/testify/suite"
	"google.golang.org/api/googleapi"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

const (
	validPath   = "path/to/valid-file"
	invalidPath = "path/to/invalid-file"
)

type StorageTestSuite struct {
	suite.Suite
	storage Storage
}

func (s *StorageTestSuite) SetupTest() {
	s.storage = *NewStorage()
	s.storage.bucketHandle = &mockBucketHandle{}
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

func (s *StorageTestSuite) TestStorage_GetFailure() {
	res := s.storage.Get(context.Background(), invalidPath)

	s.NotNil(res.Error())
	s.Equal([]byte(nil), res.Data())
	s.Equal(http.StatusNotFound, res.Status())
}

type mockBucketHandle struct{}

func (m mockBucketHandle) Object(s string) ObjectHandle {
	return &mockObjectHandle{objectKey: s}
}

type mockObjectHandle struct {
	objectKey string
}

func (m mockObjectHandle) NewReader(ctx context.Context) (Reader, error) {
	if m.objectKey == validPath {
		return ioutil.NopCloser(strings.NewReader("someData")), nil
	}
	if m.objectKey == invalidPath {
		return nil, &googleapi.Error{
			Code:    404,
			Message: "Not Found",
		}
	}
	return nil, &googleapi.Error{
		Code:    400,
		Message: "Bad Request",
	}
}

func (m mockObjectHandle) NewRangeReader(ctx context.Context, i int64, i2 int64) (Reader, error) {
	// TODO
	panic("implement me")
}
