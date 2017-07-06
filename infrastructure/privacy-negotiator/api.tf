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
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 2
    protocol            = "HTTP"
    interval            = 10
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

# Cloudwatch Alarms
resource "aws_cloudwatch_log_metric_filter" "api_error" {

  name = "${var.environment}-api.error"
  pattern = "error"
  log_group_name = "${aws_cloudwatch_log_group.api.name}"

  metric_transformation {
    name = "${var.environment}-api.error"
    namespace = "${var.environment}-api"
    value = "1"
  }

}

resource "aws_cloudwatch_log_metric_filter" "api_error_reset" {

  name = "${var.environment}-api.error"
  pattern = ""
  log_group_name = "${aws_cloudwatch_log_group.api.name}"

  metric_transformation {
    name = "${var.environment}-api.error"
    namespace = "${var.environment}-api"
    value = "0"
  }

}

resource "aws_cloudwatch_metric_alarm" "api_error" {

  alarm_name = "${var.environment}.api.error"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  threshold = "1"
  period = "60"
  statistic = "Sum"
  evaluation_periods = "1"
  metric_name = "${var.environment}.api.error"
  namespace = "${var.environment}-api.error"
  alarm_description = "monitors log for api errors"
//  alarm_actions = ["arn:aws:sns:eu-west-1:812414252941:error_notification"]

}

# Autoscaling
resource "aws_iam_role" "api_autoscaling" {
  name = "api.${var.environment}.ecs_autoscaling"

  assume_role_policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "application-autoscaling.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    },
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

resource "aws_iam_role_policy_attachment" "test-attach" {
  role       = "${aws_iam_role.api_autoscaling.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceAutoscaleRole"
}

resource "aws_appautoscaling_target" "api" {
  service_namespace = "ecs"
  resource_id = "service/${var.cluster_name}/${aws_ecs_service.api.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  role_arn = "${aws_iam_role.api_autoscaling.arn}"
  min_capacity = 2
  max_capacity = 10
}

resource "aws_appautoscaling_policy" "api_up" {
  name = "scale-up"
  service_namespace = "ecs"
  resource_id = "service/${var.cluster_name}/${aws_ecs_service.api.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  adjustment_type = "ChangeInCapacity"
  cooldown = 60
  metric_aggregation_type = "Maximum"

  step_adjustment {
    metric_interval_lower_bound = 0
    scaling_adjustment = 1
  }

  depends_on = ["aws_appautoscaling_target.api"]
}

resource "aws_appautoscaling_policy" "api_down" {
  name = "scale-down"
  service_namespace = "ecs"
  resource_id = "service/${var.cluster_name}/${aws_ecs_service.api.name}"
  scalable_dimension = "ecs:service:DesiredCount"

  adjustment_type = "ChangeInCapacity"
  cooldown = 60
  metric_aggregation_type = "Maximum"

  step_adjustment {
    metric_interval_lower_bound = 0
    scaling_adjustment = -1
  }

  depends_on = ["aws_appautoscaling_target.api"]
}

resource "aws_cloudwatch_metric_alarm" "api_cpu_high" {
  alarm_name = "${var.environment}.api-cpuutilization-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods = "2"
  metric_name = "CPUUtilization"
  namespace = "AWS/ECS"
  period = "60"
  statistic = "Maximum"
  threshold = "85"

  dimensions {
    ClusterName = "${var.cluster_name}"
    ServiceName = "${aws_ecs_service.api.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.api_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "api_cpu_low" {
  alarm_name = "${var.environment}.api-cpuutilization-low"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods = "2"
  metric_name = "CPUUtilization"
  namespace = "AWS/ECS"
  period = "60"
  statistic = "Maximum"
  threshold = "30"

  dimensions {
    ClusterName = "${var.cluster_name}"
    ServiceName = "${aws_ecs_service.api.name}"
  }

  ok_actions = ["${aws_appautoscaling_policy.api_down.arn}"]
}