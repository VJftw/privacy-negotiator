provider "aws" {
  region = "${var.aws_region}"
}

data "aws_vpc" "app_cluster" {
  tags {
    cluster = "${var.cluster_name}"
  }
}

data "aws_subnet" "app_cluster" {
  vpc_id            = "${data.aws_vpc.app_cluster.id}"
  state             = "available"
  availability_zone = "${element(split(",", var.aws_availability_zones), count.index)}"

  tags {
    cluster = "${var.cluster_name}"
  }

  count = "${length(split(",", var.aws_availability_zones))}"
}

data "aws_route53_zone" "organisation" {
  name = "vjpatel.me."
}

### ECS PrivNeg Web containers
data "template_file" "ecs_def_web" {
  template = "${file("${path.module}/web.def.tpl.json")}"

  vars {
    fb_app_id    = "${var.facebook_app_id}"
    api_endpoint = "${var.api_endpoint}"

    version = "${var.version}"

    cloudwatch_log_group = "${aws_cloudwatch_log_group.web.arn}"
    cloudwatch_region    = "${var.aws_region}"

    weave_cidr = "${var.weave_cidr}"
  }
}

resource "aws_ecs_task_definition" "web" {
  family                = "web_${var.environment}"
  container_definitions = "${data.template_file.ecs_def_web.rendered}"
}

resource "aws_ecs_service" "chat_chat" {
  name            = "web_${var.environment}"
  cluster         = "${var.cluster_name}"
  task_definition = "${aws_ecs_task_definition.web.arn}"
  desired_count   = "2"
  iam_role        = "${aws_iam_role.ecs_service.arn}"

  placement_strategy {
    type  = "spread"
    field = "attribute:ecs.availability-zone"
  }

  load_balancer {
    target_group_arn = "${aws_alb_target_group.web.id}"
    container_name   = "web_${var.environment}"
    container_port   = "80"
  }

  depends_on = [
    "aws_iam_role_policy.ecs_service",
    "aws_alb_listener.web",
  ]
}

#### Log Group for Web
resource "aws_cloudwatch_log_group" "web" {
  name = "${var.environment}.web-container-logs"

  retention_in_days = 7

  tags {
    Name        = "Web"
    Environment = "${var.environment}"
  }
}

resource "aws_route53_record" "domain" {
  zone_id = "${data.aws_route53_zone.organisation.zone_id}"
  name    = "${var.domain}"
  type    = "A"

  alias {
    name                   = "${aws_alb.web.dns_name}"
    zone_id                = "${aws_alb.web.zone_id}"
    evaluate_target_health = false
  }
}

resource "aws_alb" "web" {
  name            = "web-${var.environment}-alb"
  subnets         = ["${data.aws_subnet.app_cluster.*.id}"]
  security_groups = ["${aws_security_group.alb_sg.id}"]

  provisioner "local-exec" {
    command = "sleep 10"
  }
}

resource "aws_alb_target_group" "web" {
  name     = "web-${var.environment}-tg"
  port     = 80
  protocol = "HTTP"
  vpc_id   = "${data.aws_vpc.app_cluster.id}"

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 3
    timeout             = 3
    protocol            = "HTTP"
    interval            = 5
    matcher             = "200,404"
  }
}

resource "aws_security_group" "alb_sg" {
  description = "Controls access to and from the ALB"

  vpc_id = "${data.aws_vpc.app_cluster.id}"
  name   = "web.${var.environment}.alb-sg"

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
  name = "web.${var.environment}.ecs_role"

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
  name = "web.${var.environment}.ecs_policy"
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

data "aws_acm_certificate" "privneg" {
  domain   = "privacy-negotiator.vjpatel.me"
  statuses = ["ISSUED"]
}

resource "aws_alb_listener" "web" {
  load_balancer_arn = "${aws_alb.web.id}"
  port              = "80"
  protocol          = "HTTP"

  default_action {
    target_group_arn = "${aws_alb_target_group.web.id}"
    type             = "forward"
  }
}

resource "aws_alb_listener" "web_ssl" {
  load_balancer_arn = "${aws_alb.web.id}"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2015-05"
  certificate_arn   = "${data.aws_acm_certificate.privneg.arn}"

  default_action {
    target_group_arn = "${aws_alb_target_group.web.id}"
    type             = "forward"
  }
}
