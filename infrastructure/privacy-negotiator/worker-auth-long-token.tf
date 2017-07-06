### ECS RabbitMQ containers
data "template_file" "ecs_def_worker-auth-queue" {
  template = "${file("${path.module}/worker.def.tpl.json")}"

  vars {
    environment = "${var.environment}"
    domain      = "${var.domain}"
    weave_cidr  = "${var.weave_cidr}"
    version     = "${var.version}"

    cloudwatch_log_group = "${aws_cloudwatch_log_group.worker-auth-queue.name}"
    cloudwatch_region    = "${var.aws_region}"

    rabbitmq_user     = "${var.rabbitmq_user}"
    rabbitmq_pass     = "${var.rabbitmq_pass}"
    rabbitmq_hostname = "rabbitmq-${var.environment}.weave.local"
    redis_host        = "redis-${var.environment}.weave.local"

    postgres_user     = "${var.postgres_user}"
    postgres_dbname   = "privneg_${var.environment}"
    postgres_password = "${var.postgres_password}"
    postgres_host     = "${aws_db_instance.db.address}"

    facebook_app_id     = "${var.facebook_app_id}"
    facebook_app_secret = "${var.facebook_app_secret}"

    queue = "auth-long-token"
  }
}

resource "aws_ecs_task_definition" "worker-auth-queue" {
  family                = "worker-auth-queue_${var.environment}"
  container_definitions = "${data.template_file.ecs_def_worker-auth-queue.rendered}"
}

resource "aws_ecs_service" "worker-auth-queue" {
  name            = "worker-auth-queue_${var.environment}"
  cluster         = "${var.cluster_name}"
  task_definition = "${aws_ecs_task_definition.worker-auth-queue.arn}"
  desired_count   = 2

  placement_strategy {
    type  = "spread"
    field = "attribute:ecs.availability-zone"
  }
}

#### Log Group for Privacy Negotiator worker-auth-queue
resource "aws_cloudwatch_log_group" "worker-auth-queue" {
  name = "${var.environment}.worker-auth-queue-container-logs"

  retention_in_days = 7

  tags {
    Name        = "worker-auth-queue"
    Environment = "${var.environment}"
  }
}

resource "aws_cloudwatch_log_metric_filter" "worker-auth-queue_error" {

  name = "${var.environment}-worker-auth-queue.error"
  pattern = "error"
  log_group_name = "${aws_cloudwatch_log_group.worker-auth-queue.name}"

  metric_transformation {
    name = "${var.environment}-worker-auth-queue.error"
    namespace = "${var.environment}-worker-auth-queue"
    value = "1"
  }

}

resource "aws_cloudwatch_log_metric_filter" "worker-auth-queue_error_reset" {

  name = "${var.environment}-worker-auth-queue.error"
  pattern = ""
  log_group_name = "${aws_cloudwatch_log_group.worker-auth-queue.name}"

  metric_transformation {
    name = "${var.environment}-worker-auth-queue.error"
    namespace = "${var.environment}-worker-auth-queue"
    value = "0"
  }

}

resource "aws_cloudwatch_metric_alarm" "worker-auth-queue_error" {

  alarm_name = "${var.environment}.worker-auth-queue.error"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  threshold = "1"
  period = "60"
  statistic = "Sum"
  evaluation_periods = "1"
  metric_name = "${var.environment}.worker-auth-queue.error"
  namespace = "${var.environment}-worker-auth-queue.error"
  alarm_description = "monitors log for worker-auth-queue errors"
  //  alarm_actions = ["arn:aws:sns:eu-west-1:812414252941:error_notification"]

}