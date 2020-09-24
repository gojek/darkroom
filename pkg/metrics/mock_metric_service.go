package metrics

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type MockMetricService struct {
	mock.Mock
}

func (m *MockMetricService) TrackDuration(imageProcess string, start time.Time, ImageData []byte) {
	m.Called(imageProcess, start, ImageData)
}

func (m *MockMetricService) CountImageHandlerErrors(kind string) {
	m.Called(kind)
}
