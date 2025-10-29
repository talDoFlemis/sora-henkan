#!/bin/bash

set -e

cd /opt/sora-henkan

# Determine API URL based on whether API_DOMAIN is set
if [ -n "${API_DOMAIN}" ]; then
  API_URL="https://${API_DOMAIN}"
else
  API_URL="http://${ALB_DNS_NAME}:42069"
fi

# Create docker-compose file for the application
cat >/opt/sora-henkan/docker-compose.yaml <<EOF
services:
  migrate:
    image: ${DOCKER_IMAGE_MIGRATE}
    container_name: migrate
    environment:
      MIGRATE_DATABASE_HOST: ${DB_HOST}
      MIGRATE_DATABASE_PORT: ${DB_PORT}
      MIGRATE_DATABASE_USER: ${DB_USERNAME}
      MIGRATE_DATABASE_PASSWORD: ${DB_PASSWORD}
      MIGRATE_DATABASE_NAME: ${DB_NAME}
      MIGRATE_DATABASE_SSLMODE: require
    command:
      - "/app/migrate"
      - "-direction=up"
    restart: "no"

  worker:
    image: ${DOCKER_IMAGE_WORKER}
    container_name: worker
    network_mode: host
    environment:
      WORKER_APP_NAME: worker
      WORKER_HTTP_PORT: 8081
      WORKER_DATABASE_HOST: ${DB_HOST}
      WORKER_DATABASE_PORT: ${DB_PORT}
      WORKER_DATABASE_USER: ${DB_USERNAME}
      WORKER_DATABASE_PASSWORD: ${DB_PASSWORD}
      WORKER_DATABASE_NAME: ${DB_NAME}
      WORKER_DATABASE_SSLMODE: require
      WORKER_OBJECTSTORER_ENDPOINT: s3.${AWS_REGION}.amazonaws.com
      WORKER_OBJECTSTORER_USESSL: true
      WORKER_OBJECTSTORER_ACCESSKEYID: ""
      WORKER_OBJECTSTORER_SECRETACCESSKEY: ""
      WORKER_WATERMILL_BROKER_AWS_ENDPOINT: "https://sqs.${AWS_REGION}.amazonaws.com"
      WORKER_WATERMILL_BROKER_AWS_ANONYMOUS: false
      WORKER_IMAGEPROCESSOR_BUCKETNAME: ${S3_BUCKET_NAME}
      AWS_REGION: ${AWS_REGION}
      WORKER_OPENTELEMETRY_ENABLED: true
      WORKER_OPENTELEMETRY_ENDPOINT: ${OTEL_COLLECTOR_ENDPOINT}
      GOMEMLIMIT: 700MiB
      GOGC: 50
    depends_on:
      migrate:
        condition: service_completed_successfully
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '1.5'
          memory: 700M
        reservations:
          memory: 50M

  api:
    image: ${DOCKER_IMAGE_API}
    container_name: api
    network_mode: host
    environment:
      API_APP_NAME: api
      API_DATABASE_HOST: ${DB_HOST}
      API_DATABASE_PORT: ${DB_PORT}
      API_DATABASE_USER: ${DB_USERNAME}
      API_DATABASE_PASSWORD: ${DB_PASSWORD}
      API_DATABASE_NAME: ${DB_NAME}
      API_DATABASE_SSLMODE: require
      API_OBJECTSTORER_ENDPOINT: s3.${AWS_REGION}.amazonaws.com
      API_OBJECTSTORER_USESSL: true
      API_OBJECTSTORER_ACCESSKEYID: ""
      API_OBJECTSTORER_SECRETACCESSKEY: ""
      API_IMAGEPROCESSOR_BUCKETNAME: ${S3_BUCKET_NAME}
      API_WATERMILL_BROKER_AWS_ENDPOINT: "https://sqs.${AWS_REGION}.amazonaws.com"
      API_WATERMILL_BROKER_AWS_ANONYMOUS: false
      AWS_REGION: ${AWS_REGION}
      API_OPENTELEMETRY_ENABLED: true
      API_OPENTELEMETRY_ENDPOINT: ${OTEL_COLLECTOR_ENDPOINT}
      GOMEMLIMIT: 200MiB
      GOGC: 50
    depends_on:
      migrate:
        condition: service_completed_successfully
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 200M
        reservations:
          memory: 50M

  frontend:
    image: ${DOCKER_IMAGE_FRONTEND}
    container_name: frontend
    environment:
      VITE_API_URL: $${API_URL}
      VITE_AWS_BUCKET_ENDPOINT: ${AWS_BUCKET_ENDPOINT}
    ports:
      - "80:8080"
    depends_on:
      - api
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '0.25'
          memory: 100M
        reservations:
          memory: 50M
EOF

/usr/local/bin/docker-compose up -d --pull always
echo "Sora Henkan application started"
