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
## TODO: enable this for user and roles
variable "base_user_additional_policies_arn" {
  type    = list(string)
  default = []
}
### AWS IAM Role section
variable "create_base_role" {
  type = bool
  default = false
}

variable "allow_user_assume" {
  type = bool
  default = false
}

variable "base_role_name" {
  type    = string
  default = "terraform"
}

## AWS DynamoDB and S3 section
variable "generate_bucket_name" {
  type = bool
  default = true
}
variable "state_bucket_name_override" {
  type    = string
  default = "my-very-unique-terraform-state-eu-central-1"
}

variable "s3_state_key_path" {
  type    = string
  default = "/"
}
variable "generate_lock_table_name" {
  type = bool
  default = true
}

variable "lock_table_name_override" {
  type    = string
  default = "terraform-state-lock"
}
