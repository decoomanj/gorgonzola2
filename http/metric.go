package gorgonzola

import (
	"github.com/prometheus/client_golang/prometheus"
)

//
type Metrics struct {
	*prometheus.Registry
}

func NewMetrics() *Metrics {
	return &Metrics{
		prometheus.NewRegistry(),
	}
}
