resource "aws_s3_bucket" "main_blob" {
  bucket = local.app_name
}

resource "aws_s3_bucket_acl" "main_blob" {
  depends_on = [aws_s3_bucket.main_blob]

  bucket = local.app_name
  acl    = "private"
}

resource "aws_s3_bucket_versioning" "main_blob" {
  depends_on = [aws_s3_bucket.main_blob]

  bucket = local.app_name
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "main_blob" {
  depends_on = [aws_s3_bucket.main_blob]

  bucket = local.app_name
  rule {
    id     = "tiering"
    status = "Enabled"

    transition {
      storage_class = "STANDARD_IA"
      days          = var.blob_bucket_standard_ia_days
    }

    transition {
      storage_class = "INTELLIGENT_TIERING"
      days          = var.blob_bucket_intelligent_tier_days
    }
  }

  rule {
    id     = "expiration"
    status = "Enabled"

    noncurrent_version_transition {
      noncurrent_days = var.blob_bucket_standard_ia_versioned_days
      storage_class   = "STANDARD_IA"
    }

    noncurrent_version_transition {
      noncurrent_days = var.blob_bucket_glacier_versioned_days
      storage_class   = "GLACIER"
    }

    noncurrent_version_expiration {
      noncurrent_days = var.blob_bucket_expiration_versioned_days
    }
  }
}

resource "aws_s3_bucket_intelligent_tiering_configuration" "main_blob" {
  depends_on = [aws_s3_bucket.main_blob]

  bucket = local.app_name
  name   = "EntireBucket"
  tiering {
    access_tier = "ARCHIVE_ACCESS"
    days        = var.blob_bucket_glacier_archive_days
  }
  tiering {
    access_tier = "DEEP_ARCHIVE_ACCESS"
    days        = var.blob_bucket_glacier_deep_archive_days
  }
}
