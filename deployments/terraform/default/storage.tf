resource "aws_s3_bucket" "main_blob" {
  bucket = local.app_name
}