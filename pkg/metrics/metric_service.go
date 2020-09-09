package metrics

import (
	"github.com/gorilla/mux"
	"time"
)

type MetricService interface {
	TrackDecodeDuration(start time.Time, ImageData []byte)

	TrackEncodeDuration(start time.Time, ImageData []byte)

	TrackCropDuration(start time.Time, ImageData []byte)

	TrackScaleDuration(start time.Time, ImageData []byte)

	TrackResizeDuration(start time.Time, ImageData []byte)

	TrackGrayScaleDuration(start time.Time, ImageData []byte)

	TrackBlurDuration(start time.Time, ImageData []byte)

	TrackFixOrientationDuration(start time.Time, ImageData []byte)

	TrackFlipDuration(start time.Time, ImageData []byte)

	TrackRotateDuration(start time.Time, ImageData []byte)

	AddMetricsEndPoint(metricsPath string, router *mux.Router)

	CountStorageGetErrors()

	CountProcessorErrors()
}