provider "aws" {
  region = "${var.aws_region}"
}

terraform {
  backend "s3" {
    encrypt = true
    bucket  = "privneg-terraform"
    key     = "cluster/terraform.tfstate"
    region  = "eu-west-1"
  }
}

## AWS ECS - Elastic Container Service
resource "aws_ecs_cluster" "app" {
  name = "${var.cluster_name}"
}

## Network
data "aws_availability_zones" "available" {}

resource "aws_vpc" "app" {
  cidr_block = "10.0.0.0/16"

  enable_dns_support   = true
  enable_dns_hostnames = true

  tags {
    cluster = "${var.cluster_name}"
    Name    = "${var.cluster_name}.ecs"
  }
}

resource "aws_subnet" "app" {
  count             = "3"
  cidr_block        = "${cidrsubnet(aws_vpc.app.cidr_block, 8, count.index)}"
  availability_zone = "${data.aws_availability_zones.available.names[count.index]}"
  vpc_id            = "${aws_vpc.app.id}"

  tags {
    cluster = "${var.cluster_name}"
    Name    = "${var.cluster_name}.ecs.${count.index}"
  }
}

output "subnets" {
  value = ["${aws_subnet.app.*.id}"]
}

resource "aws_internet_gateway" "app" {
  vpc_id = "${aws_vpc.app.id}"
}

resource "aws_route_table" "app" {
  vpc_id = "${aws_vpc.app.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.app.id}"
  }

  tags {
    cluster = "${var.cluster_name}"
    Name    = "${var.cluster_name}.ecs"
  }
}

resource "aws_route_table_association" "a" {
  count          = "3"
  subnet_id      = "${element(aws_subnet.app.*.id, count.index)}"
  route_table_id = "${aws_route_table.app.id}"
}

## Compute
resource "aws_spot_fleet_request" "cheap_compute" {
  iam_fleet_role                      = "${aws_iam_role.fleet.arn}"
  spot_price                          = "1.00"
  allocation_strategy                 = "diversified"
  target_capacity                     = 2
  valid_until                         = "2019-11-04T20:44:20Z"
  terminate_instances_with_expiration = true

  launch_specification {
    instance_type               = "m4.large"
    ami                         = "ami-13c8f475"
    spot_price                  = "0.03"
    user_data                   = "${data.template_file.user_data.rendered}"
    iam_instance_profile        = "${aws_iam_instance_profile.app.name}"
    availability_zone           = "eu-west-1a"
    subnet_id                   = "${aws_subnet.app.0.id}"
    vpc_security_group_ids      = ["${aws_security_group.instance_sg.id}"]
    key_name                    = "test"
    associate_public_ip_address = true
  }

  launch_specification {
    instance_type               = "m4.large"
    ami                         = "ami-13c8f475"
    spot_price                  = "0.03"
    user_data                   = "${data.template_file.user_data.rendered}"
    iam_instance_profile        = "${aws_iam_instance_profile.app.name}"
    availability_zone           = "eu-west-1b"
    subnet_id                   = "${aws_subnet.app.1.id}"
    vpc_security_group_ids      = ["${aws_security_group.instance_sg.id}"]
    key_name                    = "test"
    associate_public_ip_address = true
  }

  launch_specification {
    instance_type               = "m4.large"
    ami                         = "ami-13c8f475"
    spot_price                  = "0.03"
    user_data                   = "${data.template_file.user_data.rendered}"
    iam_instance_profile        = "${aws_iam_instance_profile.app.name}"
    availability_zone           = "eu-west-1c"
    subnet_id                   = "${aws_subnet.app.2.id}"
    vpc_security_group_ids      = ["${aws_security_group.instance_sg.id}"]
    key_name                    = "test"
    associate_public_ip_address = true
  }
}

data "template_file" "user_data" {
  template = "${file("user-data.sh")}"

  vars {
    aws_region        = "${var.aws_region}"
    ecs_cluster_name  = "${aws_ecs_cluster.app.name}"
    ecs_log_level     = "info"
    ecs_agent_version = "latest"
  }
}

resource "aws_iam_instance_profile" "app" {
  name = "${var.cluster_name}.ecs-instprofile"
  role = "${aws_iam_role.app_instance.name}"
}

// TODO: Reduce AWS Policy to follow Principal of Least Privilege!
resource "aws_iam_role_policy" "app_instance" {
  name = "${var.cluster_name}.app_policy"
  role = "${aws_iam_role.app_instance.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:*",
        "ecs:*",
        "autoscaling:DescribeAutoScalingInstances",
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "logs:CreateLogStream",
        "logs:PutLogEvents",
        "route53:*",
        "route53domains:*"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "fleet" {
  name = "${var.cluster_name}.fleet_policy"
  role = "${aws_iam_role.fleet.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
       "ec2:DescribeImages",
       "ec2:DescribeSubnets",
       "ec2:RequestSpotInstances",
       "ec2:TerminateInstances",
       "ec2:DescribeInstanceStatus",
       "iam:PassRole"
        ],
    "Resource": ["*"]
  }]
}
EOF
}

resource "aws_iam_role" "fleet" {
  name = "${var.cluster_name}-fleet"

  assume_role_policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "spotfleet.amazonaws.com",
          "ec2.amazonaws.com"
        ]
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role" "app_instance" {
  name = "${var.cluster_name}.ecs-instance-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.amazonaws.com",
          "spotfleet.amazonaws.com"
        ]
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

## Security Groups
resource "aws_security_group" "instance_sg" {
  description = "Controls access to nodes in ECS cluster"
  vpc_id      = "${aws_vpc.app.id}"
  name        = "${var.cluster_name}.ecs-sg"

  ingress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["10.0.0.0/16"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    cluster = "${var.cluster_name}"
    Name    = "${var.cluster_name}.ecs.instance"
  }
}
