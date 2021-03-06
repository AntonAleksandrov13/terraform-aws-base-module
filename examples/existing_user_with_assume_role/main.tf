variable "region" {
  default = "eu-central-1"
}

provider "aws" {
  region = var.region
}

terraform {
  required_providers {
    aws = ">= 3.65.0"
  }
}
# usually you would do something like this:
#resource "aws_iam_user" "test_user" {
#  name = "tester"
#  path = "/"
#}

# for tests we will get current AWS user, so in tests in can test role assume
data "aws_caller_identity" "current" {}

locals {
  # get only user name from user arn
  current_user = regex("([^/]+$)", data.aws_caller_identity.current.arn)[0]
}
module "base_module" {
  source            = "../.."
  # creates a role, permissions, S3 and DynamoDB
  create_base_role  = true
  # you need to create a user separately either in AWS console or using Terraform resources
  # using the following parameters, the user will be able to assume the newly created role
  # if you don't provide these params, only role will be created
  # please see examples folders for more
  allow_user_assume_on_role = true
  user_name    = local.current_user
}

output "role_name" {
  value = module.base_module.role_name
}

output "role_arn" {
  value = module.base_module.role_arn
}

output "s3_bucket_name" {
  value = module.base_module.s3_bucket_name
}

output "lock_table_name" {
  value = module.base_module.lock_table_name
}
