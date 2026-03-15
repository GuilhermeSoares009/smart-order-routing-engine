package observability

import (
    "context"
    "time"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
    "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/metric"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

const serviceName = "smart-order-routing-engine"

func Init(ctx context.Context) (func(context.Context) error, func(context.Context) error, error) {
    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceName(serviceName),
        ),
    )
    if err != nil {
        return nil, nil, err
    }

    traceExporter, err := stdouttrace.New(
        stdouttrace.WithPrettyPrint(),
        stdouttrace.WithWriter(nil),
    )
    if err != nil {
        return nil, nil, err
    }

    tracerProvider := trace.NewTracerProvider(
        trace.WithBatcher(traceExporter),
        trace.WithResource(res),
    )
    otel.SetTracerProvider(tracerProvider)

    metricExporter, err := stdoutmetric.New()
    if err != nil {
        return nil, nil, err
    }

    meterProvider := metric.NewMeterProvider(
        metric.WithReader(metric.NewPeriodicReader(metricExporter, metric.WithInterval(15*time.Second))),
        metric.WithResource(res),
    )
    otel.SetMeterProvider(meterProvider)

    otel.SetTextMapPropagator(propagation.TraceContext{})

    return tracerProvider.Shutdown, meterProvider.Shutdown, nil
}
