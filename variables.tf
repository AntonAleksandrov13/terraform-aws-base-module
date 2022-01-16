### AWS IAM User section
variable "user_name" {
  type    = string
  default = "some_user_name"
  description = ""
}

variable "additional_policies_arn" {
  type    = list(string)
  default = []
}
### AWS IAM Role section
variable "create_base_role" {
  type    = bool
  default = false
}

variable "allow_user_assume_on_role" {
  type    = bool
  default = false
}

variable "role_name" {
  type    = string
  default = "terraform"
}

## AWS DynamoDB and S3 section
variable "generate_bucket_name" {
  type    = bool
  default = true
}
variable "state_bucket_name_override" {
  type    = string
  default = "my-very-unique-terraform-state-eu-central-1"
}

variable "s3_state_key_path" {
  type    = string
  default = "/*"
}
variable "generate_lock_table_name" {
  type    = bool
  default = true
}

variable "lock_table_name_override" {
  type    = string
  default = "terraform-state-lock"
}
