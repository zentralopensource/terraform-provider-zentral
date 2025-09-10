package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMACMEIssuerDataSource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resource1Name := "zentral_mdm_acme_issuer.test1"
	resource2Name := "zentral_mdm_acme_issuer.test2"
	ds1ResourceName := "data.zentral_mdm_acme_issuer.check1_by_name"
	ds2ResourceName := "data.zentral_mdm_acme_issuer.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMACMEIssuerDataSourceConfig(firstName, secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", resource1Name, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "description", "Description"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "directory_url", "https://www.example.com/acme"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "key_type", "ECSECPrimeRandom"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "key_size", "256"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "usage_flags", "1"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "extended_key_usage.#", "1"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "extended_key_usage.0", "1.3.6.1.5.5.7.3.2"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "hardware_bound", "true"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "attest", "true"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "backend", "MICROSOFT_CA"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "ident"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "microsoft_ca.url", "https://www.example.com/ndes"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "microsoft_ca.username", "yolo"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "microsoft_ca.password", "fomo"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "okta_ca"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "static_challenge"),
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", resource2Name, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "description", "Description"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "directory_url", "https://www.example.com/acme"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "key_type", "ECSECPrimeRandom"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "key_size", "256"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "usage_flags", "1"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "extended_key_usage.#", "1"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "extended_key_usage.0", "1.3.6.1.5.5.7.3.2"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "hardware_bound", "true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "attest", "true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "backend", "IDENT"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "ident.url", "https://www.example.com/ident"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "ident.bearer_token", "YoloFomo"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "ident.request_timeout", "30"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "ident.max_retries", "3"),
					resource.TestCheckNoResourceAttr(
						ds2ResourceName, "microsoft_ca"),
					resource.TestCheckNoResourceAttr(
						ds2ResourceName, "okta_ca"),
					resource.TestCheckNoResourceAttr(
						ds2ResourceName, "static_challenge"),
				),
			},
		},
	})
}

func testAccMDMACMEIssuerDataSourceConfig(firstName string, secondName string) string {
	return fmt.Sprintf(`
resource  "zentral_mdm_acme_issuer" "test1" {
  name               = %[1]q
  description        = "Description"
  directory_url      = "https://www.example.com/acme"
  key_type           = "ECSECPrimeRandom"
  key_size           = 256
  usage_flags        = 1
  extended_key_usage = ["1.3.6.1.5.5.7.3.2"]
  hardware_bound     = true
  attest             = true
  backend            = "MICROSOFT_CA"
  microsoft_ca       = {
    "url"      = "https://www.example.com/ndes"
    "username" = "yolo"
    "password" = "fomo"
  }
}

resource  "zentral_mdm_acme_issuer" "test2" {
  name               = %[2]q
  description        = "Description"
  directory_url      = "https://www.example.com/acme"
  key_type           = "ECSECPrimeRandom"
  key_size           = 256
  usage_flags        = 1
  extended_key_usage = ["1.3.6.1.5.5.7.3.2"]
  hardware_bound     = true
  attest             = true
  backend            = "IDENT"
  ident              = {
    "url"          = "https://www.example.com/ident"
    "bearer_token" = "YoloFomo"
  }
}

data "zentral_mdm_acme_issuer" "check1_by_name" {
  name = zentral_mdm_acme_issuer.test1.name
}

data "zentral_mdm_acme_issuer" "check2_by_id" {
  id = zentral_mdm_acme_issuer.test2.id
}
`, firstName, secondName)
}
