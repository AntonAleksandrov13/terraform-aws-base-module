variable "region" {
  default = "eu-central-1"
}

variable "base_user_pgp_key" {
  default = "mDMEYbOgExYJKwYBBAHaRw8BAQdAzED..."
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
  base_user_pgp_key = var.base_user_pgp_key
}
