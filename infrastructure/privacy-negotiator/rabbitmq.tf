### ECS RabbitMQ containers
data "template_file" "ecs_def_rabbitmq" {
  template = "${file("${path.module}/rabbitmq.def.tpl.json")}"

  vars {
    environment        = "${var.environment}"
    domain = "${var.domain}"
    weave_cidr = "${var.weave_cidr}"

    cloudwatch_log_group = "${aws_cloudwatch_log_group.rabbitmq.arn}"
    cloudwatch_region    = "${var.aws_region}"

    rabbitmq_user = "${var.rabbitmq_user}"
    rabbitmq_pass = "${var.rabbitmq_pass}"
  }
}

resource "aws_ecs_task_definition" "rabbitmq" {
  family                = "rabbitmq_${var.environment}"
  container_definitions = "${data.template_file.ecs_def_rabbitmq.rendered}"
}

resource "aws_ecs_service" "rabbitmq" {
  name            = "rabbitmq_${var.environment}"
  cluster         = "${var.cluster_name}"
  task_definition = "${aws_ecs_task_definition.rabbitmq.arn}"
  desired_count   = 1

  placement_strategy {
    type  = "spread"
    field = "attribute:ecs.availability-zone"
  }
}

#### Log Group for Privacy Negotiator rabbitmq
resource "aws_cloudwatch_log_group" "rabbitmq" {
  name = "${var.environment}.rabbitmq-container-logs"

  retention_in_days = 7

  tags {
    Name        = "rabbitmq"
    Environment = "${var.environment}"
  }
}
