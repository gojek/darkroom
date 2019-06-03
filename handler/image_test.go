package handler

import (
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
	storage *storage.MockBaseStorage
}

func TestImageHandlerSuite(t *testing.T) {
	suite.Run(t, new(ImageHandlerTestSuite))
}

func (s *ImageHandlerTestSuite) SetupSuite() {
	s.storage = &storage.MockBaseStorage{}
}

func (s *ImageHandlerTestSuite) TestImageHandler() {
	r, _ := http.NewRequest(http.MethodGet, "/image-valid", nil)
	rr := httptest.NewRecorder()

	s.storage.On("Get", mock.Anything, "/image-valid").Return([]byte("someData"), http.StatusOK, nil)
	deps := service.Dependencies{Storage: s.storage}

	ImageHandler(&deps).ServeHTTP(rr, r)

	assert.Equal(s.T(), "someData", rr.Body.String())
	assert.Equal(s.T(), http.StatusOK, rr.Code)
}

func (s *ImageHandlerTestSuite) TestImageHandlerWithStorageGetError() {
	r, _ := http.NewRequest(http.MethodGet, "/image-invalid", nil)
	rr := httptest.NewRecorder()

	s.storage.On("Get", mock.Anything, "/image-invalid").Return([]byte(nil), http.StatusUnprocessableEntity, errors.New("error"))
	deps := service.Dependencies{Storage: s.storage}

	ImageHandler(&deps).ServeHTTP(rr, r)

	assert.Equal(s.T(), "", rr.Body.String())
	assert.Equal(s.T(), http.StatusUnprocessableEntity, rr.Code)
}
