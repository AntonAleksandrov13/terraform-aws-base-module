terraform {
  required_providers {
    # Because we're currently using a built-in provider as
    # a substitute for dedicated Terraform language syntax
    # for now, test suite modules must always declare a
    # dependency on this provider. This provider is only
    # available when running tests, so you shouldn't use it
    # in non-test modules.
    test = {
      source = "terraform.io/builtin/test"
    }

    # This example also uses the "http" data source to
    # verify the behavior of the hypothetical running
    # service, so we should declare that too.
    http = {
      source = "hashicorp/http"
    }
  }
}

module "main" {
  # source is always ../.. for test suite configurations,
  # because they are placed two subdirectories deep under
  # the main module directory.
  source = "../.."

  # This test suite is aiming to test the "defaults" for
  # this module, so it doesn't set any input variables
  # and just lets their default values be selected instead.
}


# The special test_assertions resource type, which belongs
# to the test provider we required above, is a temporary
# syntax for writing out explicit test assertions.
resource "test_assertions" "api_url" {
  # "component" serves as a unique identifier for this
  # particular set of assertions in the test results.
  component = "api_url"

  # equal and check blocks serve as the test assertions.
  # the labels on these blocks are unique identifiers for
  # the assertions, to allow more easily tracking changes
  # in success between runs.

  equal "scheme" {
    description = "default scheme is https"
    got         = "https"
    want        = "https"
  }
}
