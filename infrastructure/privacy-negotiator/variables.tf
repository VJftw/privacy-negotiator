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

variable "weave_cidr" {
  type = "string"
  description = "The Weave subnet to join. This should be unique across applications/environments"
}

variable "version" {
  type        = "string"
  description = "The version of the application(containers) to deploy"
}

variable "jwt_secret" {
  type = "string"
  description = "The JWT secret"
}

variable "rabbitmq_user" {
  type = "string"
  description = "The username for RabbitMQ"
}

variable "rabbitmq_pass" {
  type = "string"
  description = "The password for RabbitMQ"
}

variable "postgres_user" {
  type = "string"
  description = "The host for Postgres"
}

variable "postgres_password" {
  type = "string"
  description = "The password for Postgres"
}

variable "facebook_app_id" {
  type = "string"
  description = "The ID for the Facebook App"
}

variable "facebook_app_secret" {
  type = "string"
  description = "The secret for the Facebook App"
}

variable "aws_availability_zones" {
  default = "list"
}
