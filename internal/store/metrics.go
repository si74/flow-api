package store

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
}

func NewMetrics(p *prometheus.Registry) *Metrics {
	return &Metrics{}
}
