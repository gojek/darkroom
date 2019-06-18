package cloudfront

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	validHost   = "cloudfront.net"
	validPath   = "/path/to/valid-file"
	invalidPath = "/path/to/invalid-file"
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
