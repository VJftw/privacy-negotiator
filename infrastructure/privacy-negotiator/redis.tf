### ECS RabbitMQ containers
data "template_file" "ecs_def_redis" {
  template = "${file("${path.module}/redis.def.tpl.json")}"

  vars {
    environment        = "${var.environment}"
    domain = "${var.domain}"
    weave_cidr = "${var.weave_cidr}"

    cloudwatch_log_group = "${aws_cloudwatch_log_group.redis.arn}"
    cloudwatch_region    = "${var.aws_region}"
  }
}

resource "aws_ecs_task_definition" "redis" {
  family                = "redis_${var.environment}"
  container_definitions = "${data.template_file.ecs_def_redis.rendered}"
}

resource "aws_ecs_service" "redis" {
  name            = "redis_${var.environment}"
  cluster         = "${var.cluster_name}"
  task_definition = "${aws_ecs_task_definition.redis.arn}"
  desired_count   = 1

  placement_strategy {
    type  = "spread"
    field = "attribute:ecs.availability-zone"
  }
}

#### Log Group for Privacy Negotiator redis
resource "aws_cloudwatch_log_group" "redis" {
  name = "${var.environment}.redis-container-logs"

  retention_in_days = 7

  tags {
    Name        = "redis"
    Environment = "${var.environment}"
  }
}
