package metrics

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusMetrics(t *testing.T) {
	tests := []struct {
		name       string
		addMetrics func(MetricService)
		expMetrics []string
		expCode    int
	}{
		{
			name: "Measuring duration metrics should expose metrics on prometheus endpoint.",
			addMetrics: func(s MetricService) {
				now := time.Now()
				imageData, err := ioutil.ReadFile("../processor/native/_testdata/test.png")
				if err != nil {
					panic(err)
				}
				s.TrackDuration("cropDuration", now.Add(-6*time.Second), imageData)
			},
			expMetrics: []string{
				`image_process_duration_bucket{image_type="<=128KB.png",process="cropDuration",le="0.005"} 0`,
				`image_process_duration_bucket{image_type="<=128KB.png",process="cropDuration",le="0.01"} 0`,
				`image_process_duration_bucket{image_type="<=128KB.png",process="cropDuration",le="0.025"} 0`,
				`image_process_duration_bucket{image_type="<=128KB.png",process="cropDuration",le="0.05"} 0`,
				`image_process_duration_bucket{image_type="<=128KB.png",process="cropDuration",le="0.1"} 0`,
				`image_process_duration_bucket{image_type="<=128KB.png",process="cropDuration",le="0.25"} 0`,
				`image_process_duration_bucket{image_type="<=128KB.png",process="cropDuration",le="0.5"} 0`,
				`image_process_duration_bucket{image_type="<=128KB.png",process="cropDuration",le="1"} 0`,
				`image_process_duration_bucket{image_type="<=128KB.png",process="cropDuration",le="2.5"} 0`,
				`image_process_duration_bucket{image_type="<=128KB.png",process="cropDuration",le="5"} 0`,
				`image_process_duration_bucket{image_type="<=128KB.png",process="cropDuration",le="10"} 1`,
				`image_process_duration_bucket{image_type="<=128KB.png",process="cropDuration",le="+Inf"} 1`,
				`image_process_duration_count{image_type="<=128KB.png",process="cropDuration"} 1`,
			},
			expCode: 200,
		},
		{
			name: "Measuring storage and processor errors should expose metrics on prometheus endpoint.",
			addMetrics: func(s MetricService) {
				s.CountImageHandlerErrors("storage_get_error")
				s.CountImageHandlerErrors("processor_error")
			},
			expMetrics: []string{
				`image_handler_errors{error_type="storage_get_error"} 1`,
				`image_handler_errors{error_type="processor_error"} 1`,
			},
			expCode: 200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reg := prometheus.NewRegistry()
			m := NewPrometheus(reg)
			test.addMetrics(m)

			h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
			r := httptest.NewRequest("GET", "/metrics", nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			resp := w.Result()

			// Check all metrics are present.
			if assert.Equal(t, test.expCode, resp.StatusCode) {
				body, _ := ioutil.ReadAll(resp.Body)
				for _, expMetric := range test.expMetrics {
					assert.Contains(t, string(body), expMetric, "metric not present on the result of metrics service")
				}
			}
		})
	}
}
