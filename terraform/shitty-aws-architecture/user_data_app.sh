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

# Get instance metadata
OTEL_ENDPOINT="${OTEL_COLLECTOR_ENDPOINT}"
DB_HOST="${DB_HOST}"
DB_PORT="${DB_PORT}"

# Create application directory
mkdir -p /opt/sora-henkan

# Create docker-compose file for the application
cat > /opt/sora-henkan/docker-compose.yaml <<EOF
services:
  migrate:
    image: ${DOCKER_IMAGE_MIGRATE}
    container_name: migrate
    environment:
      MIGRATE_DATABASE_HOST: \$DB_HOST
      MIGRATE_DATABASE_PORT: \$DB_PORT
      MIGRATE_DATABASE_USER: ${DB_USERNAME}
      MIGRATE_DATABASE_PASSWORD: ${DB_PASSWORD}
      MIGRATE_DATABASE_NAME: ${DB_NAME}
    command:
      - "/app/migrate"
      - "-direction=up"
    restart: "no"

  worker:
    image: ${DOCKER_IMAGE_WORKER}
    container_name: worker
    environment:
      WORKER_DATABASE_HOST: \$DB_HOST
      WORKER_DATABASE_PORT: \$DB_PORT
      WORKER_DATABASE_USER: ${DB_USERNAME}
      WORKER_DATABASE_PASSWORD: ${DB_PASSWORD}
      WORKER_DATABASE_NAME: ${DB_NAME}
      WORKER_OBJECTSTORER_ENDPOINT: ${S3_BUCKET_NAME}
      WORKER_WATERMILL_BROKER_AWS_ENDPOINT: ${SQS_QUEUE_URL}
      AWS_REGION: ${AWS_REGION}
      WORKER_OPENTELEMETRY_ENDPOINT: \$OTEL_ENDPOINT
    depends_on:
      migrate:
        condition: service_completed_successfully
    restart: unless-stopped

  api:
    image: ${DOCKER_IMAGE_API}
    container_name: api
    environment:
      API_DATABASE_HOST: \$DB_HOST
      API_DATABASE_PORT: \$DB_PORT
      API_DATABASE_USER: ${DB_USERNAME}
      API_DATABASE_PASSWORD: ${DB_PASSWORD}
      API_DATABASE_NAME: ${DB_NAME}
      API_OBJECTSTORER_ENDPOINT: ${S3_BUCKET_NAME}
      API_WATERMILL_BROKER_AWS_ENDPOINT: ${SQS_QUEUE_URL}
      AWS_REGION: ${AWS_REGION}
      API_OPENTELEMETRY_ENDPOINT: \$OTEL_ENDPOINT
    ports:
      - "42069:42069"
    depends_on:
      migrate:
        condition: service_completed_successfully
    restart: unless-stopped

  frontenzo:
    image: traefik/whoami
    container_name: frontenzo
    ports:
      - "80:80"
    depends_on:
      - api
    restart: unless-stopped
EOF

# Start the application
cd /opt/sora-henkan
docker-compose up -d

echo "Application setup completed"
