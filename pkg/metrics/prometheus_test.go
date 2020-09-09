package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPrometheusMetrics(t *testing.T) {
	tests := []struct {
		name string
		addMetrics func(MetricService)
		expMetrics []string
		expCode int
	}{
		{
			name: "Measuring duration metrics should expose metrics on prometheus endpoint.",
			addMetrics: func(s MetricService) {
				now := time.Now()
				imageData, err := ioutil.ReadFile("../processor/native/_testdata/test.png")
				if err != nil {
					panic(err)
				}
				s.TrackCropDuration(now.Add(-6*time.Second), imageData)
			},
			expMetrics: []string{
				`image_crop_duration_bucket{image_type="<=128KB.png",le="0.005"} 0`,
				`image_crop_duration_bucket{image_type="<=128KB.png",le="0.01"} 0`,
				`image_crop_duration_bucket{image_type="<=128KB.png",le="0.025"} 0`,
				`image_crop_duration_bucket{image_type="<=128KB.png",le="0.05"} 0`,
				`image_crop_duration_bucket{image_type="<=128KB.png",le="0.1"} 0`,
				`image_crop_duration_bucket{image_type="<=128KB.png",le="0.25"} 0`,
				`image_crop_duration_bucket{image_type="<=128KB.png",le="0.5"} 0`,
				`image_crop_duration_bucket{image_type="<=128KB.png",le="1"} 0`,
				`image_crop_duration_bucket{image_type="<=128KB.png",le="2.5"} 0`,
				`image_crop_duration_bucket{image_type="<=128KB.png",le="5"} 0`,
				`image_crop_duration_bucket{image_type="<=128KB.png",le="10"} 1`,
				`image_crop_duration_bucket{image_type="<=128KB.png",le="+Inf"} 1`,
				`image_crop_duration_count{image_type="<=128KB.png"} 1`,
			},
			expCode: 200,
		},
		{
			name: "Measuring storage and processor errors should expose metrics on prometheus endpoint.",
			addMetrics: func(s MetricService) {
				s.CountProcessorErrors()
				s.CountStorageGetErrors()
				s.CountStorageGetErrors()
			},
			expMetrics: []string{
				`processor_errors_total 1`,
				`storage_get_errors_total 2`,
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
