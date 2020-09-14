package metrics

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type MockMetricService struct {
	mock.Mock
}

func (m *MockMetricService) TrackDecodeDuration(start time.Time, ImageData []byte) {
}

func (m *MockMetricService) TrackEncodeDuration(start time.Time, ImageData []byte) {
}

func (m *MockMetricService) TrackCropDuration(start time.Time, ImageData []byte) {
}

func (m *MockMetricService) TrackScaleDuration(start time.Time, ImageData []byte) {
}

func (m *MockMetricService) TrackResizeDuration(start time.Time, ImageData []byte) {
}

func (m *MockMetricService) TrackGrayScaleDuration(start time.Time, ImageData []byte) {
}

func (m *MockMetricService) TrackBlurDuration(start time.Time, ImageData []byte) {
}

func (m *MockMetricService) TrackFixOrientationDuration(start time.Time, ImageData []byte) {
}

func (m *MockMetricService) TrackFlipDuration(start time.Time, ImageData []byte) {
}

func (m *MockMetricService) TrackRotateDuration(start time.Time, ImageData []byte) {
}

func (m *MockMetricService) CountStorageGetErrors() {
	m.Called()
}

func (m *MockMetricService) CountProcessorErrors() {
	m.Called()
}

