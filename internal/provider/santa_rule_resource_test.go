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
						resourceName, "policy", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "target_type", "BINARY"),
					resource.TestCheckResourceAttr(
						resourceName, "target_identifier", "9f3e7b21a0a745297dd906dad4a4a4637bdec066cdf331b457230aa32fe68b4b"),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "custom_message", ""),
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
						resourceName, "policy", "2"),
					resource.TestCheckResourceAttr(
						resourceName, "target_type", "CERTIFICATE"),
					resource.TestCheckResourceAttr(
						resourceName, "target_identifier", "bff4a6a4d6b42e94e7d7f48e66b66c69b58fb409a785e0c65409e3bef9ad8887"),
					resource.TestCheckResourceAttr(
						resourceName, "description", "description"),
					resource.TestCheckResourceAttr(
						resourceName, "custom_message", "custom message"),
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
  policy            = 1
  target_type       = "BINARY"
  target_identifier = "9f3e7b21a0a745297dd906dad4a4a4637bdec066cdf331b457230aa32fe68b4b"
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
  policy                  = 2
  target_type             = "CERTIFICATE"
  target_identifier       = "bff4a6a4d6b42e94e7d7f48e66b66c69b58fb409a785e0c65409e3bef9ad8887"
  description             = "description"
  custom_message          = "custom message"
  primary_users           = ["un", "deux"]
  excluded_primary_users  = ["trois", "quatre"]
  serial_numbers          = ["cinq", "six"]
  excluded_serial_numbers = ["sept", "huit"]
  tag_ids                 = [zentral_tag.test.id]
  excluded_tag_ids        = [zentral_tag.test2.id]
}
`, name, tagName, tag2Name)
}
