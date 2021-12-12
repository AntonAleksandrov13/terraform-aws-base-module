package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformHelloWorldExample(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: "../examples/only_user",
		Vars: map[string]interface{}{
			"base_user_pgp_key": "mDMEYbOgExYJKwYBBAHaRw8BAQdAzEXnH4epreJ6d7nrnzszhHkQPy8NlMTJ/ZOkjxgHZ+m0IkFudG9uIDxhbnRvbmFsZWtzYW5kcm92QGdtYWlsLmNvbT6ImgQTFgoAQhYhBMdt1/dLxOPUx+wPaCftDhiwJxIGBQJhs6ATAhsDBQkDwmcABQsJCAcCAyICAQYVCgkICwIEFgIDAQIeBwIXgAAKCRAn7Q4YsCcSBt7lAP9bwH2KoMAHs87mWh8OLbN3RzXF5HJxVCzqE8jSRvTEHAEA6AWhU7T/KBBaJR5uricnI017gjka/30Q5eyf++ScbwW4OARhs6ATEgorBgEEAZdVAQUBAQdAYY8/kacV7C4gxYuwhJVLpKi+miXFF25d2k+NeQr9sXsDAQgHiH4EGBYKACYWIQTHbdf3S8Tj1MfsD2gn7Q4YsCcSBgUCYbOgEwIbDAUJA8JnAAAKCRAn7Q4YsCcSBvzUAQD54FVoMPzQZC3cW4LVU8d5qpLHKGsD31jzUwGvUTMooAEAkBEJaYoiVwTj4htlKO5LDWbeRnf2lTW0eW+EnMh6bAI=",
		},
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.Output(t, terraformOptions, "user_name")
	assert.Equal(t, "terraform", output)
}
