locals {
  table_name = var.generate_lock_table_name ? "terraform-state-lock-${data.aws_caller_identity.current.account_id}" : var.lock_table_name_override
}
resource "aws_dynamodb_table" "terraform_lock" {
   name = local.table_name
   hash_key = "LockID"
   read_capacity = 20
   write_capacity = 20

   attribute {
      name = "LockID"
      type = "S"
   }

 tags = {
     Name = local.table_name
   }
}
