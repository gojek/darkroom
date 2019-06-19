package metrics

import (
	"fmt"
	"github.com/cactus/go-statsd-client/statsd"
	"***REMOVED***/darkroom/core/pkg/logger"
	"time"
)

// https://github.com/etsy/statsd/blob/master/docs/metric_types.md#multi-metric-packets
const (
	WANStatsdFlushBytes     = 512
	LANStatsdFlushBytes     = 1432
	GigabitStatsdFlushBytes = 8932
)

var instance *statsdClient

type statsdClient struct {
	client        statsd.Statter
	collectorName string
	sampleRate    float32
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

// InitializeStatsdCollector will start publishing metrics in the form {config.Prefix}.{name}.{updateOption.Name}
func InitializeStatsdCollector(config *StatsdCollectorConfig, name string) error {
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
	instance = &statsdClient{client: c, collectorName: name, sampleRate: sampleRate}
	return nil
}

var formatter = func(on string) string {
	return fmt.Sprintf("%s.%s", instance.collectorName, on)
}

func Update(updateOption UpdateOption) {
	if instance == nil {
		return
	}
	var err error
	switch updateOption.Type {
	case Duration:
		err = instance.client.TimingDuration(formatter(updateOption.Name), updateOption.Duration, instance.sampleRate)
		break
	case Guage:
		err = instance.client.Gauge(formatter(updateOption.Name), int64(updateOption.NumValue), instance.sampleRate)
	}
	logger.Error(err)
}
