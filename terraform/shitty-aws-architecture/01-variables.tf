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
  default     = "t3.medium"
}

variable "app_instance_type" {
  description = "Instance type for application EC2"
  type        = string
  default     = "t3.micro"
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
  default     = "17.6"
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

# DynamoDB Variables
variable "dynamodb_table_name" {
  description = "Name of the DynamoDB table"
  type        = string
  default     = "main-table"
}

variable "dynamodb_billing_mode" {
  description = "DynamoDB billing mode (PROVISIONED or PAY_PER_REQUEST)"
  type        = string
  default     = "PAY_PER_REQUEST"
  validation {
    condition     = contains(["PROVISIONED", "PAY_PER_REQUEST"], var.dynamodb_billing_mode)
    error_message = "Billing mode must be either PROVISIONED or PAY_PER_REQUEST"
  }
}

variable "dynamodb_read_capacity" {
  description = "Read capacity units (only used if billing_mode is PROVISIONED)"
  type        = number
  default     = 5
}

variable "dynamodb_write_capacity" {
  description = "Write capacity units (only used if billing_mode is PROVISIONED)"
  type        = number
  default     = 5
}

variable "dynamodb_hash_key" {
  description = "Hash key (partition key) attribute name"
  type        = string
  default     = "id"
}

variable "dynamodb_hash_key_type" {
  description = "Hash key attribute type (S, N, or B)"
  type        = string
  default     = "S"
  validation {
    condition     = contains(["S", "N", "B"], var.dynamodb_hash_key_type)
    error_message = "Hash key type must be S (String), N (Number), or B (Binary)"
  }
}

variable "dynamodb_range_key" {
  description = "Range key (sort key) attribute name (leave empty if not using range key)"
  type        = string
  default     = ""
}

variable "dynamodb_range_key_type" {
  description = "Range key attribute type (S, N, or B)"
  type        = string
  default     = "S"
  validation {
    condition     = contains(["S", "N", "B"], var.dynamodb_range_key_type)
    error_message = "Range key type must be S (String), N (Number), or B (Binary)"
  }
}

variable "dynamodb_additional_attributes" {
  description = "Additional attributes for Global Secondary Indexes or Local Secondary Indexes"
  type = list(object({
    name = string
    type = string
  }))
  default = []
}

variable "dynamodb_global_secondary_indexes" {
  description = "Global Secondary Indexes configuration"
  type = list(object({
    name               = string
    hash_key           = string
    range_key          = optional(string)
    projection_type    = string
    non_key_attributes = optional(list(string))
    read_capacity      = optional(number)
    write_capacity     = optional(number)
  }))
  default = []
}

variable "dynamodb_enable_point_in_time_recovery" {
  description = "Enable point-in-time recovery for DynamoDB table"
  type        = bool
  default     = false
}

variable "dynamodb_enable_encryption" {
  description = "Enable server-side encryption for DynamoDB table"
  type        = bool
  default     = true
}

variable "dynamodb_ttl_enabled" {
  description = "Enable TTL for DynamoDB table"
  type        = bool
  default     = false
}

variable "dynamodb_ttl_attribute" {
  description = "Attribute name for TTL (leave empty to disable TTL)"
  type        = string
  default     = ""
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
  default     = ""
  # default     = "sorahenkan.flemis.cloud"
}

variable "api_domain" {
  description = "Domain name for API service (e.g., api.sorahenkan.flemis.cloud). Leave empty for port-based routing only."
  type        = string
  default     = ""
  # default     = "api-sorahenkan.flemis.cloud"
}

# Tags
variable "common_tags" {
  description = "Common tags to apply to all resources"
  type        = map(string)
  default = {
    Project   = "sora-henkan"
    ManagedBy = "Terraform"
  }
}

# DynamoDB Variables
variable "dynamodb_logs_table" {
  description = "Name of the DynamoDB table for API logs"
  type        = string
  default     = "api-logs"
}
