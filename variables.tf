
variable "region" {
  type    = string
  default = "eu-central-1"
}
### AWS IAM User section
variable "create_base_user" {
  type = bool
  default = false
}

variable "base_user_name" {
  type    = string
  default = "terraform"
}

variable "base_user_path" {
  type    = string
  default = "/"
}

variable "base_user_pgp_key" {
  type     = string
  default = "place_holder"
}

## not used right now
## todo: enable this for user and roles
variable "base_user_additional_policies_arn" {
  type    = list(string)
  default = []
}
### AWS IAM Role section
variable "create_base_role" {
  type = bool
  default = true
}
variable "state_bucket_name" {
  type    = string
  default = "my-very-unique-terraform-state-eu-central-1"
}

variable "s3_state_key_path" {
  type    = string
  default = "/"
}

variable "terraform_lock_table_name" {
  type    = string
  default = "terraform-state-lock"
}
