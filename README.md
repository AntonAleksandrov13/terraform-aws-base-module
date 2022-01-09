# terraform-aws-base-module

This repository contains a set of AWS resources required for state storing and locking. Using this module, you can create all necessary resources to get started with remote state in AWS S3.

## What does it deploy?

This module deploys:

1. AWS IAM role
2. AWS IAM policy for S3
   access [based on Terraform documentation](https://www.terraform.io/language/settings/backends/s3)
3. AWS IAM policy for DynamoDB
   access [based on Terraform documentation](https://www.terraform.io/language/settings/backends/s3)
4. S3 bucket for remote state storage
5. DynamoDB table for state locking

## Implementation

Run this module first to create all necessary resources for working with AWS and then proceed with your infrastructure.

```terraform
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

module "base-module" {
  source                    = "AntonAleksandrov13/base-module/aws"
  version                   = "1.0.0"
  create_base_role          = true
  # you need to create a user separately either in AWS console or using Terraform resources
  # using the following parameters, the user will be able to assume the newly created role
  # if you don't provide these params, only role will be created
  # please see examples folders for more
  allow_user_assume_on_role = true
  user_name                 = "existing_user_name"
}
```

## Additional configurations

If you want to give your role more permissions, you can use `additional_policies_arn` variable to attach more policies
to the role. This is useful in case you want to use this role for CI/CD.

Use the following configurations to enable this:

```terraform
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

module "base-module" {
  source                  = "AntonAleksandrov13/base-module/aws"
  version                 = "1.0.0"
  create_base_role        = true
  #...
  additional_policies_arn = [aws_iam_policy.test_policy.arn]
  #...
}
```



