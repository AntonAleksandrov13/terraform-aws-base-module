package test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnlyRoleCreation(t *testing.T) {
	svc := sts.New(session.Must(session.NewSession()))
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				t.Errorf(aerr.Error())
			}
		} else {
			t.Errorf(err.Error())
		}
	}
	callerAccount := *result.Account
	//Construct the terraform options with default retryable errors to handle the most common
	//retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: "../examples/only_role",
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform roleNameReturned` to get the values of roleNameReturned variables and check they have the expected values.
	roleNameReturned := terraform.Output(t, terraformOptions, "role_name")
	assert.Equal(t, "terraform", roleNameReturned)

	roleARNReturned := terraform.Output(t, terraformOptions, "role_arn")
	assert.Equal(t, fmt.Sprintf("arn:aws:iam::%v:role/terraform", callerAccount), roleARNReturned)
}

func TestExistingUserCanAssumeRole(t *testing.T) {
	svc := sts.New(session.Must(session.NewSession()))
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				t.Errorf(aerr.Error())
			}
		} else {
			t.Errorf(err.Error())
		}
	}
	callerAccount := *result.Account
	//Construct the terraform options with default retryable errors to handle the most common
	//retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: "../examples/existing_user_can_assume_role",
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform roleNameReturned` to get the values of roleNameReturned variables and check they have the expected values.
	roleNameReturned := terraform.Output(t, terraformOptions, "role_name")
	assert.Equal(t, "terraform", roleNameReturned)

	roleARNReturned := terraform.Output(t, terraformOptions, "role_arn")
	assert.Equal(t, fmt.Sprintf("arn:aws:iam::%v:role/terraform", callerAccount), roleARNReturned)
}
