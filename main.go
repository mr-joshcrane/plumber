package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsxrayexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func performAction(action string) {
	out, err := exec.Command("/bin/bash", action).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}

func main() {
	action := os.Args[1]

	ctx := context.Background()
	buildInfo := component.DefaultBuildInfo()
	logger := zap.NewExample()
	set := component.ExporterCreateSettings{ logger, buildInfo,	}
	config := awsxrayexporter.NewFactory().CreateDefaultConfig()

	tExporter, err := awsxrayexporter.NewFactory().CreateTracesExporter(
		ctx, set, config,
	)
	if err != nil {
		log.Fatal(err)
	}
	
	tp := sdktrace.NewTracerProvider()


	// Handle this error in a sensible manner where possible
	defer func() { _ = tp.Shutdown(ctx) }()


	// Handle this error in a sensible manner where possible

	otel.SetTracerProvider(tp)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)


	tracer := otel.Tracer("Pipeline")
	ctx = baggage.ContextWithValues(ctx)

	func(ctx context.Context) {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "Pipeline Operation")
		defer span.End()
		performAction(action)

		span.AddEvent("Performing Action", trace.WithAttributes(attribute.String("command", action)))
		traces := pdata.NewTraces()
		tExporter.ConsumeTraces(ctx, traces)
	}(ctx)
}
