package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	// "go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	// "go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	apitrace "go.opentelemetry.io/otel/trace"
)

func performAction(action string) {
	out, err := exec.Command("/bin/bash", action).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}

func initTracer() {
	// Create new OTLP Exporter struct
	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint("localhost:4317"),
		// otlpgrpc.WithDialOption(grpc.WithBlock()), // useful for testing
	)

	exporter, err := otlp.NewExporter(context.Background(), driver)

	if err != nil {
		log.Panicln(err.Error())
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		//sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.NeverSample())),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithIDGenerator(xray.NewIDGenerator()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})
}

func main() {
	action := os.Args[1]

	initTracer()
	ctx := context.Background()
	ctx = baggage.ContextWithValues(ctx)
	tracer := otel.Tracer("Pipeline")

	var span apitrace.Span
	ctx, span = tracer.Start(ctx, "Pipeline Operation")
	
	defer span.End()
	
	
	
	func(ctx context.Context) {
		ctx, span = tracer.Start(ctx, "Pipeline Action")
		defer span.End()

		span.AddEvent("Performing Action", apitrace.WithAttributes(attribute.String("command", action)))
		performAction(action)
		ctx, span = tracer.Start(ctx, "Pipeline Action1")
		defer span.End()
		ctx, span = tracer.Start(ctx, "Pipeline Action2")
		defer span.End()
		ctx, span = tracer.Start(ctx, "Pipeline Action3")
		span.SetAttributes(attribute.String("thingo", "2"))
		span.AddEvent("Nice Op", apitrace.WithAttributes(attribute.String("bogons", "100")))
		defer span.End()
	}(ctx)
}
