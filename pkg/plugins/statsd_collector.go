package plugins


import (
	"github.com/cactus/go-statsd-client/statsd"
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
}

func InitializeStatsdCollector(config *StatsdCollectorConfig) (*StatsdCollectorClient, error) {
	return &StatsdCollectorClient{}, nil
}