package metric

import (
	"fmt"

	metricsdk "go.opentelemetry.io/otel/metric"
)

//go:generate mockgen -destination=mocks/mock_metric.go -source=metric.go Metric
type Metric interface {
	// Add metrics implementations here
}
type metric struct {
	// Add add metric dependencies here
}

func NewMetric(meter metricsdk.Meter) (*metric, error) {
	if meter == nil {
		return nil, fmt.Errorf("meter should not be nil")
	}

	return &metric{}, nil
}
