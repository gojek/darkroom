package metrics

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/afex/hystrix-go/plugins"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/gojek/darkroom/pkg/logger"
)

// https://github.com/etsy/statsd/blob/master/docs/metric_types.md#multi-metric-packets
const (
	WANStatsdFlushBytes     = 512
	LANStatsdFlushBytes     = 1432
	GigabitStatsdFlushBytes = 8932
	DefaultScope            = "default"
)

var instance *statsdClient

type statsdClient struct {
	client     statsd.Statter
	sampleRate float32
}

// StatsdCollectorConfig provides configuration that the Statsd client will need.
type StatsdCollectorConfig struct {
	// StatsdAddr is the tcp address of the Statsd server
	StatsdAddr string
	// Prefix is the prefix that will be prepended to all metrics sent from this collector.
	Prefix string
	// StatsdSampleRate sets statsd sampling. If 0, defaults to 1.0. (no sampling)
	SampleRate float32
	// FlushBytes sets message size for statsd packets. If 0, defaults to LANFlushSize.
	FlushBytes int
}

// InitializeStatsdCollector will start publishing metrics in the form {config.Prefix}.{updateOption.Scope|default}.{updateOption.Name}
func InitializeStatsdCollector(config *StatsdCollectorConfig) (MetricService, error) {
	flushBytes := config.FlushBytes
	if flushBytes == 0 {
		flushBytes = LANStatsdFlushBytes
	}

	sampleRate := config.SampleRate
	if sampleRate == 0 {
		sampleRate = 1
	}

	c, err := statsd.NewBufferedClient(config.StatsdAddr, config.Prefix, 1*time.Second, flushBytes)
	if err != nil {
		// TODO Add logger for error
		logger.Errorf("failed to initialize statsd collector with error: %s", err.Error())
		return nil, err
	}
	instance = &statsdClient{client: c, sampleRate: sampleRate}
	return instance, nil
}

func RegisterHystrixMetrics(config *StatsdCollectorConfig, prefix string) error {
	c, err := plugins.InitializeStatsdCollector(&plugins.StatsdCollectorConfig{
		StatsdAddr: config.StatsdAddr,
		Prefix:     prefix,
	})
	if err != nil {
		logger.Errorf("failed to initialize statsd collector for hystrix metrics with error: %s", err.Error())
		return err
	}
	metricCollector.Registry.Register(c.NewStatsdCollector)
	return nil
}

func (s statsdClient) TrackDuration(imageProcess string, start time.Time, ImageData []byte) {
	metricTag := s.getMetricTag(imageProcess, ImageData)
	err := s.client.TimingDuration(metricTag, time.Since(start), s.sampleRate)
	if err != nil {
		logger.Errorf("MetricService.TrackDuration got an error: %s", err)
	}
}

func (s statsdClient) CountImageHandlerErrors(kind string) {
	err := s.client.Inc(kind, 1, s.sampleRate)
	if err != nil {
		logger.Errorf("MetricService.CountImageHandlerErrors got an error: %s", err)
	}
}

func (s statsdClient) getMetricTag(imageProcess string, ImageData []byte) string {
	ext := strings.Split(http.DetectContentType(ImageData), "/")[1]
	tag := fmt.Sprintf("%s.%s.%s", imageProcess, GetImageSizeCluster(ImageData), ext)
	return tag
}

