package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strings"
	"time"
)

type prometheusService struct {
	decodeDuration *prometheus.HistogramVec
	encodeDuration *prometheus.HistogramVec
	cropDuration *prometheus.HistogramVec
	scaleDuration *prometheus.HistogramVec
	resizeDuration *prometheus.HistogramVec
	grayscaleDuration *prometheus.HistogramVec
	blurDuration *prometheus.HistogramVec
	fixOrientationDuration *prometheus.HistogramVec
	flipDuration *prometheus.HistogramVec
	rotateDuration *prometheus.HistogramVec
	storageGetErrors prometheus.Counter
	processorErrors prometheus.Counter
	reg *prometheus.Registry
}

func NewPrometheus(reg *prometheus.Registry) MetricService {
	p := &prometheusService{
		decodeDuration:   prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "image_decode_duration",
				Help: "Time taken to decode requested image",
			}, []string{"image_type"}),
		encodeDuration:   prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "image_encode_duration",
				Help: "Time taken to encode data back to image",
			}, []string{"image_type"}),
		cropDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "image_crop_duration",
				Help: "Time taken to apply cropping to image data",
			}, []string{"image_type"}),
		scaleDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "image_scale_duration",
				Help: "Time taken to apply scaling to image data",
			}, []string{"image_type"}),
		resizeDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "image_resize_duration",
				Help: "Time taken to apply resizing to image data",
			}, []string{"image_type"}),
		grayscaleDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "image_grayscale_duration",
				Help: "Time taken to apply grayscale filter to image data",
			}, []string{"image_type"}),
		blurDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "image_blur_duration",
				Help: "Time taken to apply blur to image data",
			}, []string{"image_type"}),
		fixOrientationDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "image_orientation_duration",
				Help: "Time to taken to apply orientation to image data",
			}, []string{"image_type"}),
		flipDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "image_flip_duration",
				Help: "Time taken to apply flipping to image data",
			}, []string{"image_type"}),
		rotateDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "image_rotate_duration",
				Help: "Time taken to apply rotation to image data",
			}, []string{"image_type"}),
		storageGetErrors: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "storage_get_errors_total",
				Help: "The total number of storage get errors",
			}),
		processorErrors: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "processor_errors_total",
				Help: "The total number of storage get errors",
			}),

		reg: reg,
	}
	p.registerMetrics()

	return p
}

func (p prometheusService) registerMetrics() {
	p.reg.MustRegister(
		p.decodeDuration,
		p.encodeDuration,
		p.cropDuration,
		p.scaleDuration,
		p.resizeDuration,
		p.grayscaleDuration,
		p.blurDuration,
		p.fixOrientationDuration,
		p.flipDuration,
		p.rotateDuration,
		p.storageGetErrors,
		p.processorErrors,
	)
}

func (p prometheusService) TrackDecodeDuration(start time.Time, ImageData []byte) {
	labelValue := p.getLabelValue(ImageData)
	p.decodeDuration.WithLabelValues(labelValue).Observe(time.Since(start).Seconds())
}

func (p prometheusService) TrackEncodeDuration(start time.Time, ImageData []byte) {
	labelValue := p.getLabelValue(ImageData)
	p.encodeDuration.WithLabelValues(labelValue).Observe(time.Since(start).Seconds())
}

func (p prometheusService) TrackCropDuration(start time.Time, ImageData []byte) {
	labelValue := p.getLabelValue(ImageData)
	p.cropDuration.WithLabelValues(labelValue).Observe(time.Since(start).Seconds())
}

func (p prometheusService) TrackScaleDuration(start time.Time, ImageData []byte) {
	labelValue := p.getLabelValue(ImageData)
	p.scaleDuration.WithLabelValues(labelValue).Observe(time.Since(start).Seconds())
}

func (p prometheusService) TrackResizeDuration(start time.Time, ImageData []byte) {
	labelValue := p.getLabelValue(ImageData)
	p.resizeDuration.WithLabelValues(labelValue).Observe(time.Since(start).Seconds())
}

func (p prometheusService) TrackGrayScaleDuration(start time.Time, ImageData []byte) {
	labelValue := p.getLabelValue(ImageData)
	p.grayscaleDuration.WithLabelValues(labelValue).Observe(time.Since(start).Seconds())
}

func (p prometheusService) TrackBlurDuration(start time.Time, ImageData []byte) {
	labelValue := p.getLabelValue(ImageData)
	p.blurDuration.WithLabelValues(labelValue).Observe(time.Since(start).Seconds())
}

func (p prometheusService) TrackFixOrientationDuration(start time.Time, ImageData []byte) {
	labelValue := p.getLabelValue(ImageData)
	p.fixOrientationDuration.WithLabelValues(labelValue).Observe(time.Since(start).Seconds())
}

func (p prometheusService) TrackFlipDuration(start time.Time, ImageData []byte) {
	labelValue := p.getLabelValue(ImageData)
	p.flipDuration.WithLabelValues(labelValue).Observe(time.Since(start).Seconds())
}

func (p prometheusService) TrackRotateDuration(start time.Time, ImageData []byte) {
	labelValue := p.getLabelValue(ImageData)
	p.rotateDuration.WithLabelValues(labelValue).Observe(time.Since(start).Seconds())
}

func (p prometheusService) CountStorageGetErrors() {
	p.storageGetErrors.Inc()
}

func (p prometheusService) CountProcessorErrors() {
	p.processorErrors.Inc()
}

func (p prometheusService) getLabelValue(ImageData []byte) string {
	ext := strings.Split(http.DetectContentType(ImageData), "/")[1]
	labelValue := fmt.Sprintf("%s.%s", GetImageSizeCluster(ImageData), ext)
	return labelValue
}


