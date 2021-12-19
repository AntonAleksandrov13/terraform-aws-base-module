locals {
  bucket_name = var.generate_bucket_name ? "terraform-${data.aws_caller_identity.current.account_id}" : var.state_bucket_name_override
}

resource "aws_s3_bucket" "state_storage" {
  bucket = local.bucket_name
  acl    = "private"

  tags = {
    Name = local.bucket_name
  }
}

resource "aws_s3_bucket_public_access_block" "state_storage" {
  bucket = aws_s3_bucket.state_storage.id

  block_public_acls   = true
  block_public_policy = true
}
