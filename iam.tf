locals {
  user_count = var.create_base_user ? 1 : 0
}

resource "aws_iam_user" "base_user" {
  count = local.user_count
  name  = var.base_user_name
  path  = var.base_user_path
}

resource "aws_iam_access_key" "base_user" {
  count   = local.user_count
  user    = aws_iam_user.base_user[count.index].name
  pgp_key = var.base_user_pgp_key
}

# resource "aws_iam_user_policy" "s3_access" {
#   name = "S3AccessTerraform"
#   user = aws_iam_user.base_user.name

#   policy = <<EOF
# {
#   "Version": "2012-10-17",
#   "Statement": [
#     {
#       "Effect": "Allow",
#       "Action": "s3:ListBucket",
#       "Resource": "arn:aws:s3:::${var.state_bucket_name}"
#     },
#     {
#       "Effect": "Allow",
#       "Action": ["s3:GetObject", "s3:PutObject", "s3:DeleteObject"],
#       "Resource": "arn:aws:s3:::${var.state_bucket_name}${var.s3_state_key_path}"
#     }
#   ]
# }
# EOF
# }

# resource "aws_iam_user_policy" "dynamodb_access" {
#   name = "DynamoDBAccessTerraform"
#   user = aws_iam_user.base_user.name

#   policy = <<EOF
# {
#   "Version": "2012-10-17",
#   "Statement": [
#     {
#       "Effect": "Allow",
#       "Action": [
#         "dynamodb:GetItem",
#         "dynamodb:PutItem",
#         "dynamodb:DeleteItem"
#       ],
#       "Resource": "arn:aws:dynamodb:*:*:table/${var.terraform_lock_table_name}"
#     }
#   ]
# }
# EOF
# }

# resource "aws_iam_user_policy_attachment" "base_user_additional_policies" {
#   count      = length(var.base_user_additional_policies_arn)
#   role       = aws_iam_role.aws_iam_user.base_user.id
#   policy_arn = var.base_user_additional_policies_arn[count.index]
# }
