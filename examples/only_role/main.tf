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

module "base_module" {
  source            = "../.."
  create_base_role  = true
}

output "role_name" {
  value = module.base_module.role_name
}

output "role_arn" {
  value = module.base_module.role_arn
}
