variable "environment" {
  type        = "string"
  description = "The stage of development this service should be tagged as"
}

variable "cluster_name" {
  type        = "string"
  description = "The cluster to launch this service into"
}

variable "aws_region" {
  type        = "string"
  description = "The AWS region"
}

variable "domain" {
  type        = "string"
  description = "The domain name for the service (environment)"
}

// variable "container_version" {
//   type        = "string"
//   description = "The version of the container to deploy"
// }

variable "weave_cidr" {
  type = "string"
  description = "The Weave subnet to join. This should be unique across applications/environments"
}

variable "aws_availability_zones" {
  default = "list"
}
