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
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"
)

type Item struct {
	LockID    string
	LockValue string
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

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

func getRoleAttachedPolicies(session *session.Session, roleName string) ([]string, error) {
	svc := iam.New(session)
	policyNameList := []string{}
	err := svc.ListAttachedRolePoliciesPages(
		&iam.ListAttachedRolePoliciesInput{
			RoleName: &roleName,
		},
		func(page *iam.ListAttachedRolePoliciesOutput, lastPage bool) bool {
			if page != nil && len(page.AttachedPolicies) > 0 {
				for _, policy := range page.AttachedPolicies {
					policyNameList = append(policyNameList, *policy.PolicyName)
				}
				return true
			}
			return false
		},
	)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return []string{}, aerr
		}
	}
	return policyNameList, nil
}

func currentUserAssumeRole(session client.ConfigProvider, role string) (*sts.AssumeRoleOutput, error) {
	svc := sts.New(session)
	sessionName := getRandomString(6) + "_session"
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

func getRandomString(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func TestRoleCreation(t *testing.T) {
	sess := session.Must(session.NewSession())
	callerAccount, err := getAWSAccountNumber(sess)
	require.NoError(t, err)

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

func TestAdditionalPolicyAttachment(t *testing.T) {
	sess := session.Must(session.NewSession())
	callerAccount, err := getAWSAccountNumber(sess)
	require.NoError(t, err)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/role_with_additional_policies",
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
	additionalPolicyName := terraform.Output(t, terraformOptions, "test_policy_name")
	policyNameList := retry.DoWithRetryInterface(t, "retry", 2, 10*time.Second, func() (interface{}, error) {
		policyNameList, err := getRoleAttachedPolicies(sess, roleNameReturned)
		if err != nil {
			return []string{}, fmt.Errorf("could not fetch attached list of attached policies or policy is not attached")
		}
		return policyNameList, nil
	})
	found := false
	switch reflect.TypeOf(policyNameList).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(policyNameList)
		for i := 0; i < s.Len(); i++ {
			if s.Index(i).String() == additionalPolicyName {
				found = true
				break
			}
		}
	}
	require.Equal(t, true, found)
}

func TestExistingUserCanAssumeRole(t *testing.T) {
	sess := session.Must(session.NewSession())
	callerAccount, err := getAWSAccountNumber(sess)
	require.NoError(t, err)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/existing_user_with_assume_role",
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
	require.NoError(t, err)

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
		TerraformDir: "../examples/existing_user_with_assume_role",
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
	// can the assumed role write to S3?
	_, err := retry.DoWithRetryE(t, "retry", 2, 10*time.Second, func() (string, error) {
		err := uploadFileToS3Bucket(sess, "test.txt", terraform.Output(t, terraformOptions, "s3_bucket_name"))
		if err != nil {
			return "", fmt.Errorf("could not upload a file to s3 bucket")
		}
		return "", nil
	})
	require.NoError(t, err)

	// can the assumed role delete from S3?
	_, err = retry.DoWithRetryE(t, "retry", 2, 10*time.Second, func() (string, error) {
		err = deleteFileFromS3Bucket(sess, "test.txt", terraform.Output(t, terraformOptions, "s3_bucket_name"))
		if err != nil {
			return "", fmt.Errorf("could notdelete a file from s3 bucket")
		}
		return "", nil
	})
	require.NoError(t, err)

}

func TestExistingUserCRUDDynamoDBLockTable(t *testing.T) {
	sess := session.Must(session.NewSession())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/existing_user_with_assume_role",
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	roleARNReturned := terraform.Output(t, terraformOptions, "role_arn")
	tableNameReturned := terraform.Output(t, terraformOptions, "lock_table_name")
	// dirty trick to bypass dynamodb reachability issue
	time.Sleep(10 * time.Second)
	// run role assume and create a new session
	sess = session.Must(session.NewSession(&aws.Config{
		Credentials: stscreds.NewCredentials(sess, roleARNReturned),
	}))

	id := getRandomString(6)
	value := getRandomString(12)
	// can the assumed role write to the dynamodb table?
	_, err := retry.DoWithRetryE(t, "retry", 2, 10*time.Second, func() (string, error) {
		err := addLockTableItem(sess, id, value, tableNameReturned)
		if err != nil {
			return "", fmt.Errorf("could not add item to Dynamodb")
		}
		return "", nil
	})
	require.NoError(t, err)

	// can the assumed role delete data from the dynamodb table?
	_, err = retry.DoWithRetryE(t, "retry", 2, 10*time.Second, func() (string, error) {
		err = deleteLockTableItem(sess, id, tableNameReturned)
		if err != nil {
			return "", fmt.Errorf("could not add item to Dynamodb")
		}
		return "", nil
	})
	require.NoError(t, err)
}
