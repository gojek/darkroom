package runtime

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var instance *prometheus.Registry
var once sync.Once

func PrometheusRegistry() *prometheus.Registry {
	once.Do(func() {
		instance = prometheus.NewRegistry()
	})
	return instance
}
