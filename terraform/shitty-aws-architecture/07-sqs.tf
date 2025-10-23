# SQS Queue for Image Processing
resource "aws_sqs_queue" "image_queue" {
  name                       = "${var.sqs_queue_name}-${var.environment}"
  visibility_timeout_seconds = var.sqs_visibility_timeout
  message_retention_seconds  = var.sqs_message_retention
  receive_wait_time_seconds  = 10 # Enable long polling

  tags = merge(
    var.common_tags,
    {
      Name = "${var.sqs_queue_name}-${var.environment}"
    }
  )
}

# Dead Letter Queue
resource "aws_sqs_queue" "image_queue_dlq" {
  name                      = "${var.sqs_queue_name}-dlq-${var.environment}"
  message_retention_seconds = var.sqs_message_retention

  tags = merge(
    var.common_tags,
    {
      Name = "${var.sqs_queue_name}-dlq-${var.environment}"
    }
  )
}

# SQS Queue Policy for Dead Letter Queue
resource "aws_sqs_queue_redrive_policy" "image_queue" {
  queue_url = aws_sqs_queue.image_queue.id

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.image_queue_dlq.arn
    maxReceiveCount     = 3
  })
}
