resource "aws_kms_key" "blob_bucket" {
  description             = "Key used to encrypt data within Neutrino CloudSync blob bucket"
  deletion_window_in_days = var.blob_bucket_encrypt_key_removal_days
}
