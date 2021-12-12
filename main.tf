terraform {
  required_providers {
    aws        = ">= 3.22.0"
  }
}

data "aws_caller_identity" "current" {}