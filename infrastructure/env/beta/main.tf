module "main" {
  source = "../../privacy-negotiator"

  aws_region = "eu-west-1"
  aws_availability_zones = "eu-west-1a,eu-west-1b,eu-west-1c"

  cluster_name = "privacy-negotiator"

  environment = "beta"
  domain      = "beta.privacy-negotiator.vjpatel.me"

  container_version = "${var.container_version}"
  weave_cidr = "10.32.101.0/24"

  facebook_app_id = "${var.facebook_app_id}"
  // facebook_app_secret = "${var.facebook_app_secret}"

}
