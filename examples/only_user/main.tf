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
  source            = "../.."
  create_base_user  = true
  base_user_pgp_key = "ACTUAL_BASE64_PGP_KEY_SHOULD_BE_HERE"
}
