package test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

type Item struct {
	LockID    string
	LockValue string
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

func uploadFileToS3Bucket(session *session.Session, filename string, bucketName string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Could not close file properly in defer statement")
		}
	}(file)
	uploader := s3manager.NewUploader(session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucketName)
	return nil
}

func deleteFileFromS3Bucket(sess *session.Session, filename string, bucketName string) error {
	svc := s3.New(sess)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &filename,
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: &bucketName,
		Key:    &filename,
	})
	if err != nil {
		return err
	}

	return nil
}

func addLockTableItem(sess *session.Session, lockID string, lockValue string, tableName string) error {
	svc := dynamodb.New(sess)
	item := Item{
		LockID:    lockID,
		LockValue: lockValue,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: &tableName,
	})
	if err != nil {
		return err
	}

	return nil
}

func deleteLockTableItem(sess *session.Session, lockID string, tableName string) error {
	svc := dynamodb.New(sess)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"LockID": {
				S: &lockID,
			},
		},
		TableName: &tableName,
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
}

func TestOnlyRoleCreation(t *testing.T) {
	sess := session.Must(session.NewSession())
	callerAccount, err := getAWSAccountNumber(sess)
	if err != nil {
		require.NoError(t, err)
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
	require.Equal(t, "terraform", roleNameReturned)

	roleARNReturned := terraform.Output(t, terraformOptions, "role_arn")
	require.Equal(t, fmt.Sprintf("arn:aws:iam::%v:role/terraform", callerAccount), roleARNReturned)
}

func TestExistingUserCanAssumeRole(t *testing.T) {
	sess := session.Must(session.NewSession())
	callerAccount, err := getAWSAccountNumber(sess)
	if err != nil {
		require.NoError(t, err)
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/existing_user_can_assume_role",
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// is the role name output matching the expected name?
	roleNameReturned := terraform.Output(t, terraformOptions, "role_name")
	require.Equal(t, "terraform", roleNameReturned)
	// is the role arn output matching the expected arn?
	roleARNReturned := terraform.Output(t, terraformOptions, "role_arn")
	require.Equal(t, fmt.Sprintf("arn:aws:iam::%v:role/terraform", callerAccount), roleARNReturned)

	// run role assume using current user credentials
	assumeRole, err := currentUserAssumeRole(sess, roleARNReturned)
	if err != nil {
		require.NoError(t, err)
	}
	userCanAssumeRole := false
	// check if this property is present.
	// it confirms that the role has been assumed correctly
	if assumeRole.AssumedRoleUser != nil {
		userCanAssumeRole = true
	}
	// can current user assume this role?
	require.Equal(t, true, userCanAssumeRole)
}

func TestExistingUserReadWriteS3Bucket(t *testing.T) {
	sess := session.Must(session.NewSession())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/existing_user_can_assume_role",
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	roleARNReturned := terraform.Output(t, terraformOptions, "role_arn")

	// run role assume and create a new session
	sess = session.Must(session.NewSession(&aws.Config{
		Credentials: stscreds.NewCredentials(sess, roleARNReturned),
	}))
	// dirty trick to bypass s3 reachability issue
	time.Sleep(5 * time.Second)
	err := uploadFileToS3Bucket(sess, "test.txt", terraform.Output(t, terraformOptions, "s3_bucket_name"))
	// can the assumed role write to S3?
	require.NoError(t, err)

	// can the assumed role delete from S3?
	err = deleteFileFromS3Bucket(sess, "test.txt", terraform.Output(t, terraformOptions, "s3_bucket_name"))
	require.NoError(t, err)

}

func TestExistingUserCRUDDynamoDBLockTable(t *testing.T) {
	sess := session.Must(session.NewSession())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/existing_user_can_assume_role",
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	roleARNReturned := terraform.Output(t, terraformOptions, "role_arn")
	tableNameReturned := terraform.Output(t, terraformOptions, "lock_table_name")
	// dirty trick to bypass dynamodb reachability issue
	time.Sleep(5 * time.Second)
	// run role assume and create a new session
	sess = session.Must(session.NewSession(&aws.Config{
		Credentials: stscreds.NewCredentials(sess, roleARNReturned),
	}))

	id := "some_lock_id"
	value := "some_lock_value"
	// can the assumed role write to the dynamodb table?
	err := addLockTableItem(sess, id, value, tableNameReturned)
	require.NoError(t, err)

	// can the assumed role delete data from the dynamodb table?
	err = deleteLockTableItem(sess, id, tableNameReturned)
	require.NoError(t, err)
}
