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
