package pkg

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"log"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

func NewTracerProvider() trace.TracerProvider {
	ctx := context.Background()
	projectID := viper.GetString("gcp.project.id")
	exp, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		log.Fatalf("texporter.New: %v", err)
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
	defer tp.ForceFlush(ctx) // flushes any pending spans

	otel.SetTracerProvider(tp)
	//otel.SetTextMapPropagator(jaeger_propagator.Jaeger{})

	return otel.GetTracerProvider()
}
