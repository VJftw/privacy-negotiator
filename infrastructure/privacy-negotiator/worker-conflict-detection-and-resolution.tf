### ECS RabbitMQ containers
data "template_file" "ecs_def_worker-conflict-detection-and-resolution" {
  template = "${file("${path.module}/worker.def.tpl.json")}"

  vars {
    environment = "${var.environment}"
    domain      = "${var.domain}"
    weave_cidr  = "${var.weave_cidr}"
    version     = "${var.version}"

    cloudwatch_log_group = "${aws_cloudwatch_log_group.worker-conflict-detection-and-resolution.name}"
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

    queue = "conflict-detection-and-resolution"
  }
}

resource "aws_ecs_task_definition" "worker-conflict-detection-and-resolution" {
  family                = "worker-conflict-detection-and-resolution_${var.environment}"
  container_definitions = "${data.template_file.ecs_def_worker-conflict-detection-and-resolution.rendered}"
}

resource "aws_ecs_service" "worker-conflict-detection-and-resolution" {
  name            = "worker-conflict-detection-and-resolution_${var.environment}"
  cluster         = "${var.cluster_name}"
  task_definition = "${aws_ecs_task_definition.worker-conflict-detection-and-resolution.arn}"
  desired_count   = 2

  placement_strategy {
    type  = "spread"
    field = "attribute:ecs.availability-zone"
  }
}

#### Log Group for Privacy Negotiator worker-conflict-detection-and-resolution
resource "aws_cloudwatch_log_group" "worker-conflict-detection-and-resolution" {
  name = "${var.environment}.worker-conflict-detection-and-resolution-container-logs"

  retention_in_days = 7

  tags {
    Name        = "worker-conflict-detection-and-resolution"
    Environment = "${var.environment}"
  }
}
