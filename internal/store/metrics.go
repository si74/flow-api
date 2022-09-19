package store

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	// FlowKeyCount is the current # of flows in the flowstore
	flows prometheus.Gauge
	// TODO(sneha): add gauge indicating total # of data flow points in a flow map
}

func NewMetrics(reg *prometheus.Registry) *Metrics {
	metrics := &Metrics{
		flows: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "flowd",
				Name:      "flowstore_size",
				Help:      "Number of total flow datapoints in the flowstore",
			},
		),
	}

	reg.MustRegister(metrics.flows)

	return metrics
}
