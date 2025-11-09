terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.17.0"
    }
  }
}

provider "aws" {
  region = var.aws_region

  # Use LocalStack endpoint if enabled
  skip_credentials_validation = var.use_localstack
  skip_metadata_api_check     = var.use_localstack
  skip_requesting_account_id  = var.use_localstack

  endpoints {
    ec2            = var.use_localstack ? var.localstack_endpoint : null
    s3             = var.use_localstack ? var.localstack_endpoint : null
    sqs            = var.use_localstack ? var.localstack_endpoint : null
    iam            = var.use_localstack ? var.localstack_endpoint : null
    autoscaling    = var.use_localstack ? var.localstack_endpoint : null
    cloudwatch     = var.use_localstack ? var.localstack_endpoint : null
    elb            = var.use_localstack ? var.localstack_endpoint : null
    elbv2          = var.use_localstack ? var.localstack_endpoint : null
    dynamodb       = var.use_localstack ? var.localstack_endpoint : null
  }

  # For LocalStack, use dummy credentials
  access_key = var.use_localstack ? "test" : null
  secret_key = var.use_localstack ? "test" : null
}
