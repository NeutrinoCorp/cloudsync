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

variable "org_short_name" {
  description = "Short name of the organization or individual (e.g. ncorp, aruiz) used to generate semantic infrastructure naming and/or labels"
  default     = "ncorp"
}

variable "blob_bucket_encrypt_key_removal_days" {
  description = "Days to wait for the blob bucket encryption key deletion"
  default     = 14
}

variable "blob_bucket_encrypt_key_enable_rotation" {
  description = "Enable key rotation of the blob bucket encryption key"
  default     = true
}

variable "blob_bucket_standard_ia_days" {
  description = "Days to move archives from Standard access to Standard Infrequent Access tier for blob bucket"
  type        = number
  default     = 30
}

variable "blob_bucket_intelligent_tier_days" {
  description = "Days to move archives from Standard Infrequent Access to Intelligent tiering for blob bucket"
  type        = number
  default     = 60
}

variable "blob_bucket_glacier_archive_days" {
  description = "Days to move archives from Intelligent tiering to Glacier Archive tier for blob bucket"
  type        = number
  default     = 180
}

variable "blob_bucket_glacier_deep_archive_days" {
  description = "Days to move archives from Glacier Archive tier to Glacier Deep Archive tier for blob bucket"
  type        = number
  default     = 365
}

variable "blob_bucket_standard_ia_versioned_days" {
  description = "Days to move archives from Standard access to Standard Infrequent Access tier for versioned objects in blob bucket"
  type        = number
  default     = 30
}

variable "blob_bucket_glacier_versioned_days" {
  description = "Days to move archives from Standard Infrequent Access tier to Glacier tier for versioned objects in blob bucket"
  type        = number
  default     = 60
}

variable "blob_bucket_expiration_versioned_days" {
  description = "Days to expire (permanently delete) versioned objects for versioned objects in blob bucket"
  type        = number
  default     = 90
}
