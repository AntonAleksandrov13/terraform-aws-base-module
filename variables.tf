### AWS IAM User section
variable "user_name" {
  type        = string
  default     = "some_user_name"
  description = "AWS IAM user name which can assume role. Does nothing without allow_user_assume_on_role=true."
}

variable "allow_user_assume_on_role" {
  type        = bool
  default     = false
  description = "Allows to an AWS IAM user to assume the newly created IAM role. See user_name variable to specify AWS IAM user."
}

variable "additional_policies_arn" {
  type        = list(string)
  default     = []
  description = "List of AWS IAM policy arns that will be attached to the newly created IAM role."
}
### AWS IAM Role section
variable "create_base_role" {
  type        = bool
  default     = false
  description = "Boolean determines whether to create a new IAM role. Note: only S3 and DynamoDB tables will be created in this case."
}

variable "role_name" {
  type        = string
  default     = "terraform"
  description = "The name of a new IAM role."
}

## AWS DynamoDB and S3 section
variable "generate_bucket_name" {
  type        = bool
  default     = true
  description = "Boolean determines whether to generate S3 bucket name. If enabled, S3 bucket will be named using the following pattern: terraform-{YOUR_ACCOUNT_NUMBER}."
}
variable "state_bucket_name_override" {
  type        = string
  default     = "my-very-unique-terraform-state-eu-central-1"
  description = "Overrides the generated S3 bucket name"
}

variable "s3_state_key_path" {
  type        = string
  default     = "/*"
  description = "S3 prefix used in IAM policy for S3 access. Determines which prefix can be read by AWS IAM policy for S3."
}
variable "generate_lock_table_name" {
  type        = bool
  default     = true
  description = "Boolean determines whether to generate DynamoDB table name. If enabled, the table will be named using the following pattern: terraform-state-lock-{YOUR_ACCOUNT_NUMBER}."
}

variable "lock_table_name_override" {
  type        = string
  default     = "terraform-state-lock"
  description = "Overrides the generated DynamoDB table name."
}
