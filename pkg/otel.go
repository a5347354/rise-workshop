package pkg

import (
	"github.com/spf13/viper"
	jaeger_propagator "go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func NewTracerProvider() trace.TracerProvider {
	exp, err := jaeger.New(jaeger.WithAgentEndpoint())
	if err != nil {
		panic(err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(viper.GetString("service.id")),
			attribute.String("environment", func() string {
				if viper.GetBool("debug") {
					return "test"
				}
				return "production"
			}()),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(jaeger_propagator.Jaeger{})

	return otel.GetTracerProvider()
}
