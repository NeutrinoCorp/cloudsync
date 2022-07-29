provider "aws" {
  region     = var.aws_region
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key

  default_tags {
    tags = {
      platform    = "cloudsync"
      environment = var.environment
    }
  }
}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}
