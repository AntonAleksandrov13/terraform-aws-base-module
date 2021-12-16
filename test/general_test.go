package test

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestOnlyRoleCreation(t *testing.T) {
	sess := session.Must(session.NewSession())
	callerAccount, err := getAWSAccountNumber(sess)
	if err != nil {
		t.Error(err)
	}
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
	sess := session.Must(session.NewSession())
	callerAccount, err := getAWSAccountNumber(sess)
	if err != nil {
		t.Error(err)
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/existing_user_can_assume_role",
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// is role name output matching the expected name?
	roleNameReturned := terraform.Output(t, terraformOptions, "role_name")
	assert.Equal(t, "terraform", roleNameReturned)
	// is the role arn output matching the expected arn?
	roleARNReturned := terraform.Output(t, terraformOptions, "role_arn")
	assert.Equal(t, fmt.Sprintf("arn:aws:iam::%v:role/terraform", callerAccount), roleARNReturned)

	// run role assume using current user credentials
	assumeRole, err := currentUserAssumeRole(sess, roleARNReturned)
	if err != nil {
		t.Error(err)
	}
	userCanAssumeRole := false
	// check if this property is present.
	// it confirms that the role has been assumed correctly
	if assumeRole.AssumedRoleUser != nil {
		userCanAssumeRole = true
	}
	// can current user assume this role?
	assert.Equal(t, true, userCanAssumeRole)

	sess = session.Must(session.NewSession(&aws.Config{
		Credentials: stscreds.NewCredentials(sess, roleARNReturned),
	}))
	uploadFile(sess, "test.txt", terraform.Output(t, terraformOptions, "s3_bucket_name"))
	// can we write to S3?

}

func getAWSAccountNumber(session client.ConfigProvider) (string, error) {
	svc := sts.New(session)
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return "", aerr
		}
	}
	return *result.Account, nil
}

func currentUserAssumeRole(session client.ConfigProvider, role string) (*sts.AssumeRoleOutput, error) {
	svc := sts.New(session)
	sessionName := "test_session"
	result, err := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         &role,
		RoleSessionName: &sessionName,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func uploadFile(session *session.Session, uploadFile string, bucketName string) error {
	upFile, err := os.Open(uploadFile)
	if err != nil {
		return err
	}
	defer upFile.Close()

	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(bucketName),
		Key:                aws.String("/"),
		ACL:                aws.String("private"),
		Body:               bytes.NewReader(fileBuffer),
		ContentLength:      aws.Int64(fileSize),
		ContentType:        aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition: aws.String("attachment"),
	})
	fmt.Println(err)
	return err
}
