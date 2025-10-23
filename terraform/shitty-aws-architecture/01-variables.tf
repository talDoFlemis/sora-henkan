variable "aws_region" {
  description = "AWS region for resources"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
  default     = "sora-henkan"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "CIDR blocks for public subnets"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for private subnets"
  type        = list(string)
  default     = ["10.0.10.0/24", "10.0.11.0/24"]
}

variable "availability_zones" {
  description = "Availability zones"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b"]
}

# EC2 Variables
variable "otel_instance_type" {
  description = "Instance type for OTEL collector EC2"
  type        = string
  default     = "t3.small"
}

variable "app_instance_type" {
  description = "Instance type for application EC2"
  type        = string
  default     = "t3.medium"
}

# Note: SSH key pair not required
# EC2 instances in private subnets can be accessed via AWS Systems Manager Session Manager
# which uses the LabRole IAM permissions

# Auto Scaling Variables
variable "asg_min_size" {
  description = "Minimum number of instances in ASG"
  type        = number
  default     = 1
}

variable "asg_max_size" {
  description = "Maximum number of instances in ASG"
  type        = number
  default     = 3
}

variable "asg_desired_capacity" {
  description = "Desired number of instances in ASG"
  type        = number
  default     = 1
}

variable "scale_up_cpu_threshold" {
  description = "CPU threshold percentage to trigger scale up"
  type        = number
  default     = 70
}

variable "scale_down_cpu_threshold" {
  description = "CPU threshold percentage to trigger scale down"
  type        = number
  default     = 25
}

variable "scale_up_cooldown" {
  description = "Cooldown period after scale up (seconds)"
  type        = number
  default     = 60
}

variable "scale_down_cooldown" {
  description = "Cooldown period after scale down (seconds)"
  type        = number
  default     = 300
}

# S3 Variables
variable "image_bucket_name" {
  description = "Name of the S3 bucket for storing images"
  type        = string
  default     = "sora-henkan-images"
}

# SQS Variables
variable "sqs_queue_name" {
  description = "Name of the SQS queue"
  type        = string
  default     = "sora-henkan-image-queue"
}

variable "sqs_visibility_timeout" {
  description = "SQS visibility timeout in seconds"
  type        = number
  default     = 300
}

variable "sqs_message_retention" {
  description = "SQS message retention period in seconds"
  type        = number
  default     = 1209600 # 14 days
}

# Database Variables
variable "db_username" {
  description = "Database username"
  type        = string
  default     = "sorahenkan"
  sensitive   = true
}

variable "db_password" {
  description = "Database password"
  type        = string
  default     = "sorahenkan_password"
  sensitive   = true
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "sorahenkan"
}

variable "db_engine_version" {
  description = "PostgreSQL engine version"
  type        = string
  default     = "16.3"
}

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro"
}

variable "db_allocated_storage" {
  description = "Allocated storage for RDS in GB"
  type        = number
  default     = 20
}

variable "db_backup_retention_period" {
  description = "Number of days to retain backups"
  type        = number
  default     = 7
}

variable "db_multi_az" {
  description = "Enable Multi-AZ deployment for high availability"
  type        = bool
  default     = false
}

variable "db_deletion_protection" {
  description = "Enable deletion protection"
  type        = bool
  default     = false
}

variable "db_skip_final_snapshot" {
  description = "Skip final snapshot when destroying"
  type        = bool
  default     = true
}

# LocalStack Variables
variable "localstack_endpoint" {
  description = "LocalStack endpoint URL"
  type        = string
  default     = "http://localhost:4566"
}

variable "use_localstack" {
  description = "Whether to use LocalStack for local development"
  type        = bool
  default     = false
}

# Domain Variables
variable "frontend_domain" {
  description = "Domain name for frontend service (e.g., sorahenkan.flemis.cloud). Leave empty for port-based routing only."
  type        = string
  default     = "sorahenkan.flemis.cloud"
}

variable "api_domain" {
  description = "Domain name for API service (e.g., api.sorahenkan.flemis.cloud). Leave empty for port-based routing only."
  type        = string
  default     = "api-sorahenkan.flemis.cloud"
}

# Tags
variable "common_tags" {
  description = "Common tags to apply to all resources"
  type        = map(string)
  default = {
    Project     = "sora-henkan"
    ManagedBy   = "Terraform"
  }
}
