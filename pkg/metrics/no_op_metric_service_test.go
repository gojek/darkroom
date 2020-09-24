package metrics

import (
	"testing"
	"time"
)

func TestNoOpMetricService(t *testing.T) {
	ms := NoOpMetricService{}
	ms.CountImageHandlerErrors("handler_error")
	ms.TrackDuration("error", time.Now(), []byte(nil))
}
