output "user_name" {
  value = var.create_base_user == true ? aws_iam_user.base_user[0].name : null
}

output "aws_access_key_id" {
  value       = var.create_base_user == true ? aws_iam_access_key.base_user[0].id : null
  description = "Access key ID"
  sensitive   = true
}

output "base64_aws_secret_access_key" {
  sensitive   = true
  value       = var.create_base_user == true ? aws_iam_access_key.base_user[0].encrypted_secret : null
  description = "Encrypted secret, base64 encoded, if base_user_pgp_key was specified. This attribute is not available for imported resources."
}

output "ses_smtp_base64_aws_secret_access_key" {
  sensitive   = true
  value       = var.create_base_user == true ? aws_iam_access_key.base_user[0].encrypted_ses_smtp_password_v4 : null
  description = "Encrypted SES SMTP password, base64 encoded, if base_user_pgp_key was specified. This attribute is not available for imported resources."
}

output "ses_smtp_password_aws_secret_access_key" {
  value       = var.create_base_user == true ? aws_iam_access_key.base_user[0].ses_smtp_password_v4 : null
  description = "Secret access key converted into an SES SMTP password by applying AWS's documented Sigv4 conversion algorithm."
  sensitive   = true
}

output "role_name" {
  value = var.create_base_role == true ? aws_iam_role.base_role[0].name : null
}

output "role_arn" {
  value = var.create_base_role == true ? aws_iam_role.base_role[0].arn : null
}

output "s3_bucket_name" {
  value = aws_s3_bucket.state_storage.id
}

output "lock_table_name" {
  value = aws_dynamodb_table.terraform_lock.id
}