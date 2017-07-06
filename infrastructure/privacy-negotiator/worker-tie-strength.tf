### ECS RabbitMQ containers
data "template_file" "ecs_def_worker-tie-strength" {
  template = "${file("${path.module}/worker.def.tpl.json")}"

  vars {
    environment = "${var.environment}"
    domain      = "${var.domain}"
    weave_cidr  = "${var.weave_cidr}"
    version     = "${var.version}"

    cloudwatch_log_group = "${aws_cloudwatch_log_group.worker-tie-strength.name}"
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

    queue = "tie-strength"
  }
}

resource "aws_ecs_task_definition" "worker-tie-strength" {
  family                = "worker-tie-strength_${var.environment}"
  container_definitions = "${data.template_file.ecs_def_worker-tie-strength.rendered}"
}

resource "aws_ecs_service" "worker-tie-strength" {
  name            = "worker-tie-strength_${var.environment}"
  cluster         = "${var.cluster_name}"
  task_definition = "${aws_ecs_task_definition.worker-tie-strength.arn}"
  desired_count   = 2

  placement_strategy {
    type  = "spread"
    field = "attribute:ecs.availability-zone"
  }
}

#### Log Group for Privacy Negotiator worker-tie-strength
resource "aws_cloudwatch_log_group" "worker-tie-strength" {
  name = "${var.environment}.worker-tie-strength-container-logs"

  retention_in_days = 7

  tags {
    Name        = "worker-tie-strength"
    Environment = "${var.environment}"
  }
}

resource "aws_cloudwatch_log_metric_filter" "worker-tie-strength_error" {

  name = "${var.environment}-worker-tie-strength.error"
  pattern = "error"
  log_group_name = "${aws_cloudwatch_log_group.worker-tie-strength.name}"

  metric_transformation {
    name = "${var.environment}-worker-tie-strength.error"
    namespace = "${var.environment}-worker-tie-strength"
    value = "1"
  }

}

resource "aws_cloudwatch_log_metric_filter" "worker-tie-strength_error_reset" {

  name = "${var.environment}-worker-tie-strength.error"
  pattern = ""
  log_group_name = "${aws_cloudwatch_log_group.worker-tie-strength.name}"

  metric_transformation {
    name = "${var.environment}-worker-tie-strength.error"
    namespace = "${var.environment}-worker-tie-strength"
    value = "0"
  }

}

resource "aws_cloudwatch_metric_alarm" "worker-tie-strength_error" {

  alarm_name = "${var.environment}.worker-tie-strength.error"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  threshold = "1"
  period = "60"
  statistic = "Sum"
  evaluation_periods = "1"
  metric_name = "${var.environment}.worker-tie-strength.error"
  namespace = "${var.environment}-worker-tie-strength.error"
  alarm_description = "monitors log for worker-tie-strength errors"
  //  alarm_actions = ["arn:aws:sns:eu-west-1:812414252941:error_notification"]

}