package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gojek/darkroom/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	validPath    = "path/to/valid-file"
	invalidPath  = "path/to/invalid-file"
	validRange   = "bytes=0-100"
	invalidRange = "none"
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
	s.storage.service = &mockGetObject{}
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

func (s *StorageTestSuite) TestStorage_GetRange() {
	opt := &storage.GetPartialObjectRequestOptions{Range: validRange}
	res := s.storage.GetPartialObject(context.Background(), validPath, opt)
	metadata := storage.ResponseMetadata{
		AcceptRanges:  "bytes",
		ContentLength: "101",
		ContentRange:  "bytes 100-200/247103",
		ContentType:   "image/png",
		ETag:          "32705ce195789d7bf07f3d44783c2988",
		LastModified:  "Wed, 21 Oct 2015 07:28:00 GMT",
	}

	assert.Nil(s.T(), res.Error())
	assert.Equal(s.T(), []byte("someData"), res.Data())
	assert.Equal(s.T(), http.StatusPartialContent, res.Status())
	assert.Equal(s.T(), &metadata, res.Metadata())
}

func (s *StorageTestSuite) TestStorage_GetRangeFailure() {
	opt := &storage.GetPartialObjectRequestOptions{Range: invalidRange}
	res := s.storage.GetPartialObject(context.Background(), validPath, opt)

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

type mockGetObject struct {
	mock.Mock
	s3iface.S3API
}

func (d *mockGetObject) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	if *input.Range == validRange {
		t, _ := time.Parse(http.TimeFormat, "Wed, 21 Oct 2015 07:28:00 GMT")
		return &s3.GetObjectOutput{
			AcceptRanges:  aws.String("bytes"),
			ContentLength: aws.Int64(101),
			ContentRange:  aws.String("bytes 100-200/247103"),
			ContentType:   aws.String("image/png"),
			ETag:          aws.String("32705ce195789d7bf07f3d44783c2988"),
			LastModified:  aws.Time(t),
			Body:          ioutil.NopCloser(bytes.NewReader([]byte("someData"))),
		}, nil
	}
	return nil, errors.New("error")
}
