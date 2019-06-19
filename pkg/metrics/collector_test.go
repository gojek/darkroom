package metrics

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestMetricCollectorRegistry_Register(t *testing.T) {
	sc := InitializeStatsdCollector(&StatsdCollectorConfig{})
	mcr := metricCollectorRegistry{lock: &sync.RWMutex{}}

	mcr.Register(sc.NewStatsdMetricCollector)
	assert.Equal(t, 1, len(mcr.registry))
}
