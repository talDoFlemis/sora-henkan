resource "aws_lb_target_group_attachment" "jaeger" {
  target_group_arn = aws_lb_target_group.jaeger.arn
  target_id        = aws_instance.otel_collector.id
  port             = 16686
}
