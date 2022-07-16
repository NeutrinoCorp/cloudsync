locals {
  env_short_names = {
    "production" : "prod"
    "development" : "dev"
  }

  app_name = "ncorp-${local.env_short_names[var.environment]}-cloudsync"
}
