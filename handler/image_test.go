package handler

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"***REMOVED***/darkroom/server/service"
	"***REMOVED***/darkroom/server/storage"
	"testing"
)

type ImageHandlerTestSuite struct {
	suite.Suite
	deps        *service.Dependencies
	storage     *storage.MockBaseStorage
	manipulator *mockManipulator
}

func TestImageHandlerSuite(t *testing.T) {
	suite.Run(t, new(ImageHandlerTestSuite))
}

func (s *ImageHandlerTestSuite) SetupTest() {
	s.storage = &storage.MockBaseStorage{}
	s.manipulator = &mockManipulator{}
	s.deps = &service.Dependencies{Storage: s.storage, Manipulator: s.manipulator}
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

	ImageHandler(s.deps).ServeHTTP(rr, r)

	assert.Equal(s.T(), "", rr.Body.String())
	assert.Equal(s.T(), http.StatusUnprocessableEntity, rr.Code)
}

func (s *ImageHandlerTestSuite) TestImageHandlerWithQueryParameters() {
	r, _ := http.NewRequest(http.MethodGet, "/image-valid?w=100&h=100", nil)
	rr := httptest.NewRecorder()

	params := make(map[string]string)
	params["w"] = "100"
	params["h"] = "100"
	s.storage.On("Get", mock.Anything, "/image-valid").Return([]byte("validData"), http.StatusOK, nil)
	s.manipulator.On("Process", mock.Anything, []byte("validData"), params).Return([]byte("processedData"), nil)

	ImageHandler(s.deps).ServeHTTP(rr, r)

	assert.Equal(s.T(), "processedData", rr.Body.String())
	assert.Equal(s.T(), http.StatusOK, rr.Code)
}

func (s *ImageHandlerTestSuite) TestImageHandlerWithQueryParametersAndProcessingError() {
	r, _ := http.NewRequest(http.MethodGet, "/image-valid?w=100&h=100", nil)
	rr := httptest.NewRecorder()

	params := make(map[string]string)
	params["w"] = "100"
	params["h"] = "100"
	s.storage.On("Get", mock.Anything, "/image-valid").Return([]byte("validData"), http.StatusOK, nil)
	s.manipulator.On("Process", mock.Anything, []byte("validData"), params).Return([]byte(nil), errors.New("error"))

	ImageHandler(s.deps).ServeHTTP(rr, r)

	assert.Equal(s.T(), "", rr.Body.String())
	assert.Equal(s.T(), http.StatusUnprocessableEntity, rr.Code)
}

type mockManipulator struct {
	mock.Mock
}

func (m *mockManipulator) Process(ctx context.Context, data []byte, params map[string]string) ([]byte, error) {
	args := m.Called(ctx, data, params)
	if args.Get(1) == nil {
		return args.Get(0).([]byte), nil
	}
	return args.Get(0).([]byte), args.Get(1).(error)
}
