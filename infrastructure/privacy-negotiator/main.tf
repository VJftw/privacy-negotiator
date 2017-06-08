provider "aws" {
  region = "${var.aws_region}"
}

// data "aws_vpc" "app_cluster" {
//   tags {
//     cluster = "${var.cluster_name}"
//   }
// }
//
// data "aws_subnet" "app_cluster" {
//   vpc_id            = "${data.aws_vpc.app_cluster.id}"
//   state             = "available"
//   availability_zone = "${element(split(",", var.aws_availability_zones), count.index)}"
//
//   tags {
//     cluster = "${var.cluster_name}"
//   }
//
//   count = "${length(split(",", var.aws_availability_zones))}"
// }

data "aws_route53_zone" "organisation" {
  name = "vjpatel.me."
}
