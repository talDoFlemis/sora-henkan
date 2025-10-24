# Security Group for OTEL Collector
resource "aws_security_group" "otel_collector" {
  name        = "${var.project_name}-otel-collector-sg-${var.environment}"
  description = "Security group for OTEL collector EC2 - only accessible from app instances"
  vpc_id      = aws_vpc.main.id

  # OTEL Collector OTLP gRPC - only from app instances
  ingress {
    description     = "OTLP gRPC from application instances"
    from_port       = 4317
    to_port         = 4317
    protocol        = "tcp"
    security_groups = [aws_security_group.app_server.id]
  }

  # OTEL Collector OTLP HTTP - only from app instances
  ingress {
    description     = "OTLP HTTP from application instances"
    from_port       = 4318
    to_port         = 4318
    protocol        = "tcp"
    security_groups = [aws_security_group.app_server.id]
  }

  ingress {
    description     = "Jaeger UI from LB"
    from_port       = 16686
    to_port         = 16686
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  # OTEL Collector health check - only from VPC
  ingress {
    description = "Health check from VPC"
    from_port   = 13133
    to_port     = 13133
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
  }

  # Allow all outbound traffic
  egress {
    description = "All outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-otel-collector-sg-${var.environment}"
    }
  )
}

# Security Group for Application Servers
resource "aws_security_group" "app_server" {
  name        = "${var.project_name}-app-server-sg-${var.environment}"
  description = "Security group for application EC2 instances"
  vpc_id      = aws_vpc.main.id

  # HTTP from load balancer
  ingress {
    description     = "HTTP from ALB"
    from_port       = 80
    to_port         = 80
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  # Application API port
  ingress {
    description     = "API from ALB"
    from_port       = 42069
    to_port         = 42069
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  # Allow all outbound traffic
  egress {
    description = "All outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-app-server-sg-${var.environment}"
    }
  )
}

# Security Group for Application Load Balancer
resource "aws_security_group" "alb" {
  name        = "${var.project_name}-alb-sg-${var.environment}"
  description = "Security group for Application Load Balancer"
  vpc_id      = aws_vpc.main.id

  # HTTP access
  ingress {
    description = "HTTP from internet"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTPS access
  ingress {
    description = "HTTPS from internet"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # API port access (for direct ALB access without domain)
  ingress {
    description = "API port from internet"
    from_port   = 42069
    to_port     = 42069
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Jaeger UI port
  ingress {
    description = "Jaeger UI from internet"
    from_port   = 16686
    to_port     = 16686
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow all outbound traffic
  egress {
    description = "All outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-alb-sg-${var.environment}"
    }
  )
}

# Security Group for RDS
resource "aws_security_group" "rds" {
  name        = "${var.project_name}-rds-sg-${var.environment}"
  description = "Security group for RDS PostgreSQL - ONLY accessible from app instances"
  vpc_id      = aws_vpc.main.id

  # PostgreSQL access from app instances ONLY
  ingress {
    description     = "PostgreSQL from application instances ONLY"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.app_server.id]
  }

  # Allow all outbound traffic
  egress {
    description = "All outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-rds-sg-${var.environment}"
    }
  )
}
