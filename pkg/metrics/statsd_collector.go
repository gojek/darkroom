package metrics

import (
	"fmt"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/gojek/darkroom/pkg/logger"
	"strings"
	"time"
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
func InitializeStatsdCollector(config *StatsdCollectorConfig) error {
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
		c = &statsd.NoopClient{}
	}
	instance = &statsdClient{client: c, sampleRate: sampleRate}
	return nil
}

func formatter(updateOption UpdateOption) string {
	scope := strings.Trim(updateOption.Scope, ".")
	if updateOption.Scope == "" {
		scope = DefaultScope
	}
	return fmt.Sprintf("%s.%s", scope, strings.Trim(updateOption.Name, "."))
}

// Update takes an UpdateOption and pushes the metrics to the statd client if initialised
func Update(updateOption UpdateOption) {
	if instance == nil {
		return
	}
	var err error
	switch updateOption.Type {
	case Duration:
		err = instance.client.TimingDuration(formatter(updateOption), updateOption.Duration, instance.sampleRate)
	case Gauge:
		err = instance.client.Gauge(formatter(updateOption), int64(updateOption.NumValue), instance.sampleRate)
	case Count:
		err = instance.client.Inc(formatter(updateOption), 1, instance.sampleRate)
	}
	if err != nil {
		logger.Errorf("metrics.Update got an error: %s", err)
	}
}
