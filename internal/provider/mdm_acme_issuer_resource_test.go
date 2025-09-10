package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMACMEIssuerResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	thirdName := acctest.RandString(12)
	resourceName := "zentral_mdm_acme_issuer.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMACMEIssuerResourceFirstConfig(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "directory_url", "https://www.example.com/acme1"),
					resource.TestCheckResourceAttr(
						resourceName, "key_type", "ECSECPrimeRandom"),
					resource.TestCheckResourceAttr(
						resourceName, "key_size", "384"),
					resource.TestCheckResourceAttr(
						resourceName, "usage_flags", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "extended_key_usage.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "extended_key_usage.0", "1.3.6.1.5.5.7.3.2"),
					resource.TestCheckResourceAttr(
						resourceName, "hardware_bound", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "attest", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "backend", "STATIC_CHALLENGE"),
					resource.TestCheckNoResourceAttr(
						resourceName, "ident"),
					resource.TestCheckNoResourceAttr(
						resourceName, "microsoft_ca"),
					resource.TestCheckNoResourceAttr(
						resourceName, "okta_ca"),
					resource.TestCheckResourceAttr(
						resourceName, "static_challenge.challenge", "Yolo"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccMDMACMEIssuerResourceSecondConfig(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "Description"),
					resource.TestCheckResourceAttr(
						resourceName, "directory_url", "https://www.example.com/acme2"),
					resource.TestCheckResourceAttr(
						resourceName, "key_type", "RSA"),
					resource.TestCheckResourceAttr(
						resourceName, "key_size", "2048"),
					resource.TestCheckResourceAttr(
						resourceName, "usage_flags", "5"),
					resource.TestCheckResourceAttr(
						resourceName, "extended_key_usage.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "hardware_bound", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "attest", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "backend", "OKTA_CA"),
					resource.TestCheckNoResourceAttr(
						resourceName, "ident"),
					resource.TestCheckNoResourceAttr(
						resourceName, "microsoft_ca"),
					resource.TestCheckResourceAttr(
						resourceName, "okta_ca.url", "https://www.example.com/okta"),
					resource.TestCheckResourceAttr(
						resourceName, "okta_ca.username", "Yolo"),
					resource.TestCheckResourceAttr(
						resourceName, "okta_ca.password", "Fomo"),
					resource.TestCheckNoResourceAttr(
						resourceName, "static_challenge"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccMDMACMEIssuerResourceThirdConfig(thirdName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", thirdName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "Description"),
					resource.TestCheckResourceAttr(
						resourceName, "directory_url", "https://www.example.com/acme2"),
					resource.TestCheckResourceAttr(
						resourceName, "key_type", "RSA"),
					resource.TestCheckResourceAttr(
						resourceName, "key_size", "4096"),
					resource.TestCheckResourceAttr(
						resourceName, "usage_flags", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "extended_key_usage.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "hardware_bound", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "attest", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "backend", "IDENT"),
					resource.TestCheckResourceAttr(
						resourceName, "ident.url", "https://www.example.com/ident"),
					resource.TestCheckResourceAttr(
						resourceName, "ident.bearer_token", "YoloFomo"),
					resource.TestCheckResourceAttr(
						resourceName, "ident.request_timeout", "123"),
					resource.TestCheckResourceAttr(
						resourceName, "ident.max_retries", "5"),
					resource.TestCheckNoResourceAttr(
						resourceName, "microsoft_ca"),
					resource.TestCheckNoResourceAttr(
						resourceName, "okta_ca"),
					resource.TestCheckNoResourceAttr(
						resourceName, "static_challenge"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMDMACMEIssuerResourceFirstConfig(name string) string {
	return fmt.Sprintf(`
resource  "zentral_mdm_acme_issuer" "test" {
  name               = %[1]q
  directory_url      = "https://www.example.com/acme1"
  key_type           = "ECSECPrimeRandom"
  key_size           = 384
  usage_flags        = 1
  extended_key_usage = ["1.3.6.1.5.5.7.3.2"]
  backend          = "STATIC_CHALLENGE"
  static_challenge = {
      challenge = "Yolo"
  }
}
`, name)
}

func testAccMDMACMEIssuerResourceSecondConfig(name string) string {
	return fmt.Sprintf(`
resource  "zentral_mdm_acme_issuer" "test" {
  name               = %[1]q
  description        = "Description"
  directory_url      = "https://www.example.com/acme2"
  key_type           = "RSA"
  key_size           = 2048
  usage_flags        = 5
  hardware_bound     = false
  attest             = false
  backend            = "OKTA_CA"
  okta_ca            = {
      url      = "https://www.example.com/okta"
      username = "Yolo"
      password = "Fomo"
  }
}
`, name)
}

func testAccMDMACMEIssuerResourceThirdConfig(name string) string {
	return fmt.Sprintf(`
resource  "zentral_mdm_acme_issuer" "test" {
  name               = %[1]q
  description        = "Description"
  directory_url      = "https://www.example.com/acme2"
  key_type           = "RSA"
  key_size           = 4096
  usage_flags        = 1
  hardware_bound     = false
  attest             = false
  backend            = "IDENT"
  ident              = {
      url             = "https://www.example.com/ident"
      bearer_token    = "YoloFomo"
      request_timeout = 123
      max_retries     = 5
  }
}
`, name)
}
