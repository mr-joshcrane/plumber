#!/bin/bash
GO111MODULE=on go get github.com/open-telemetry/opentelemetry-collector-builder
cat > ~/.otelcol-builder.yaml <<EOF
dist:
  go: "/usr/local/bin/go"
replaces:
  - github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray => github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray v0.28.0
  - github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/awsutil => github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/awsutil v0.28.0
exporters:
  - gomod: "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsxrayexporter v0.27.0"
EOF
opentelemetry-collector-builder --output-path=/tmp/dist
cat > /tmp/otelcol.yaml <<EOF
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: localhost:4317

exporters:
  logging:
    loglevel: debug
  awsxray:
  awsxray/customname:
    region: ap-southeast-2


service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [logging, awsxray]
EOF
/tmp/dist/otelcol-custom --config=/tmp/otelcol.yaml