package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSantaRuleResource(t *testing.T) {
	name := acctest.RandString(12)
	tagName := acctest.RandString(12)
	tag2Name := acctest.RandString(12)
	resourceName := "zentral_santa_rule.test"
	cfgResourceName := "zentral_santa_configuration.test"
	tagResourceName := "zentral_tag.test"
	tag2ResourceName := "zentral_tag.test2"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccSantaRuleResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "configuration_id", cfgResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "policy", "ALLOWLIST"),
					resource.TestCheckResourceAttr(
						resourceName, "cel_expr", ""),
					resource.TestCheckResourceAttr(
						resourceName, "target_type", "CDHASH"),
					resource.TestCheckResourceAttr(
						resourceName, "target_identifier", "9f3e7b21a0a745297dd906dad4a4a4637bdec066"),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "custom_message", ""),
					resource.TestCheckResourceAttr(
						resourceName, "custom_url", ""),
					resource.TestCheckNoResourceAttr(
						resourceName, "ruleset_id"),
					resource.TestCheckResourceAttr(
						resourceName, "primary_users.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_primary_users.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "serial_numbers.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_serial_numbers.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "1"),
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
				Config: testAccSantaRuleResourceConfigFull(name, tagName, tag2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "configuration_id", cfgResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "policy", "CEL"),
					resource.TestCheckResourceAttr(
						resourceName, "cel_expr", "target.signing_time >= timestamp('2025-05-31T00:00:00Z')"),
					resource.TestCheckResourceAttr(
						resourceName, "target_type", "SIGNINGID"),
					resource.TestCheckResourceAttr(
						resourceName, "target_identifier", "platform:com.apple.curl"),
					resource.TestCheckResourceAttr(
						resourceName, "description", "description"),
					resource.TestCheckResourceAttr(
						resourceName, "custom_message", "custom message"),
					resource.TestCheckResourceAttr(
						resourceName, "custom_url", "https://zentral.com"),
					resource.TestCheckNoResourceAttr(
						resourceName, "ruleset_id"),
					resource.TestCheckResourceAttr(
						resourceName, "primary_users.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "primary_users.*", "un"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "primary_users.*", "deux"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_primary_users.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "excluded_primary_users.*", "trois"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "excluded_primary_users.*", "quatre"),
					resource.TestCheckResourceAttr(
						resourceName, "serial_numbers.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "serial_numbers.*", "cinq"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "serial_numbers.*", "six"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_serial_numbers.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "excluded_serial_numbers.*", "sept"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "excluded_serial_numbers.*", "huit"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "tag_ids.*", tagResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "excluded_tag_ids.*", tag2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "2"),
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

func testAccSantaRuleResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_santa_configuration" "test" {
  name = %[1]q
}

resource "zentral_santa_rule" "test" {
  configuration_id  = zentral_santa_configuration.test.id
  policy            = "ALLOWLIST"
  target_type       = "CDHASH"
  target_identifier = "9f3e7b21a0a745297dd906dad4a4a4637bdec066"
}
`, name)
}

func testAccSantaRuleResourceConfigFull(name string, tagName string, tag2Name string) string {
	return fmt.Sprintf(`
resource "zentral_santa_configuration" "test" {
  name = %[1]q
}

resource "zentral_taxonomy" "test" {
  name = %[1]q
}

resource "zentral_tag" "test" {
  taxonomy_id = zentral_taxonomy.test.id
  name        = %[2]q
}

resource "zentral_tag" "test2" {
  taxonomy_id = zentral_taxonomy.test.id
  name        = %[3]q
}

resource "zentral_santa_rule" "test" {
  configuration_id        = zentral_santa_configuration.test.id
  policy                  = "CEL"
  cel_expr                = "target.signing_time >= timestamp('2025-05-31T00:00:00Z')"
  target_type             = "SIGNINGID"
  target_identifier       = "platform:com.apple.curl"
  description             = "description"
  custom_message          = "custom message"
  custom_url			  = "https://zentral.com"
  primary_users           = ["un", "deux"]
  excluded_primary_users  = ["trois", "quatre"]
  serial_numbers          = ["cinq", "six"]
  excluded_serial_numbers = ["sept", "huit"]
  tag_ids                 = [zentral_tag.test.id]
  excluded_tag_ids        = [zentral_tag.test2.id]
}
`, name, tagName, tag2Name)
}
