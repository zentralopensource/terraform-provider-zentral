package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMSCEPIssuerResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	thirdName := acctest.RandString(12)
	fourthName := acctest.RandString(12)
	resourceName := "zentral_mdm_scep_issuer.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMSCEPIssuerResourceFirstConfig(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "url", "https://www.example.com/scep1"),
					resource.TestCheckResourceAttr(
						resourceName, "key_usage", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "key_size", "2048"),
					resource.TestCheckResourceAttr(
						resourceName, "backend", "STATIC_CHALLENGE"),
					resource.TestCheckNoResourceAttr(
						resourceName, "digicert"),
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
				Config: testAccMDMSCEPIssuerResourceSecondConfig(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "Description"),
					resource.TestCheckResourceAttr(
						resourceName, "url", "https://www.example.com/scep2"),
					resource.TestCheckResourceAttr(
						resourceName, "key_usage", "5"),
					resource.TestCheckResourceAttr(
						resourceName, "key_size", "4096"),
					resource.TestCheckResourceAttr(
						resourceName, "backend", "MICROSOFT_CA"),
					resource.TestCheckNoResourceAttr(
						resourceName, "digicert"),
					resource.TestCheckNoResourceAttr(
						resourceName, "ident"),
					resource.TestCheckResourceAttr(
						resourceName, "microsoft_ca.url", "https://www.example.com/ndes"),
					resource.TestCheckResourceAttr(
						resourceName, "microsoft_ca.username", "Yolo"),
					resource.TestCheckResourceAttr(
						resourceName, "microsoft_ca.password", "Fomo"),
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
			// Update and Read
			{
				Config: testAccMDMSCEPIssuerResourceThirdConfig(thirdName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", thirdName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "Description"),
					resource.TestCheckResourceAttr(
						resourceName, "url", "https://www.example.com/scep2"),
					resource.TestCheckResourceAttr(
						resourceName, "key_usage", "4"),
					resource.TestCheckResourceAttr(
						resourceName, "key_size", "2048"),
					resource.TestCheckResourceAttr(
						resourceName, "backend", "IDENT"),
					resource.TestCheckNoResourceAttr(
						resourceName, "digicert"),
					resource.TestCheckResourceAttr(
						resourceName, "ident.url", "https://www.example.com/ident"),
					resource.TestCheckResourceAttr(
						resourceName, "ident.bearer_token", "YoloFomo"),
					resource.TestCheckResourceAttr(
						resourceName, "ident.request_timeout", "30"),
					resource.TestCheckResourceAttr(
						resourceName, "ident.max_retries", "3"),
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
			// Update and Read
			{
				Config: testAccMDMSCEPIssuerResourceFourthConfig(fourthName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", fourthName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "Description"),
					resource.TestCheckResourceAttr(
						resourceName, "url", "https://www.example.com/scep3"),
					resource.TestCheckResourceAttr(
						resourceName, "key_usage", "4"),
					resource.TestCheckResourceAttr(
						resourceName, "key_size", "2048"),
					resource.TestCheckResourceAttr(
						resourceName, "backend", "DIGICERT"),
					resource.TestCheckResourceAttr(
						resourceName, "digicert.api_base_url", "https://www.example.com/api/"),
					resource.TestCheckResourceAttr(
						resourceName, "digicert.api_token", "YoloFomo"),
					resource.TestCheckResourceAttr(
						resourceName, "digicert.profile_guid", "531049b7-8fff-4d4f-8250-0e3e048ac957"),
					resource.TestCheckResourceAttr(
						resourceName, "digicert.business_unit_guid", "88d2b54e-cb5b-4197-8ed9-8a77d6f15806"),
					resource.TestCheckResourceAttr(
						resourceName, "digicert.seat_type", "USER_SEAT"),
					resource.TestCheckResourceAttr(
						resourceName, "digicert.seat_id_mapping", "common_name"),
					resource.TestCheckResourceAttr(
						resourceName, "digicert.default_seat_email", "yolo@example.com"),
					resource.TestCheckNoResourceAttr(
						resourceName, "ident"),
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

func testAccMDMSCEPIssuerResourceFirstConfig(name string) string {
	return fmt.Sprintf(`
resource  "zentral_mdm_scep_issuer" "test" {
  name             = %[1]q
  url              = "https://www.example.com/scep1"
  backend          = "STATIC_CHALLENGE"
  static_challenge = {
      challenge = "Yolo"
  }
}
`, name)
}

func testAccMDMSCEPIssuerResourceSecondConfig(name string) string {
	return fmt.Sprintf(`
resource  "zentral_mdm_scep_issuer" "test" {
  name             = %[1]q
  description      = "Description"
  url              = "https://www.example.com/scep2"
  key_size         = 4096
  key_usage        = 5
  backend          = "MICROSOFT_CA"
  microsoft_ca = {
      url      = "https://www.example.com/ndes"
      username = "Yolo"
      password = "Fomo"
  }
}
`, name)
}

func testAccMDMSCEPIssuerResourceThirdConfig(name string) string {
	return fmt.Sprintf(`
resource  "zentral_mdm_scep_issuer" "test" {
  name             = %[1]q
  description      = "Description"
  url              = "https://www.example.com/scep2"
  key_size         = 2048
  key_usage        = 4
  backend          = "IDENT"
  ident = {
      url          = "https://www.example.com/ident"
      bearer_token = "YoloFomo"
  }
}
`, name)
}

func testAccMDMSCEPIssuerResourceFourthConfig(name string) string {
	return fmt.Sprintf(`
resource  "zentral_mdm_scep_issuer" "test" {
  name             = %[1]q
  description      = "Description"
  url              = "https://www.example.com/scep3"
  key_size         = 2048
  key_usage        = 4
  backend          = "DIGICERT"
  digicert = {
      api_base_url       = "https://www.example.com/api/"
      api_token          = "YoloFomo"
      profile_guid       = "531049b7-8fff-4d4f-8250-0e3e048ac957"
      business_unit_guid = "88d2b54e-cb5b-4197-8ed9-8a77d6f15806"
      seat_type          = "USER_SEAT"
      seat_id_mapping    = "common_name"
      default_seat_email = "yolo@example.com"
  }
}
`, name)
}
