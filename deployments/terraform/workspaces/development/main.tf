terraform {
  backend "s3" {
    bucket         = "ncorp-cloudsync-tf-state"
    key            = "cloudsync/dev/terraform.tfstate"
    dynamodb_table = "ncorp-cloudsync-tf_state"
    region         = "us-east-2"
    encrypt        = true
  }
}

module "main" {
  source         = "../../default"
  aws_account    = var.aws_account
  aws_region     = var.aws_region
  aws_access_key = var.aws_access_key
  aws_secret_key = var.aws_secret_key
  environment    = "development"
}