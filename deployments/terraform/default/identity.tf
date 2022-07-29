data "aws_iam_policy_document" "cli" {
  version = "2012-10-17"
  statement {
    sid     = "BlobStorageAccess"
    effect  = "Allow"
    actions = [
      "s3:PutObject",
      "s3:AbortMultipartUpload",
      "s3:ListBucket",
      "s3:ListMultipartUploadParts",
      "s3:GetObject"
    ]
    resources = [
      "arn:aws:s3:::${aws_s3_bucket.main_blob.bucket}",
      "arn:aws:s3:::${aws_s3_bucket.main_blob.bucket}/*"
    ]
  }
}

resource "aws_iam_user" "cli" {
  name = "${local.app_name}-CLIAccess"
}

resource "aws_iam_user_policy" "cli" {
  name   = aws_iam_user.cli.name
  policy = data.aws_iam_policy_document.cli.json
  user   = aws_iam_user.cli.name
}
