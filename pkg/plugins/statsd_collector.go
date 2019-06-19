package plugins

import (
	"github.com/cactus/go-statsd-client/statsd"
	"time"
)

// https://github.com/etsy/statsd/blob/master/docs/metric_types.md#multi-metric-packets
const (
	WANStatsdFlushBytes     = 512
	LANStatsdFlushBytes     = 1432
	GigabitStatsdFlushBytes = 8932
)

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
