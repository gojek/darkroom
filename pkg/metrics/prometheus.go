package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strings"
	"time"
)

type prometheusService struct {
	imageProcessDuration *prometheus.HistogramVec
	imageHandlerErrorCounter *prometheus.CounterVec
	reg *prometheus.Registry
}

func NewPrometheus(reg *prometheus.Registry) MetricService {
	p := &prometheusService{
		imageProcessDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "image_process_duration",
				Help: "Time taken by each stage to process requested image",
			}, []string{"process", "image_type"}),
		imageHandlerErrorCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "image_handler_errors",
				Help: "The total number of errors for each storage and processor",
			}, []string{"error_type"}),

		reg: reg,
	}
	p.registerMetrics()

	return p
}

func (p prometheusService) registerMetrics() {
	p.reg.MustRegister(
		p.imageProcessDuration,
		p.imageHandlerErrorCounter,
	)
}

func (p prometheusService) TrackDuration(imageProcess string, start time.Time, ImageData []byte) {
	imageType := p.getImageType(ImageData)
	p.imageProcessDuration.WithLabelValues(imageProcess, imageType).Observe(time.Since(start).Seconds())
}

func (p prometheusService) CountImageHandlerErrors(kind string) {
	p.imageHandlerErrorCounter.WithLabelValues(kind).Inc()
}

func (p prometheusService) getImageType(ImageData []byte) string {
	ext := strings.Split(http.DetectContentType(ImageData), "/")[1]
	labelValue := fmt.Sprintf("%s.%s", GetImageSizeCluster(ImageData), ext)
	return labelValue
}


