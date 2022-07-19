locals {
  env_short_names = {
    "production" : "prod"
    "development" : "dev"
  }

  app_name = local.env_short_names[var.environment] == "" ? "${var.org_short_name}-cloudsync" : "${var.org_short_name}-${local.env_short_names[var.environment]}-cloudsync"
}
