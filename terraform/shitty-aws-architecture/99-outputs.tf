output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = aws_subnet.private[*].id
}

output "nat_gateway_id" {
  description = "ID of the NAT Gateway"
  value       = aws_nat_gateway.main.id
}

output "otel_collector_private_ip" {
  description = "Private IP of the OTEL Collector instance"
  value       = aws_instance.otel_collector.private_ip
}

output "otel_collector_instance_id" {
  description = "Instance ID of the OTEL Collector"
  value       = aws_instance.otel_collector.id
}

output "load_balancer_dns" {
  description = "DNS name of the Application Load Balancer"
  value       = aws_lb.app.dns_name
}

output "load_balancer_frontend_url" {
  description = "Frontend URL (port-based access)"
  value       = "http://${aws_lb.app.dns_name}:80"
}

output "load_balancer_api_url" {
  description = "API URL (port-based access)"
  value       = "http://${aws_lb.app.dns_name}:42069"
}

output "load_balancer_jaeger_url" {
  description = "Jaeger UI URL (port-based access)"
  value       = "http://${aws_lb.app.dns_name}:16686"
}

output "frontend_domain_url" {
  description = "Frontend domain URL (if configured)"
  value       = var.frontend_domain != "" ? "https://${var.frontend_domain}" : "Not configured - use port-based access"
}

output "api_domain_url" {
  description = "API domain URL (if configured)"
  value       = var.api_domain != "" ? "https://${var.api_domain}" : "Not configured - use port-based access"
}

output "autoscaling_group_name" {
  description = "Name of the Auto Scaling Group"
  value       = aws_autoscaling_group.app.name
}

output "s3_bucket_name" {
  description = "Name of the S3 bucket for images"
  value       = aws_s3_bucket.images.id
}

output "s3_bucket_arn" {
  description = "ARN of the S3 bucket for images"
  value       = aws_s3_bucket.images.arn
}

output "sqs_queue_url" {
  description = "URL of the SQS queue"
  value       = aws_sqs_queue.image_queue.url
}

output "sqs_queue_arn" {
  description = "ARN of the SQS queue"
  value       = aws_sqs_queue.image_queue.arn
}

output "sqs_dlq_url" {
  description = "URL of the SQS Dead Letter Queue"
  value       = aws_sqs_queue.image_queue_dlq.url
}

output "launch_template_id" {
  description = "ID of the Launch Template"
  value       = aws_launch_template.app.id
}

output "launch_template_latest_version" {
  description = "Latest version of the Launch Template"
  value       = aws_launch_template.app.latest_version
}

output "iam_role_arn" {
  description = "ARN of the LabRole used for EC2 instances"
  value       = data.aws_iam_role.lab_role.arn
}

output "iam_instance_profile_name" {
  description = "Name of the LabInstanceProfile used for EC2 instances"
  value       = data.aws_iam_instance_profile.lab_profile.name
}

output "security_group_alb_id" {
  description = "ID of the ALB security group"
  value       = aws_security_group.alb.id
}

output "security_group_app_server_id" {
  description = "ID of the application server security group"
  value       = aws_security_group.app_server.id
}

output "security_group_otel_collector_id" {
  description = "ID of the OTEL collector security group"
  value       = aws_security_group.otel_collector.id
}

output "otel_collector_otlp_grpc_endpoint" {
  description = "OTEL Collector OTLP gRPC endpoint (internal)"
  value       = "http://${aws_instance.otel_collector.private_ip}:4317"
}

output "otel_collector_otlp_http_endpoint" {
  description = "OTEL Collector OTLP HTTP endpoint (internal)"
  value       = "http://${aws_instance.otel_collector.private_ip}:4318"
}

# RDS Outputs
output "rds_endpoint" {
  description = "RDS instance endpoint"
  value       = aws_db_instance.main.endpoint
}

output "rds_address" {
  description = "RDS instance address"
  value       = aws_db_instance.main.address
}

output "rds_port" {
  description = "RDS instance port"
  value       = aws_db_instance.main.port
}

output "rds_database_name" {
  description = "RDS database name"
  value       = aws_db_instance.main.db_name
}

# DynamoDB Outputs
output "dynamodb_table_name" {
  description = "Name of the DynamoDB table"
  value       = aws_dynamodb_table.main.name
}

output "dynamodb_table_arn" {
  description = "ARN of the DynamoDB table"
  value       = aws_dynamodb_table.main.arn
}

output "dynamodb_table_id" {
  description = "ID of the DynamoDB table"
  value       = aws_dynamodb_table.main.id
}

output "dynamodb_vpc_endpoint_id" {
  description = "ID of the DynamoDB VPC endpoint"
  value       = aws_vpc_endpoint.dynamodb.id
}

# Access Instructions
output "access_instructions" {
  description = "Instructions for accessing your application"
  value       = <<-EOT
  
  ====================================
  APPLICATION ACCESS
  ====================================
  
  PORT-BASED ACCESS (direct to ALB):
  Frontend: http://${aws_lb.app.dns_name}:80
  API:      http://${aws_lb.app.dns_name}:42069
  
  ${var.frontend_domain != "" || var.api_domain != "" ? "DOMAIN-BASED ACCESS (via Cloudflare):" : ""}
  ${var.frontend_domain != "" ? "Frontend: https://${var.frontend_domain}" : ""}
  ${var.api_domain != "" ? "API:      https://${var.api_domain}" : ""}
  
  ${var.frontend_domain != "" || var.api_domain != "" ? "====================================\nCLOUDFLARE DNS CONFIGURATION\n====================================\n\nAdd these CNAME records in Cloudflare:" : ""}
  ${var.frontend_domain != "" ? "\nFrontend Domain:\nName:    ${replace(var.frontend_domain, ".flemis.cloud", "")}\nType:    CNAME\nContent: ${aws_lb.app.dns_name}\nProxy:   ON (Proxied - Orange Cloud)\nTTL:     Auto" : ""}
  ${var.api_domain != "" ? "\nAPI Domain:\nName:    ${replace(var.api_domain, ".flemis.cloud", "")}\nType:    CNAME\nContent: ${aws_lb.app.dns_name}\nProxy:   ON (Proxied - Orange Cloud)\nTTL:     Auto" : ""}
  
  ${var.frontend_domain != "" || var.api_domain != "" ? "\nSSL/TLS Settings in Cloudflare:\nGo to SSL/TLS → Overview\nSet mode to: Flexible\n\nBenefits:\n- Free SSL/TLS from Cloudflare\n- CDN and caching\n- DDoS protection\n- Analytics" : "\nTo enable domain-based routing:\n- Set frontend_domain variable to your frontend domain\n- Set api_domain variable to your API domain\n- Add corresponding CNAME records in Cloudflare"}
  
  ====================================
  EOT
}
