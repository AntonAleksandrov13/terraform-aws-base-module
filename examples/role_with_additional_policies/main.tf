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

# for test we will get current AWS user, so in tests in can test role assume
data "aws_caller_identity" "current" {}

locals {
  # get only user name from user arn
  current_user = regex("([^/]+$)", data.aws_caller_identity.current.arn)[0]
}

resource "aws_iam_policy" "test_policy" {
  name        = "test_policy"
  path        = "/"
  description = "My test policy"
  policy      = jsonencode({
    Version   = "2012-10-17"
    Statement = [
      {
        Action   = [
          "ec2:Describe*",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}

resource "aws_iam_policy" "another_test_policy" {
  name        = "another_test_policy"
  path        = "/"
  description = "My test policy"

  policy = jsonencode({
    Version   = "2012-10-17"
    Statement = [
      {
        Action   = [
          "ec2:Describe*",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}
module "base_module" {
  source                    = "../.."
  create_base_role          = true
  allow_user_assume_on_role = true
  user_name                 = local.current_user
  additional_policies_arn   = [aws_iam_policy.test_policy.arn, aws_iam_policy.another_test_policy.arn]
}

output "role_name" {
  value = module.base_module.role_name
}

output "role_arn" {
  value = module.base_module.role_arn
}

output "test_policy_name" {
  value = aws_iam_policy.test_policy.name
}

output "another_test_policy_name" {
  value = aws_iam_policy.another_test_policy.name
}