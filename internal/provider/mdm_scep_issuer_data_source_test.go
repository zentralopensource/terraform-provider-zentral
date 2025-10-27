package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMSCEPIssuerDataSource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	thirdName := acctest.RandString(12)
	resource1Name := "zentral_mdm_scep_issuer.test1"
	resource2Name := "zentral_mdm_scep_issuer.test2"
	resource3Name := "zentral_mdm_scep_issuer.test3"
	ds1ResourceName := "data.zentral_mdm_scep_issuer.check1_by_name"
	ds2ResourceName := "data.zentral_mdm_scep_issuer.check2_by_id"
	ds3ResourceName := "data.zentral_mdm_scep_issuer.check3_by_name"
	ds4ResourceName := "data.zentral_mdm_scep_issuer.check4_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMSCEPIssuerDataSourceConfig(firstName, secondName, thirdName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", resource1Name, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "description", "Description"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "url", "https://www.example.com/scep"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "key_usage", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "key_size", "2048"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "backend", "STATIC_CHALLENGE"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "digicert"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "ident"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "microsoft_ca"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "okta_ca"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "static_challenge.challenge", "Yolo"),
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", resource2Name, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "description", "Description"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "url", "https://www.example.com/scep"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "key_usage", "0"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "key_size", "2048"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "backend", "OKTA_CA"),
					resource.TestCheckNoResourceAttr(
						ds2ResourceName, "digicert"),
					resource.TestCheckNoResourceAttr(
						ds2ResourceName, "ident"),
					resource.TestCheckNoResourceAttr(
						ds2ResourceName, "microsoft_ca"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "okta_ca.url", "https://www.example.com/okta"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "okta_ca.username", "yolo"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "okta_ca.password", "fomo"),
					resource.TestCheckNoResourceAttr(
						ds2ResourceName, "static_challenge"),
					// Provisioned SCEP issuer read by name
					resource.TestCheckResourceAttr(
						ds3ResourceName, "description", ""),
					resource.TestCheckResourceAttr(
						ds3ResourceName, "url", "https://github.com/zentralopensource/terraform-provider-zentral"),
					resource.TestCheckResourceAttr(
						ds3ResourceName, "key_usage", "0"),
					resource.TestCheckResourceAttr(
						ds3ResourceName, "key_size", "2048"),
					resource.TestCheckNoResourceAttr(
						ds3ResourceName, "backend"),
					resource.TestCheckNoResourceAttr(
						ds3ResourceName, "digicert"),
					resource.TestCheckNoResourceAttr(
						ds3ResourceName, "microsoft_ca"),
					resource.TestCheckNoResourceAttr(
						ds3ResourceName, "okta_ca"),
					resource.TestCheckNoResourceAttr(
						ds3ResourceName, "static_challenge"),
					// Read by ID 2
					resource.TestCheckResourceAttrPair(
						ds4ResourceName, "id", resource3Name, "id"),
					resource.TestCheckResourceAttr(
						ds4ResourceName, "name", thirdName),
					resource.TestCheckResourceAttr(
						ds4ResourceName, "description", "Description"),
					resource.TestCheckResourceAttr(
						ds4ResourceName, "url", "https://www.example.com/scep"),
					resource.TestCheckResourceAttr(
						ds4ResourceName, "key_usage", "0"),
					resource.TestCheckResourceAttr(
						ds4ResourceName, "key_size", "2048"),
					resource.TestCheckResourceAttr(
						ds4ResourceName, "backend", "DIGICERT"),
					resource.TestCheckResourceAttr(
						ds4ResourceName, "digicert.api_base_url", "https://one.digicert.com/mpki/api/"),
					resource.TestCheckResourceAttr(
						ds4ResourceName, "digicert.api_token", "secret"),
					resource.TestCheckResourceAttr(
						ds4ResourceName, "digicert.profile_guid", "bea51415-5f66-47d1-8c84-fe774e28deb5"),
					resource.TestCheckResourceAttr(
						ds4ResourceName, "digicert.business_unit_guid", "5483b378-e99f-47b8-bafb-d0452ae025ca"),
					resource.TestCheckResourceAttr(
						ds4ResourceName, "digicert.seat_type", "DEVICE_SEAT"),
					resource.TestCheckResourceAttr(
						ds4ResourceName, "digicert.seat_id_mapping", "common_name"),
					resource.TestCheckResourceAttr(
						ds4ResourceName, "digicert.default_seat_email", "yolo@example.com"),
					resource.TestCheckNoResourceAttr(
						ds4ResourceName, "ident"),
					resource.TestCheckNoResourceAttr(
						ds4ResourceName, "microsoft_ca"),
					resource.TestCheckNoResourceAttr(
						ds4ResourceName, "okta_ca"),
					resource.TestCheckNoResourceAttr(
						ds4ResourceName, "static_challenge"),
				),
			},
		},
	})
}

func testAccMDMSCEPIssuerDataSourceConfig(firstName string, secondName string, thirdName string) string {
	return fmt.Sprintf(`
resource  "zentral_mdm_scep_issuer" "test1" {
  name             = %[1]q
  description      = "Description"
  url              = "https://www.example.com/scep"
  backend          = "STATIC_CHALLENGE"
  static_challenge = {
      challenge = "Yolo"
  }
}

resource  "zentral_mdm_scep_issuer" "test2" {
  name             = %[2]q
  description      = "Description"
  url              = "https://www.example.com/scep"
  backend          = "OKTA_CA"
  okta_ca          = {
    "url"      = "https://www.example.com/okta"
    "username" = "yolo"
    "password" = "fomo"
  }
}

resource  "zentral_mdm_scep_issuer" "test3" {
  name             = %[3]q
  description      = "Description"
  url              = "https://www.example.com/scep"
  backend          = "DIGICERT"
  digicert         = {
    "api_token"          = "secret"
    "profile_guid"       = "bea51415-5f66-47d1-8c84-fe774e28deb5"
    "business_unit_guid" = "5483b378-e99f-47b8-bafb-d0452ae025ca"
    "default_seat_email" = "yolo@example.com"
  }
}

data "zentral_mdm_scep_issuer" "check1_by_name" {
  name = zentral_mdm_scep_issuer.test1.name
}

data "zentral_mdm_scep_issuer" "check2_by_id" {
  id = zentral_mdm_scep_issuer.test2.id
}

// provisioned SCEP issuer
data "zentral_mdm_scep_issuer" "check3_by_name" {
  name = "TF provider GitHub"
}

data "zentral_mdm_scep_issuer" "check4_by_id" {
 id = zentral_mdm_scep_issuer.test3.id
}
`, firstName, secondName, thirdName)
}
