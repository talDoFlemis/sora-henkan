resource "aws_dynamodb_table" "api_logs" {
  name         = var.dynamodb_logs_table
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-api-logs-${var.environment}"
    }
  )
}
