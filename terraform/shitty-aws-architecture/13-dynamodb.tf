# DynamoDB Table
resource "aws_dynamodb_table" "main" {
  name           = "${var.project_name}-${var.environment}-${var.dynamodb_table_name}"
  billing_mode   = var.dynamodb_billing_mode
  hash_key       = var.dynamodb_hash_key
  range_key      = var.dynamodb_range_key != "" ? var.dynamodb_range_key : null
  
  # Provisioned throughput (only used if billing_mode is PROVISIONED)
  read_capacity  = var.dynamodb_billing_mode == "PROVISIONED" ? var.dynamodb_read_capacity : null
  write_capacity = var.dynamodb_billing_mode == "PROVISIONED" ? var.dynamodb_write_capacity : null

  # Primary key attribute
  attribute {
    name = var.dynamodb_hash_key
    type = var.dynamodb_hash_key_type
  }

  # Range key attribute (if provided)
  dynamic "attribute" {
    for_each = var.dynamodb_range_key != "" ? [1] : []
    content {
      name = var.dynamodb_range_key
      type = var.dynamodb_range_key_type
    }
  }

  # Additional attributes for GSI/LSI
  dynamic "attribute" {
    for_each = var.dynamodb_additional_attributes
    content {
      name = attribute.value.name
      type = attribute.value.type
    }
  }

  # Point-in-time recovery
  point_in_time_recovery {
    enabled = var.dynamodb_enable_point_in_time_recovery
  }

  # Server-side encryption
  server_side_encryption {
    enabled = var.dynamodb_enable_encryption
  }

  # TTL configuration
  dynamic "ttl" {
    for_each = var.dynamodb_ttl_attribute != "" ? [1] : []
    content {
      attribute_name = var.dynamodb_ttl_attribute
      enabled        = var.dynamodb_ttl_enabled
    }
  }

  # Global Secondary Indexes
  dynamic "global_secondary_index" {
    for_each = var.dynamodb_global_secondary_indexes
    content {
      name               = global_secondary_index.value.name
      hash_key           = global_secondary_index.value.hash_key
      range_key          = lookup(global_secondary_index.value, "range_key", null)
      projection_type    = global_secondary_index.value.projection_type
      non_key_attributes = lookup(global_secondary_index.value, "non_key_attributes", null)
      read_capacity      = var.dynamodb_billing_mode == "PROVISIONED" ? lookup(global_secondary_index.value, "read_capacity", 5) : null
      write_capacity     = var.dynamodb_billing_mode == "PROVISIONED" ? lookup(global_secondary_index.value, "write_capacity", 5) : null
    }
  }

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-${var.environment}-${var.dynamodb_table_name}"
    }
  )
}

# VPC Endpoint for DynamoDB (Gateway Endpoint - no cost)
# This allows EC2 instances in private subnets to access DynamoDB without going through internet
resource "aws_vpc_endpoint" "dynamodb" {
  vpc_id       = aws_vpc.main.id
  service_name = "com.amazonaws.${var.aws_region}.dynamodb"
  
  # Gateway endpoints use route tables
  route_table_ids = concat(
    aws_route_table.private[*].id,
    [aws_route_table.public.id]
  )

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-dynamodb-endpoint-${var.environment}"
    }
  )
}

# VPC Endpoint Policy for DynamoDB (restrict access to only our table)
resource "aws_vpc_endpoint_policy" "dynamodb" {
  vpc_endpoint_id = aws_vpc_endpoint.dynamodb.id
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = "*"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan",
          "dynamodb:BatchGetItem",
          "dynamodb:BatchWriteItem",
          "dynamodb:DescribeTable"
        ]
        Resource = [
          aws_dynamodb_table.main.arn,
          "${aws_dynamodb_table.main.arn}/index/*"
        ]
      }
    ]
  })
}