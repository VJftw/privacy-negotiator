terraform {
  backend "s3" {
    encrypt = true
    bucket  = "privneg-terraform"
    key     = "beta/terraform.tfstate"
    region  = "eu-west-1"
  }
}

module "main" {
  source = "../../privacy-negotiator"

  aws_region = "eu-west-1"
  aws_availability_zones = "eu-west-1a,eu-west-1b,eu-west-1c"

  cluster_name = "privacy-negotiator"

  environment = "beta"
  domain      = "beta.privacy-negotiator.vjpatel.me"

  version = "${var.version}"
  weave_cidr = "10.32.101.0/24"

  jwt_secret = "${var.jwt_secret}"
  rabbitmq_user = "${var.rabbitmq_user}"
  rabbitmq_pass = "${var.rabbitmq_pass}"

  postgres_user = "${var.postgres_user}"
  postgres_password = "${var.postgres_password}"

  facebook_app_id = "${var.facebook_app_id}"
  facebook_app_secret = "${var.facebook_app_secret}"
}
