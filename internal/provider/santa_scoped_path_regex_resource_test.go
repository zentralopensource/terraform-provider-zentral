package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSantaScopedPathRegexResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	tagName := acctest.RandString(12)
	tag2Name := acctest.RandString(12)
	resourceName := "zentral_santa_scoped_path_regex.test"
	cfgResourceName := "zentral_santa_configuration.test"
	tagResourceName := "zentral_tag.test"
	tag2ResourceName := "zentral_tag.test2"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccSantaScopedPathRegexResourceConfigBare(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "configuration_id", cfgResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "policy", "ALLOW"),
					resource.TestCheckResourceAttr(
						resourceName, "regex", "/Applications/.+"),
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
				Config: testAccSantaScopedPathRegexResourceConfigFull(firstName, secondName, tagName, tag2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "configuration_id", cfgResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "description"),
					resource.TestCheckResourceAttr(
						resourceName, "policy", "BLOCK"),
					resource.TestCheckResourceAttr(
						resourceName, "regex", "/Users/[^/]+/Downloads/.+"),
					resource.TestCheckResourceAttr(
						resourceName, "primary_users.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "primary_users.*", "un"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "primary_users.*", "deux"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_primary_users.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "excluded_primary_users.*", "trois"),
					resource.TestCheckResourceAttr(
						resourceName, "serial_numbers.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "serial_numbers.*", "quatre"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_serial_numbers.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "excluded_serial_numbers.*", "cinq"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "excluded_serial_numbers.*", "six"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "tag_ids.*", tagResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "excluded_tag_ids.*", tag2ResourceName, "id"),
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

func testAccSantaScopedPathRegexResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_santa_configuration" "test" {
  name = %[1]q
}

resource "zentral_santa_scoped_path_regex" "test" {
  configuration_id = zentral_santa_configuration.test.id
  name             = %[1]q
  policy           = "ALLOW"
  regex            = "/Applications/.+"
}
`, name)
}

func testAccSantaScopedPathRegexResourceConfigFull(cfgName string, name string, tagName string, tag2Name string) string {
	return fmt.Sprintf(`
resource "zentral_santa_configuration" "test" {
  name = %[1]q
}

resource "zentral_taxonomy" "test" {
  name = %[1]q
}

resource "zentral_tag" "test" {
  taxonomy_id = zentral_taxonomy.test.id
  name        = %[3]q
}

resource "zentral_tag" "test2" {
  taxonomy_id = zentral_taxonomy.test.id
  name        = %[4]q
}

resource "zentral_santa_scoped_path_regex" "test" {
  configuration_id        = zentral_santa_configuration.test.id
  name                    = %[2]q
  description             = "description"
  policy                  = "BLOCK"
  regex                   = "/Users/[^/]+/Downloads/.+"
  primary_users           = ["un", "deux"]
  excluded_primary_users  = ["trois"]
  serial_numbers          = ["quatre"]
  excluded_serial_numbers = ["cinq", "six"]
  tag_ids                 = [zentral_tag.test.id]
  excluded_tag_ids        = [zentral_tag.test2.id]
}
`, cfgName, name, tagName, tag2Name)
}
