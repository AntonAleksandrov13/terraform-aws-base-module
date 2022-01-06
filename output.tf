
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