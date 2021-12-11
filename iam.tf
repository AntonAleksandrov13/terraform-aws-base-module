locals {
  create_user_count       = var.create_base_user ? 1 : 0
  create_role_count       = var.create_base_role ? 1 : 0
  allow_user_assume_count = var.allow_user_assume ? 1 : 0
}
### AWS IAM User section
resource "aws_iam_user" "base_user" {
  count = local.create_user_count
  name  = var.base_user_name
  path  = var.base_user_path
}

resource "aws_iam_access_key" "base_user" {
  count   = local.create_user_count
  user    = aws_iam_user.base_user[count.index].name
  pgp_key = var.base_user_pgp_key
}

### AWS IAM Role section
resource "aws_iam_role" "base_role" {
  name               = var.base_role_name
  count              = local.create_role_count
  assume_role_policy = var.allow_user_assume ? data.aws_iam_policy_document.user_trust_relationship[0].json : data.aws_iam_policy_document.ec2_trust_relationship.json
}

data "aws_iam_policy_document" "ec2_trust_relationship" {
  statement {
    sid = ""

    actions = [
      "sts:AssumeRole",
    ]
    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "user_trust_relationship" {
  count = local.allow_user_assume_count
  statement {
    sid = ""

    actions = [
      "sts:AssumeRole",
    ]
    principals {
      type        = "AWS"
      # NOTE: that here the module does use a dependency on aws_iam_user.base_user[0].name resource
      # this is intentional. sincer sometimes you might want to allow already existing user to be allowed to use this role
      # so, you can still use -var `base_user_name=already_existing_user`.
      # this functionality comes with the price: if the user does not exists "Invalid principal in policy" error will be returned
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:user/${var.base_user_name}"]
    }
  }
}

### AWS IAM Policies section
resource "aws_iam_policy" "s3_access" {
  name = "S3AccessTerraform"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "s3:ListBucket",
      "Resource": "arn:aws:s3:::${var.state_bucket_name}"
    },
    {
      "Effect": "Allow",
      "Action": ["s3:GetObject", "s3:PutObject", "s3:DeleteObject"],
      "Resource": "arn:aws:s3:::${var.state_bucket_name}${var.s3_state_key_path}"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "dynamodb_access" {
  name   = "DynamoDBAccessTerraform"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "dynamodb:GetItem",
        "dynamodb:PutItem",
        "dynamodb:DeleteItem"
      ],
      "Resource": "arn:aws:dynamodb:*:*:table/${var.terraform_lock_table_name}"
    }
  ]
}
EOF
}
