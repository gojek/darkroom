package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/stretchr/testify/suite"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

const (
	validPath   = "path/to/valid-file"
	invalidPath = "path/to/invalid-file"
	bucketName  = "bucket-name"
)

type StorageTestSuite struct {
	suite.Suite
	storage Storage
}

func (s *StorageTestSuite) SetupTest() {
	s.storage = *NewStorage(Options{BucketName: bucketName})
	s.storage.bucketHandle = &mockBucketHandle{}
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (s *StorageTestSuite) TestNewStorage() {
	s.NotNil(s.storage)
}

func (s *StorageTestSuite) TestNewStorageHasBucketHandle() {
	ns := *NewStorage(Options{BucketName: bucketName})
	s.NotNil(ns)
	s.NotNil(ns.bucketHandle)
}

func (s *StorageTestSuite) TestNewStorageHasCorrectBucketName() {
	ns := *NewStorage(Options{BucketName: bucketName})
	s.NotNil(ns)
	c, _ := storage.NewClient(
		context.TODO(),
		option.WithHTTPClient(newTestClient(bucketResponseMocker)),
	)
	ns.bucketHandle = bucketHandle{c.Bucket(bucketName)}
	attrs, err := ns.bucketHandle.Attrs(context.TODO())
	s.NoError(err)
	s.Equal(bucketName, attrs.Name)
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

func (m mockBucketHandle) Attrs(context.Context) (*storage.BucketAttrs, error) {
	return &storage.BucketAttrs{Name: bucketName}, nil
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
