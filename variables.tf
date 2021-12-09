
variable "region" {
  type = string
  default = "eu-central-1"
}

variable "base_user_name" {
  type=string
  default = "terraform"
}

variable "base_user_path" {
  type=string
  default = "/"
}

variable "base_user_pgp_key" {
  type=string
  default = "keybase:some_person_that_exists"
}

variable "base_user_additional_policies_arn"{
    type=list(string)
    default = []
}

variable "state_bucket_name"{
    type = string
    default="terraform-state-eu-central-1"
}

variable "s3_state_key_path"{
    type = string
    default = "/"
}

variable "terraform_lock_table_name"{
    type = string
    default = "terraform-state-lock"
}
