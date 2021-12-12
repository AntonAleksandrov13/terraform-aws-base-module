variable "region" {
  default = "eu-central-1"
}

provider "aws" {
  region = var.region
}

terraform {
  required_providers {
    aws = ">= 3.22.0"
  }
}

resource "aws_iam_user" "test_user" {
  name = "tester"
  path = "/"
}

module "base-module" {
  source            = "../.."
  create_base_role  = true
  allow_user_assume = true
  base_user_name = aws_iam_user.test_user.name
}
