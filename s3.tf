resource "aws_s3_bucket" "state-bucket" {
  bucket = var.state_bucket_name
  acl    = "private"

  tags = {
    Name        = var.state_bucket_name
  }
}