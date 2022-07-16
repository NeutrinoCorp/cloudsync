variable "aws_account" {
  description = "Amazon Web Services account ID to be used for infrastructure provisioning"
  type        = string
}

variable "aws_region" {
  description = "Amazon Web Services region to be used for infrastructure provisioning"
  type        = string
}

variable "aws_access_key" {
  description = "Amazon Web Services access key to be able to perform provisioning actions"
  type        = string
}

variable "aws_secret_key" {
  description = "Amazon Web Services secret access key to be able to perform provisioning actions"
  type        = string
}

variable "environment" {
  description = "Stage (e.g. development, production) to be used"
  type        = string
  default     = "development"
}