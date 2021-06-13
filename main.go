package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsxrayexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsutils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
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
	// exporter, err := stdout.NewExporter(
	// 	stdout.WithPrettyPrint(),
	// )

	// if err != nil {
	// 	log.Fatalf("failed to initialize stdout export pipeline: %v", err)

	// }
	ctx := context.Background()
	set := component.ExporterCreateSettings{}
	configSettings := config.NewExporterSettings(config.NewIDWithName("Exporters", "pipeline"))
	config := config

	exporter, err := awsxrayexporter.NewFactory().CreateTracesExporter(
		ctx, set, config,
	)


	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp))

	// Handle this error in a sensible manner where possible
	defer func() { _ = tp.Shutdown(ctx) }()

	pusher := controller.New(
		processor.New(
			simple.NewWithExactDistribution(),
			exporter,
		),
		controller.WithExporter(exporter),
		controller.WithCollectPeriod(5*time.Second),
	)

	err = pusher.Start(ctx)
	if err != nil {
		log.Fatalf("failed to initialize metric controller: %v", err)
	}

	// Handle this error in a sensible manner where possible
	defer func() { _ = pusher.Stop(ctx) }()

	otel.SetTracerProvider(tp)
	global.SetMeterProvider(pusher.MeterProvider())
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
	}(ctx)
}
