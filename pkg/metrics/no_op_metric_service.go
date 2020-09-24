package metrics

import (
	"time"
)

type NoOpMetricService struct{}

func (NoOpMetricService) TrackDuration(string, time.Time, []byte) {
}

func (NoOpMetricService) CountImageHandlerErrors(string) {
}
