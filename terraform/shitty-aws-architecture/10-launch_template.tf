# Launch Template for Application Servers
resource "aws_launch_template" "app" {
  name_prefix   = "${var.project_name}-app-lt-${var.environment}"
  image_id      = data.aws_ami.amazon_linux_2.id
  instance_type = var.app_instance_type

  iam_instance_profile {
    name = data.aws_iam_instance_profile.lab_profile.name
  }

  vpc_security_group_ids = [aws_security_group.app_server.id]

  user_data = base64encode(templatefile("${path.module}/user_data_app.sh", {
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
  }))

  block_device_mappings {
    device_name = "/dev/xvda"

    ebs {
      volume_size           = 30
      volume_type           = "gp3"
      delete_on_termination = true
      encrypted             = true
    }
  }

  monitoring {
    enabled = true
  }

  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"
    http_put_response_hop_limit = 1
  }

  tag_specifications {
    resource_type = "instance"

    tags = merge(
      var.common_tags,
      {
        Name = "${var.project_name}-app-${var.environment}"
        Type = "app-server"
      }
    )
  }

  tag_specifications {
    resource_type = "volume"

    tags = merge(
      var.common_tags,
      {
        Name = "${var.project_name}-app-volume-${var.environment}"
      }
    )
  }

  lifecycle {
    create_before_destroy = true
  }

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-app-lt-${var.environment}"
    }
  )
}
