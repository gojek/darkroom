package metrics

import (
	"time"
)

type MetricService interface {

	TrackDuration(imageProcess string, start time.Time, ImageData []byte)

    CountImageHandlerErrors(kind string)

}