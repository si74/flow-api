package flowd

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	requests        *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewMetrics(reg *prometheus.Registry) *Metrics {
	metrics := &Metrics{
		requests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "flowd",
				Name:      "http_requests_total",
				Help:      "Flow http requests total",
			},
			[]string{"type", "method", "status_code"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "flowd",
				Name:      "http_request_duration",
				Help:      "Flow http request duration",
				// Note(sneha): Not defining histogram buckets but
				// using default values here.
			},
			[]string{"type", "method", "status_code"},
		),
	}

	reg.MustRegister(metrics.requests)
	reg.MustRegister(metrics.requestDuration)

	return metrics
}
