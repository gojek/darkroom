package webfolder

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
	validBaseURL = "https://example.com/path/to/images"
	validPath    = "/path/to/valid-file"
	invalidPath  = "/path/to/invalid-file"
	validRange   = "bytes=100-200"
)

type StorageTestSuite struct {
	suite.Suite
	storage Storage
	client  *mockClient
}

func (s *StorageTestSuite) SetupTest() {
	s.client = &mockClient{}
	s.storage = *NewStorage(
		WithBaseURL(validBaseURL),
		WithHeimdallClient(s.client),
	)
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (s *StorageTestSuite) TestNewStorage() {
	assert.NotNil(s.T(), s.storage)
}

func (s *StorageTestSuite) TestStorage_GetNotFound() {
	s.client.On("Get", fmt.Sprintf("%s%s", validBaseURL, invalidPath), http.Header(nil)).
		Return(&http.Response{StatusCode: http.StatusNotFound}, errors.New("not found"))

	res := s.storage.Get(context.TODO(), invalidPath)

	assert.NotNil(s.T(), res.Error())
	assert.Equal(s.T(), http.StatusNotFound, res.Status())
	assert.Nil(s.T(), res.Data())
}

func (s *StorageTestSuite) TestStorage_GetNoResponse() {
	s.client.On("Get", fmt.Sprintf("%s%s", validBaseURL, invalidPath), http.Header(nil)).
		Return(nil, errors.New("response body read failure"))

	res := s.storage.Get(context.TODO(), invalidPath)

	assert.NotNil(s.T(), res.Error())
	assert.Equal(s.T(), http.StatusUnprocessableEntity, res.Status())
	assert.Nil(s.T(), res.Data())
}

func (s *StorageTestSuite) TestStorage_GetSuccessResponse() {
	s.client.On("Get", fmt.Sprintf("%s%s", validBaseURL, validPath), http.Header(nil)).
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
	s.client.On("Get", fmt.Sprintf("%s%s", validBaseURL, validPath), http.Header(nil)).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("response body"))),
		}, nil)

	opt := storage.GetPartiallyRequestOptions{Range: validRange}
	res := s.storage.GetPartially(context.TODO(), validPath, &opt)

	assert.Nil(s.T(), res.Error())
	assert.Equal(s.T(), http.StatusOK, res.Status())
	assert.Equal(s.T(), []byte("response body"), res.Data())
	assert.Nil(s.T(), res.Metadata())
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
