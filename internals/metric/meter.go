package metric

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func NewExporter(ctx context.Context, enable bool, endpoint string) (metricsdk.Exporter, error) {
	if !enable {
		return stdoutmetric.New()
	}

	// since service it's not exposed, it's ok
	// @TODO: use TLS for this kind of connection
	return otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithInsecure(),
	)
}

func NewMeterProvider(exp metricsdk.Exporter, appName string) (*metricsdk.MeterProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(appName),
		),
	)

	if err != nil {
		return nil, err
	}

	return metricsdk.NewMeterProvider(
		metricsdk.WithResource(r),
		metricsdk.WithReader(metricsdk.NewPeriodicReader(exp)),
	), nil

}
