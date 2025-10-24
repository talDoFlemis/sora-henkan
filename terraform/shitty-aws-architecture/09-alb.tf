# Application Load Balancer
resource "aws_lb" "app" {
  name               = "${var.project_name}-alb-${var.environment}"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = aws_subnet.public[*].id

  enable_deletion_protection = false
  enable_http2              = true

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-alb-${var.environment}"
    }
  )
}

# Target Group for Frontend (Port 80)
resource "aws_lb_target_group" "frontend" {
  name     = "${var.project_name}-tg-frontend-${var.environment}"
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id

  health_check {
    enabled             = true
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 5
    interval            = 30
    path                = "/"
    protocol            = "HTTP"
    matcher             = "200"
  }

  deregistration_delay = 30

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-tg-frontend-${var.environment}"
    }
  )
}

# Target Group for API (Port 42069)
resource "aws_lb_target_group" "api" {
  name     = "${var.project_name}-tg-api-${var.environment}"
  port     = 42069
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id

  health_check {
    enabled             = true
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 5
    interval            = 30
    path                = "/healthz"
    protocol            = "HTTP"
    port                = "42069"
    matcher             = "200"
  }

  deregistration_delay = 30

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-tg-api-${var.environment}"
    }
  )
}

# ALB Listener for HTTP (Port 80) - Frontend
resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.app.arn
  port              = "80"
  protocol          = "HTTP"

  # Default action forwards to frontend
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.frontend.arn
  }
}

# ALB Listener for API port (42069)
resource "aws_lb_listener" "api_port" {
  load_balancer_arn = aws_lb.app.arn
  port              = "42069"
  protocol          = "HTTP"

  # Default action forwards to API
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }
}

# Listener Rule: Forward API domain to API target group (only if domain is set)
resource "aws_lb_listener_rule" "api_host" {
  count = var.api_domain != "" ? 1 : 0

  listener_arn = aws_lb_listener.http.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }

  condition {
    host_header {
      values = [var.api_domain]
    }
  }
}

# Listener Rule: Forward frontend domain to Frontend target group (only if domain is set)
resource "aws_lb_listener_rule" "frontend_host" {
  count = var.frontend_domain != "" ? 1 : 0

  listener_arn = aws_lb_listener.http.arn
  priority     = 200

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.frontend.arn
  }

  condition {
    host_header {
      values = [var.frontend_domain]
    }
  }
}
