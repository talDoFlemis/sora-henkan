# EC2 Instance for AMI Builder
resource "aws_instance" "ami_builder" {
  ami           = data.aws_ami.amazon_linux_2.id
  instance_type = var.app_instance_type

  iam_instance_profile = data.aws_iam_instance_profile.lab_profile.name

  vpc_security_group_ids = [aws_security_group.app_server.id]
  subnet_id              = aws_subnet.public[0].id # Needs public internet to pull docker images

  user_data = base64encode(templatefile("${path.module}/user_data_ami_builder.sh", {
    AWS_REGION              = var.aws_region
    DB_HOST                 = aws_db_instance.main.address
    DB_PORT                 = aws_db_instance.main.port
    DB_USERNAME             = var.db_username
    DB_PASSWORD             = var.db_password
    DB_NAME                 = var.db_name
    S3_BUCKET_NAME          = aws_s3_bucket.images.bucket
    SQS_QUEUE_URL           = aws_sqs_queue.image_queue.url
    OTEL_COLLECTOR_ENDPOINT = "${aws_instance.otel_collector.private_ip}:4317"
    DOCKER_IMAGE_MIGRATE    = "ghcr.io/taldoflemis/sora-henkan/migrate:latest"
    DOCKER_IMAGE_WORKER     = "ghcr.io/taldoflemis/sora-henkan/worker:latest"
    DOCKER_IMAGE_API        = "ghcr.io/taldoflemis/sora-henkan/api:latest"
    DOCKER_IMAGE_FRONTEND   = "ghcr.io/taldoflemis/sora-henkan/frontend:latest"
    API_DOMAIN              = var.api_domain
    ALB_DNS_NAME            = aws_lb.app.dns_name
  }))

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-ami-builder-${var.environment}"
      Type = "ami-builder"
    }
  )
}

# Create AMI from the builder instance
resource "aws_ami_from_instance" "app" {
  name               = "${var.project_name}-app-ami-${var.environment}"
  source_instance_id = aws_instance.ami_builder.id

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-app-ami-${var.environment}"
    }
  )
}
