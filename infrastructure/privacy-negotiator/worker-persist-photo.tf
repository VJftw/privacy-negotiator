### ECS RabbitMQ containers
data "template_file" "ecs_def_worker-persist-photo" {
  template = "${file("${path.module}/worker.def.tpl.json")}"

  vars {
    environment = "${var.environment}"
    domain      = "${var.domain}"
    weave_cidr  = "${var.weave_cidr}"
    version     = "${var.version}"

    cloudwatch_log_group = "${aws_cloudwatch_log_group.worker-persist-photo.name}"
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

    queue = "persist-photo"
  }
}

resource "aws_ecs_task_definition" "worker-persist-photo" {
  family                = "worker-persist-photo_${var.environment}"
  container_definitions = "${data.template_file.ecs_def_worker-persist-photo.rendered}"
}

resource "aws_ecs_service" "worker-persist-photo" {
  name            = "worker-persist-photo_${var.environment}"
  cluster         = "${var.cluster_name}"
  task_definition = "${aws_ecs_task_definition.worker-persist-photo.arn}"
  desired_count   = 2

  placement_strategy {
    type  = "spread"
    field = "attribute:ecs.availability-zone"
  }
}

#### Log Group for Privacy Negotiator worker-persist-photo
resource "aws_cloudwatch_log_group" "worker-persist-photo" {
  name = "${var.environment}.worker-persist-photo-container-logs"

  retention_in_days = 7

  tags {
    Name        = "worker-persist-photo"
    Environment = "${var.environment}"
  }
}
