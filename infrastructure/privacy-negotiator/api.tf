### ECS Privacy Negotiator api containers
data "template_file" "ecs_def_api" {
  template = "${file("${path.module}/api.def.tpl.json")}"

  vars {
    environment = "${var.environment}"
    domain      = "${var.domain}"
    version     = "${var.version}"
    weave_cidr  = "${var.weave_cidr}"

    cloudwatch_log_group = "${aws_cloudwatch_log_group.api.arn}"
    cloudwatch_region    = "${var.aws_region}"

    jwt_secret        = "${var.jwt_secret}"
    rabbitmq_user     = "${var.rabbitmq_user}"
    rabbitmq_pass     = "${var.rabbitmq_pass}"
    rabbitmq_hostname = "rabbitmq-${var.environment}.weave.local"
    redis_host        = "redis-${var.environment}.weave.local"
  }
}

resource "aws_ecs_task_definition" "api" {
  family                = "api_${var.environment}"
  container_definitions = "${data.template_file.ecs_def_api.rendered}"
}

resource "aws_ecs_service" "api" {
  name                               = "api_${var.environment}"
  cluster                            = "${var.cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.api.arn}"
  desired_count                      = 4
  iam_role                           = "${aws_iam_role.ecs_service.arn}"
  deployment_minimum_healthy_percent = 50

  placement_strategy {
    type  = "spread"
    field = "attribute:ecs.availability-zone"
  }

  load_balancer {
    target_group_arn = "${aws_alb_target_group.api.id}"
    container_name   = "api_${var.environment}"
    container_port   = "80"
  }

  depends_on = [
    "aws_iam_role_policy.ecs_service",
    "aws_alb_listener.api",
  ]
}

#### Log Group for Privacy Negotiator API
resource "aws_cloudwatch_log_group" "api" {
  name = "${var.environment}.api-container-logs"

  retention_in_days = 7

  tags {
    Name        = "API"
    Environment = "${var.environment}"
  }
}

resource "aws_route53_record" "domain" {
  zone_id = "${data.aws_route53_zone.organisation.zone_id}"
  name    = "api.${var.domain}"
  type    = "A"

  alias {
    name                   = "${aws_alb.api.dns_name}"
    zone_id                = "${aws_alb.api.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_alb" "api" {
  name            = "api-${var.environment}-alb"
  subnets         = ["${data.aws_subnet.app_cluster.*.id}"]
  security_groups = ["${aws_security_group.alb_sg.id}"]

  provisioner "local-exec" {
    command = "sleep 10"
  }
}

resource "aws_alb_target_group" "api" {
  name     = "api-${var.environment}-tg"
  port     = 80
  protocol = "HTTP"
  vpc_id   = "${data.aws_vpc.app_cluster.id}"

  health_check {
    path                = "/v1/health"
    healthy_threshold   = 5
    unhealthy_threshold = 2
    timeout             = 5
    protocol            = "HTTP"
    interval            = 30
    matcher             = "200"
  }
}

resource "aws_security_group" "alb_sg" {
  description = "Controls access to and from the ALB"

  vpc_id = "${data.aws_vpc.app_cluster.id}"
  name   = "api.${var.environment}.alb-sg"

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 80
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    protocol    = "tcp"
    from_port   = 443
    to_port     = 443
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"

    cidr_blocks = [
      "0.0.0.0/0",
    ]
  }
}

resource "aws_iam_role" "ecs_service" {
  name = "api.${var.environment}.ecs_role"

  assume_role_policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "ecs_service" {
  name = "api.${var.environment}.ecs_policy"
  role = "${aws_iam_role.ecs_service.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:Describe*",
        "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
        "elasticloadbalancing:DeregisterTargets",
        "elasticloadbalancing:Describe*",
        "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
        "elasticloadbalancing:RegisterTargets"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_alb_listener" "api" {
  load_balancer_arn = "${aws_alb.api.id}"
  port              = "80"
  protocol          = "HTTP"

  default_action {
    target_group_arn = "${aws_alb_target_group.api.id}"
    type             = "forward"
  }
}

resource "aws_alb_listener" "front_end_ssl" {
  load_balancer_arn = "${aws_alb.api.id}"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2015-05"
  certificate_arn   = "${data.aws_acm_certificate.api.arn}"

  default_action {
    target_group_arn = "${aws_alb_target_group.api.id}"
    type             = "forward"
  }
}

data "aws_acm_certificate" "api" {
  domain   = "api.${var.domain}"
  statuses = ["ISSUED"]
}
