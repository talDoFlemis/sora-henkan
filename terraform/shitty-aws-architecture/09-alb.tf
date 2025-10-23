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
    path                = "/health"
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

# ALB Listener for HTTP
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

# Listener Rule: Forward /api/* to API target group
resource "aws_lb_listener_rule" "api" {
  listener_arn = aws_lb_listener.http.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }

  condition {
    path_pattern {
      values = ["/api/*"]
    }
  }
}
