#!/bin/bash
set -e

# Update system
yum update -y

# Install Docker
yum install -y docker
systemctl start docker
systemctl enable docker

# Install cloudwatch agent
yum install -y amazon-cloudwatch-agent

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Create application directory
mkdir -p /opt/sora-henkan

docker pull "${DOCKER_IMAGE_MIGRATE}"
docker pull "${DOCKER_IMAGE_API}"
docker pull "${DOCKER_IMAGE_WORKER}"
docker pull "${DOCKER_IMAGE_FRONTEND}"

echo "Application setup completed"
