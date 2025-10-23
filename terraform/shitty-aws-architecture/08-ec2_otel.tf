# Data source for latest Amazon Linux 2 AMI
data "aws_ami" "amazon_linux_2" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

# EC2 Instance for OTEL Collector (in private subnet)
resource "aws_instance" "otel_collector" {
  ami                    = data.aws_ami.amazon_linux_2.id
  instance_type          = var.otel_instance_type
  subnet_id              = aws_subnet.private[0].id
  vpc_security_group_ids = [aws_security_group.otel_collector.id]
  iam_instance_profile   = data.aws_iam_instance_profile.lab_profile.name
  # key_name is not required - can use AWS Systems Manager Session Manager for access

  user_data = templatefile("${path.module}/user_data_otel.sh", {
    AWS_REGION = var.aws_region
  })

  root_block_device {
    volume_type = "gp3"
    volume_size = 20
    encrypted   = true
  }

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-otel-collector-${var.environment}"
      Type = "otel-collector"
    }
  )

  lifecycle {
    create_before_destroy = true
  }
}
