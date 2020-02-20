package cloudfront

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gojek/darkroom/pkg/storage"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	validHost   = "cloudfront.net"
	validPath   = "/path/to/valid-file"
	invalidPath = "/path/to/invalid-file"
	validRange  = "bytes=100-200"
)

type StorageTestSuite struct {
	suite.Suite
	storage Storage
	client  *mockClient
}

func (s *StorageTestSuite) SetupTest() {
	s.client = &mockClient{}
	s.storage = *NewStorage(
		WithCloudfrontHost(validHost),
		WithHeimdallClient(s.client),
		WithSecureProtocol(),
	)
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (s *StorageTestSuite) TestNewStorage() {
	assert.NotNil(s.T(), s.storage)
}

func (s *StorageTestSuite) TestStorage_GetNotFound() {
	s.client.On("Get", fmt.Sprintf("%s://%s%s", s.storage.getProtocol(), validHost, invalidPath), http.Header(nil)).
		Return(&http.Response{StatusCode: http.StatusNotFound}, errors.New("not found"))

	res := s.storage.Get(context.TODO(), invalidPath)

	assert.NotNil(s.T(), res.Error())
	assert.Equal(s.T(), http.StatusNotFound, res.Status())
	assert.Nil(s.T(), res.Data())
}

func (s *StorageTestSuite) TestStorage_GetNoResponse() {
	s.storage.secureProtocol = false // Use http
	s.client.On("Get", fmt.Sprintf("%s://%s%s", s.storage.getProtocol(), validHost, invalidPath), http.Header(nil)).
		Return(nil, errors.New("response body read failure"))

	res := s.storage.Get(context.TODO(), invalidPath)

	assert.NotNil(s.T(), res.Error())
	assert.Equal(s.T(), http.StatusUnprocessableEntity, res.Status())
	assert.Nil(s.T(), res.Data())
}

func (s *StorageTestSuite) TestStorage_GetForbidden() {
	s.client.On("Get", fmt.Sprintf("%s://%s%s", s.storage.getProtocol(), validHost, invalidPath), http.Header(nil)).
		Return(&http.Response{
			StatusCode: http.StatusForbidden,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("response body"))),
		}, nil)

	res := s.storage.Get(context.TODO(), invalidPath)

	assert.NotNil(s.T(), res.Error())
	assert.Equal(s.T(), http.StatusForbidden, res.Status())
	assert.Equal(s.T(), []byte(nil), res.Data())
}

func (s *StorageTestSuite) TestStorage_GetSuccessResponse() {
	s.client.On("Get", fmt.Sprintf("%s://%s%s", s.storage.getProtocol(), validHost, validPath), http.Header(nil)).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("response body"))),
		}, nil)

	res := s.storage.Get(context.TODO(), validPath)

	assert.Nil(s.T(), res.Error())
	assert.Equal(s.T(), http.StatusOK, res.Status())
	assert.Equal(s.T(), []byte("response body"), res.Data())
}

func (s *StorageTestSuite) TestStorage_GetPartialObjectSuccessResponse() {
	metadata := storage.ResponseMetadata{
		AcceptRanges:  "bytes",
		ContentLength: "1024",
		ContentType:   "image/png",
		ContentRange:  "bytes 100-200/1024",
		ETag:          "32705ce195789d7bf07f3d44783c2988",
		LastModified:  "Wed, 21 Oct 2015 07:28:00 GMT",
	}

	reqHeader := http.Header{}
	reqHeader.Add(storage.HeaderRange, validRange)

	respHeader := http.Header{}
	respHeader.Add(storage.HeaderAcceptRanges, metadata.AcceptRanges)
	respHeader.Add(storage.HeaderContentLength, metadata.ContentLength)
	respHeader.Add(storage.HeaderContentType, metadata.ContentType)
	respHeader.Add(storage.HeaderContentRange, metadata.ContentRange)
	respHeader.Add(storage.HeaderETag, metadata.ETag)
	respHeader.Add(storage.HeaderLastModified, metadata.LastModified)

	s.client.On("Get", fmt.Sprintf("%s://%s%s", s.storage.getProtocol(), validHost, validPath), reqHeader).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("response body"))),
			Header:     respHeader,
		}, nil)

	opt := storage.GetPartialObjectRequestOptions{Range: validRange}
	res := s.storage.GetPartialObject(context.TODO(), validPath, &opt)

	assert.Nil(s.T(), res.Error())
	assert.Equal(s.T(), http.StatusOK, res.Status())
	assert.Equal(s.T(), []byte("response body"), res.Data())
	assert.Equal(s.T(), &metadata, res.Metadata())
}

func (s *StorageTestSuite) TestStorage_getURL() {
	cases := []struct {
		secureProtocol bool
		cloudfrontHost string
		path           string
		expected       string
	}{
		{secureProtocol: true, cloudfrontHost: "www.darkroom.com", path: "images/v1", expected: "https://www.darkroom.com/images/v1"},
		{secureProtocol: true, cloudfrontHost: "www.darkroom.com", path: "/images/v1", expected: "https://www.darkroom.com/images/v1"},
		{secureProtocol: true, cloudfrontHost: "www.darkroom.com/", path: "images/v1", expected: "https://www.darkroom.com/images/v1"},
		{secureProtocol: true, cloudfrontHost: "www.darkroom.com/", path: "/images/v1", expected: "https://www.darkroom.com/images/v1"},
	}
	for _, c := range cases {
		s.storage.cloudfrontHost = c.cloudfrontHost
		s.storage.secureProtocol = c.secureProtocol
		assert.Equal(s.T(), c.expected, s.storage.getURL(c.path))
	}
}

// Mocks
type mockClient struct {
	mock.Mock
}

func (m *mockClient) Get(url string, headers http.Header) (*http.Response, error) {
	args := m.Called(url, headers)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *mockClient) Post(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	return nil, nil
}

func (m *mockClient) Put(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	return nil, nil
}

func (m *mockClient) Patch(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	return nil, nil
}

func (m *mockClient) Delete(url string, headers http.Header) (*http.Response, error) {
	return nil, nil
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return nil, nil
}
