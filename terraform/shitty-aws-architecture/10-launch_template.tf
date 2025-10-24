# Launch Template for Application Servers
resource "aws_launch_template" "app" {
  name_prefix   = "${var.project_name}-app-lt-${var.environment}"
  image_id      = aws_ami_from_instance.app.id
  instance_type = var.app_instance_type

  iam_instance_profile {
    name = data.aws_iam_instance_profile.lab_profile.name
  }

  vpc_security_group_ids = [aws_security_group.app_server.id]

  user_data = base64encode(file("${path.module}/user_data_app.sh"))

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
