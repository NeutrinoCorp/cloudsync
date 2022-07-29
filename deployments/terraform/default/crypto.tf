data "aws_iam_policy_document" "kms_blob_bucket" {
  statement {
    sid    = "EnableRoot"
    effect = "Allow"
    principals {
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
      type        = "AWS"
    }
    actions   = ["kms:*"]
    resources = ["*"]
  }
  statement {
    sid    = "EnableCLIApp"
    effect = "Allow"
    principals {
      identifiers = [aws_iam_user.cli.arn]
      type        = "AWS"
    }
    actions = [
      "kms:GenerateDataKey",
      "kms:Decrypt"
    ]
    resources = ["*"]
  }
  version = "2012-10-17"
}

resource "aws_kms_key" "blob_bucket" {
  description             = "Key used to encrypt data within Neutrino CloudSync blob bucket"
  deletion_window_in_days = var.blob_bucket_encrypt_key_removal_days
  enable_key_rotation     = var.blob_bucket_encrypt_key_enable_rotation

  policy = data.aws_iam_policy_document.kms_blob_bucket.json
}

resource "aws_kms_alias" "blob_bucket" {
  target_key_id = aws_kms_key.blob_bucket.id
  name          = "alias/${local.app_name}"
}
