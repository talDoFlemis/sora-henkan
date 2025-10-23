#!/bin/bash
set -e

# Update system
yum update -y

# Install Docker
yum install -y docker
systemctl start docker
systemctl enable docker

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Create directory for OTEL collector
mkdir -p /opt/otel-collector

# Create OTEL Collector configuration
cat > /opt/otel-collector/otel-collector-config.yaml <<EOF
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    timeout: 10s
    send_batch_size: 1024
  
  memory_limiter:
    check_interval: 1s
    limit_mib: 512

exporters:
  logging:
    loglevel: debug
  awsxray:
    region: ${AWS_REGION}
  awsemf:
    region: ${AWS_REGION}
    log_group_name: "/otel/metrics"
    log_stream_name: "otel-collector-metrics"
    namespace: "sora-henkan"
  awscloudwatchlogs:
    region: ${AWS_REGION}
    log_group_name: "/otel/logs"
    log_stream_name: "otel-collector-logs"

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [logging, awsxray]
    
    metrics:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [logging, awsemf]

    logs:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [logging]
  
  extensions: []
EOF

# Create docker-compose file for OTEL collector
cat > /opt/otel-collector/docker-compose.yaml <<EOF
services:
  otel-collector:
    image: otel/opentelemetry-collector-contrib:latest
    container_name: otel-collector
    command: ["--config=/etc/otel-collector-config.yaml"]
    volumes:
      - ./otel-collector-config.yaml:/etc/otel-collector-config.yaml
    ports:
      - "4317:4317"   # OTLP gRPC receiver
      - "4318:4318"   # OTLP HTTP receiver
      - "13133:13133" # health_check extension
    restart: unless-stopped
    environment:
      - AWS_REGION=${AWS_REGION}
EOF

cd /opt/otel-collector
docker-compose up -d

echo "OTEL Collector setup completed"
