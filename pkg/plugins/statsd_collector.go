package plugins

import (
	"github.com/cactus/go-statsd-client/statsd"
	"***REMOVED***/darkroom/core/pkg/plugins/metrics"
	"time"
)

// https://github.com/etsy/statsd/blob/master/docs/metric_types.md#multi-metric-packets
const (
	WANStatsdFlushBytes     = 512
	LANStatsdFlushBytes     = 1432
	GigabitStatsdFlushBytes = 8932
)

type StatsdCollector struct {
	client             statsd.Statter
	processingDuration string
	downloadDuration   string
	totalDuration      string
	readDuration       string
	writeDuration      string
	cropDuration       string
	resizeDuration     string
	monoDuration       string
	sampleRate         float32
}

type StatsdCollectorClient struct {
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

func InitializeStatsdCollector(config *StatsdCollectorConfig) *StatsdCollectorClient {
	flushBytes := config.FlushBytes
	if flushBytes == 0 {
		flushBytes = LANStatsdFlushBytes
	}

	sampleRate := config.SampleRate
	if sampleRate == 0 {
		sampleRate = 1
	}

	c, _ := statsd.NewBufferedClient(config.StatsdAddr, config.Prefix, 1*time.Second, flushBytes)
	// TODO Add logger for error
	return &StatsdCollectorClient{client: c, sampleRate: sampleRate}
}

// NewStatsdMetricCollector creates a collector with specific name. The
// prefix given to these stats will be {config.Prefix}.{name}.{metric}.
func (s *StatsdCollectorClient) NewStatsdMetricCollector(name string) metrics.MetricCollector {
	return &StatsdCollector{
		client:             s.client,
		processingDuration: name + ".processingDuration",
		downloadDuration:   name + ".downloadDuration",
		totalDuration:      name + ".totalDuration",
		readDuration:       name + ".readDuration",
		writeDuration:      name + ".writeDuration",
		cropDuration:       name + ".cropDuration",
		resizeDuration:     name + ".resizeDuration",
		monoDuration:       name + ".monoDuration",
		sampleRate:         s.sampleRate,
	}
}

func (sc *StatsdCollector) Update(metrics.MetricResult) {
	// TODO ("implement me")
}

func (sc *StatsdCollector) Reset() {
	// TODO ("implement me")
}
