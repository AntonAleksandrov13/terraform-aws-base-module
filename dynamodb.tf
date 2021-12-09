resource "aws_dynamodb_table" "terraform_lock_table" {
   name = var.terraform_lock_table_name
   hash_key = "LockID"
   read_capacity = 20
   write_capacity = 20

   attribute {
      name = "LockID"
      type = "S"
   }

 tags = {
     Name = var.terraform_lock_table_name
   }
}
