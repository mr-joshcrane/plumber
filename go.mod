module github.com/mr-joshcrane/plumber

go 1.16

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/awsutil => github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/awsutil v0.28.0

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray => github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray v0.28.0

require (
	github.com/open-telemetry/opentelemetry-collector-builder v0.27.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsxrayexporter v0.28.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/awsutil v0.28.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.opentelemetry.io/collector v0.28.0 // indirect
	go.opentelemetry.io/contrib/propagators/aws v0.20.0 // indirect
	go.opentelemetry.io/otel v0.20.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp v0.20.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout v0.19.0 // indirect
	go.opentelemetry.io/otel/sdk v0.20.0 // indirect
	go.opentelemetry.io/otel/trace v0.20.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
)
