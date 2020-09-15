package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/gojek/darkroom/pkg/metrics"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojek/darkroom/pkg/config"
	"github.com/gojek/darkroom/pkg/service"
	"github.com/gojek/darkroom/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ImageHandlerTestSuite struct {
	suite.Suite
	deps        *service.Dependencies
	storage     *mockStorage
	manipulator *service.MockManipulator
	mockMetricService *metrics.MockMetricService
}

func TestImageHandlerSuite(t *testing.T) {
	suite.Run(t, new(ImageHandlerTestSuite))
}

func (s *ImageHandlerTestSuite) SetupTest() {
	s.storage = &mockStorage{}
	s.manipulator = &service.MockManipulator{}
	s.mockMetricService = &metrics.MockMetricService{}
	s.deps = &service.Dependencies{Storage: s.storage, Manipulator: s.manipulator,
		MetricService: s.mockMetricService}
}

func (s *ImageHandlerTestSuite) TestImageHandler() {
	r, _ := http.NewRequest(http.MethodGet, "/image-valid", nil)
	rr := httptest.NewRecorder()

	s.storage.On("Get", mock.Anything, "/image-valid").Return([]byte("validData"), http.StatusOK, nil)

	ImageHandler(s.deps).ServeHTTP(rr, r)

	assert.Equal(s.T(), "validData", rr.Body.String())
	assert.Equal(s.T(), http.StatusOK, rr.Code)
}

func (s *ImageHandlerTestSuite) TestImageHandlerWithStorageGetError() {
	r, _ := http.NewRequest(http.MethodGet, "/image-invalid", nil)
	rr := httptest.NewRecorder()

	s.storage.On("Get", mock.Anything, "/image-invalid").Return([]byte(nil), http.StatusUnprocessableEntity, errors.New("error"))
	s.mockMetricService.On("CountImageHandlerErrors")

	ImageHandler(s.deps).ServeHTTP(rr, r)

	s.mockMetricService.AssertCalled(s.T(), "CountImageHandlerErrors")
	assert.Equal(s.T(), "", rr.Body.String())
	assert.Equal(s.T(), http.StatusUnprocessableEntity, rr.Code)
}

func (s *ImageHandlerTestSuite) TestImageHandlerWithQueryParameters() {
	r, _ := http.NewRequest(http.MethodGet, "/image-valid?w=100&h=100", nil)
	rr := httptest.NewRecorder()
	maxAge := config.CacheTime()
	processedData := []byte("processedData")

	params := make(map[string]string)
	params["w"] = "100"
	params["h"] = "100"
	s.storage.On("Get", mock.Anything, "/image-valid").Return([]byte("validData"), http.StatusOK, nil)
	s.manipulator.On("Process", mock.AnythingOfType("service.processSpec")).Return(processedData, nil)

	ImageHandler(s.deps).ServeHTTP(rr, r)

	assert.Equal(s.T(), "processedData", rr.Body.String())
	assert.Equal(s.T(), http.StatusOK, rr.Code)
	assert.Equal(s.T(), fmt.Sprintf("%d", len(processedData)), rr.Header().Get(ContentLengthHeader))
	assert.Equal(s.T(), fmt.Sprintf("public,max-age=%d", maxAge), rr.Header().Get(CacheControlHeader))
	assert.Equal(s.T(), "Accept", rr.Header().Get(VaryHeader))
}

func (s *ImageHandlerTestSuite) TestImageHandlerWithQueryParametersAndProcessingError() {
	r, _ := http.NewRequest(http.MethodGet, "/image-valid?w=100&h=100", nil)
	rr := httptest.NewRecorder()

	params := make(map[string]string)
	params["w"] = "100"
	params["h"] = "100"
	s.storage.On("Get", mock.Anything, "/image-valid").Return([]byte("validData"), http.StatusOK, nil)
	s.manipulator.On("Process", mock.AnythingOfType("service.processSpec")).Return([]byte(nil), errors.New("error"))
	s.mockMetricService.On("CountImageHandlerErrors")

	ImageHandler(s.deps).ServeHTTP(rr, r)

	s.mockMetricService.AssertCalled(s.T(), "CountImageHandlerErrors")
	assert.Equal(s.T(), "", rr.Body.String())
	assert.Equal(s.T(), http.StatusUnprocessableEntity, rr.Code)
}

type mockStorage struct {
	mock.Mock
}

func (m *mockStorage) Get(ctx context.Context, path string) storage.IResponse {
	args := m.Called(ctx, path)
	return storage.NewResponse(args[0].([]byte), args.Int(1), args.Error(2))
}

func (m *mockStorage) GetPartially(ctx context.Context, path string, opt *storage.GetPartiallyRequestOptions) storage.IResponse {
	args := m.Called(ctx, path, opt)
	return storage.NewResponse(args[0].([]byte), args.Int(1), args.Error(2)).WithMetadata(args[3].(*storage.ResponseMetadata))
}

