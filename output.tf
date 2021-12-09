output "user_name" {
    value = aws_iam_user.base_user.name
}

output "aws_access_key_id"{
    value = aws_iam_access_key.base_user.id
    description="Access key ID"
}

output "base64_aws_secret_access_key"{
    sensitive=true
    value = aws_iam_access_key.base_user.encrypted_secret
    description= "Encrypted secret, base64 encoded, if base_user_pgp_key was specified. This attribute is not available for imported resources."
}

output "ses_smtp_base64_aws_secret_access_key"{
    sensitive=true
    value = aws_iam_access_key.base_user.encrypted_ses_smtp_password_v4
    description="Encrypted SES SMTP password, base64 encoded, if base_user_pgp_key was specified. This attribute is not available for imported resources."
}

output "ses_smtp_password_aws_secret_access_key"{
    value = aws_iam_access_key.base_user.ses_smtp_password_v4
    description = "Secret access key converted into an SES SMTP password by applying AWS's documented Sigv4 conversion algorithm."
    sensitive=true
}