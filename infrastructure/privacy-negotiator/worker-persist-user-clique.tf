### ECS RabbitMQ containers
data "template_file" "ecs_def_worker-persist-user-clique" {
  template = "${file("${path.module}/worker.def.tpl.json")}"

  vars {
    environment = "${var.environment}"
    domain      = "${var.domain}"
    weave_cidr  = "${var.weave_cidr}"
    version     = "${var.version}"

    cloudwatch_log_group = "${aws_cloudwatch_log_group.worker-persist-user-clique.name}"
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

    queue = "persist-user-clique"
  }
}

resource "aws_ecs_task_definition" "worker-persist-user-clique" {
  family                = "worker-persist-user-clique_${var.environment}"
  container_definitions = "${data.template_file.ecs_def_worker-persist-user-clique.rendered}"
}

resource "aws_ecs_service" "worker-persist-user-clique" {
  name            = "worker-persist-user-clique_${var.environment}"
  cluster         = "${var.cluster_name}"
  task_definition = "${aws_ecs_task_definition.worker-persist-user-clique.arn}"
  desired_count   = 1

  placement_strategy {
    type  = "spread"
    field = "attribute:ecs.availability-zone"
  }
}

#### Log Group for Privacy Negotiator worker-persist-user-clique
resource "aws_cloudwatch_log_group" "worker-persist-user-clique" {
  name = "${var.environment}.worker-persist-user-clique-container-logs"

  retention_in_days = 7

  tags {
    Name        = "worker-persist-user-clique"
    Environment = "${var.environment}"
  }
}

resource "aws_cloudwatch_log_metric_filter" "worker-persist-user-clique_error" {

  name = "${var.environment}-worker-persist-user-clique.error"
  pattern = "error"
  log_group_name = "${aws_cloudwatch_log_group.worker-persist-user-clique.name}"

  metric_transformation {
    name = "${var.environment}-worker-persist-user-clique.error"
    namespace = "${var.environment}-worker-persist-user-clique"
    value = "1"
  }

}

resource "aws_cloudwatch_log_metric_filter" "worker-persist-user-clique_error_reset" {

  name = "${var.environment}-worker-persist-user-clique.error"
  pattern = ""
  log_group_name = "${aws_cloudwatch_log_group.worker-persist-user-clique.name}"

  metric_transformation {
    name = "${var.environment}-worker-persist-user-clique.error"
    namespace = "${var.environment}-worker-persist-user-clique"
    value = "0"
  }

}

resource "aws_cloudwatch_metric_alarm" "worker-persist-user-clique_error" {

  alarm_name = "${var.environment}.worker-persist-user-clique.error"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  threshold = "1"
  period = "60"
  statistic = "Sum"
  evaluation_periods = "1"
  metric_name = "${var.environment}.worker-persist-user-clique.error"
  namespace = "${var.environment}-worker-persist-user-clique.error"
  alarm_description = "monitors log for worker-persist-user-clique errors"
  //  alarm_actions = ["arn:aws:sns:eu-west-1:812414252941:error_notification"]

}