locals {
  create_role_count       = var.create_base_role ? 1 : 0
  allow_user_assume_count = var.allow_user_assume_on_role ? 1 : 0
}

### AWS IAM Role section
resource "aws_iam_role" "base_role" {
  name               = var.role_name
  count              = local.create_role_count
  assume_role_policy = var.allow_user_assume_on_role ? data.aws_iam_policy_document.user_trust_relationship[0].json : data.aws_iam_policy_document.ec2_trust_relationship.json
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
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:user/${var.user_name}"]
    }
  }
}
### AWS IAM Policies attachment section
resource "aws_iam_role_policy_attachment" "role_s3_access" {
  count      = local.create_role_count
  role       = aws_iam_role.base_role[count.index].name
  policy_arn = aws_iam_policy.s3_access.arn
}

resource "aws_iam_role_policy_attachment" "role_dynamodb_access" {
  count      = local.create_role_count
  role       = aws_iam_role.base_role[count.index].name
  policy_arn = aws_iam_policy.dynamodb_access.arn
}

resource "aws_iam_role_policy_attachment" "role_additional_policies" {
  count      = length(var.additional_policies_arn)
  role       = aws_iam_role.base_role[count.index].name
  policy_arn = var.additional_policies_arn[count.index]
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
      "Resource": "arn:aws:s3:::${aws_s3_bucket.state_storage.id}"
    },
    {
      "Effect": "Allow",
      "Action": ["s3:GetObject", "s3:PutObject", "s3:DeleteObject"],
      "Resource": "arn:aws:s3:::${aws_s3_bucket.state_storage.id}${var.s3_state_key_path}"
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
      "Resource": "${aws_dynamodb_table.terraform_lock.arn}"
    }
  ]
}
EOF
}
