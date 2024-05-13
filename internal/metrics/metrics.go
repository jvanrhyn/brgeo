package metrics

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// NewMeter creates a new metric.Meter that can create any metric reporter
// you might want to use in your application.
func NewMeter(ctx context.Context) (metric.Meter, error) {
	provider, err := newMeterProvider(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create meter provider: %w", err)
	}
	return provider.Meter("tracetest"), nil
}

// newMeterProvcider initialize the application resource, connects to the
// OpenTelemetry Collector and configures the metric poller that will be used
// to collect the metrics and send them to the OpenTelemetry Collector.
func newMeterProvider(ctx context.Context) (metric.MeterProvider, error) {
	// Interval which the metrics will be reported to the collector
	interval := 10 * time.Second
	resource, err := getResource()
	if err != nil {
		return nil, fmt.Errorf("could not get resource: %w", err)
	}
	collectorExporter, err := getOtelMetricsCollectorExporter(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get collector exporter: %w", err)
	}
	periodicReader := metricsdk.NewPeriodicReader(collectorExporter,
		metricsdk.WithInterval(interval),
	)
	provider := metricsdk.NewMeterProvider(
		metricsdk.WithResource(resource),
		metricsdk.WithReader(periodicReader),
	)
	return provider, nil
}

// getResource creates the resource that describes our application.
//
// You can add any attributes to your resource and all your metrics
// will contain those attributes automatically.
//
// There are some attributes that are very important to be added to the resource:
// 1. hostname: allows you to identify host-specific problems
// 2. version: allows you to pinpoint problems in specific versions
func getResource() (*resource.Resource, error) {
	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("tracetest"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("could not merge resources: %w", err)
	}
	return resource, nil
}

// getOtelMetricsCollectorExporter creates a metric exporter that relies on
// an OpenTelemetry Collector running on "localhost:4317".
func getOtelMetricsCollectorExporter(ctx context.Context) (metricsdk.Exporter, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithCompressor("gzip"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create metric exporter: %w", err)
	}
	return exporter, nil
}
