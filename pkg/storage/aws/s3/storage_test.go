package s3

import (
	"context"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"***REMOVED***/darkroom/core/pkg/storage"
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

func (s *StorageTestSuite) SetupSuite() {
	s.storage = *NewStorage(
		WithBucketName("bucket"),
		WithBucketRegion("region"),
		WithAccessKey("randomAccessKey"),
		WithSecretKey("randomSecretKey"),
		WithHystrixCommand(storage.HystrixCommand{
			Name: "TestCommand",
			Config: hystrix.CommandConfig{
				Timeout:                2000,
				MaxConcurrentRequests:  100,
				RequestVolumeThreshold: 10,
				SleepWindow:            10,
				ErrorPercentThreshold:  25,
			},
		}),
	)
	s.storage.downloader = &mockDownloader{}
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (s *StorageTestSuite) TestStorage_Get() {
	res := s.storage.Get(context.Background(), validPath)

	assert.Nil(s.T(), res.Error())
	assert.Equal(s.T(), []byte("someData"), res.Data())
	assert.Equal(s.T(), http.StatusOK, res.Status())
}

func (s *StorageTestSuite) TestStorage_GetFailure() {
	res := s.storage.Get(context.Background(), invalidPath)

	assert.NotNil(s.T(), res.Error())
	assert.Equal(s.T(), []byte(nil), res.Data())
	assert.Equal(s.T(), http.StatusUnprocessableEntity, res.Status())
}

type mockDownloader struct {
	mock.Mock
}

func (d *mockDownloader) Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
	return d.DownloadWithContext(aws.BackgroundContext(), w, input, options...)
}

func (d *mockDownloader) DownloadWithContext(ctx aws.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
	if *input.Key == validPath {
		_, _ = w.WriteAt([]byte("someData"), 0)
		return 0, nil
	}
	return 0, errors.New("error")
}
