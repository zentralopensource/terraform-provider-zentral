package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSantaRuleDataSource(t *testing.T) {
	name := acctest.RandString(12)
	tagName := acctest.RandString(12)
	tag2Name := acctest.RandString(12)
	resourceName := "zentral_santa_rule.test"
	dataSourceName := "data.zentral_santa_rule.test"
	cfgResourceName := "zentral_santa_configuration.test"
	tagResourceName := "zentral_tag.test"
	tag2ResourceName := "zentral_tag.test2"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSantaRuleDataSourceConfig(name, tagName, tag2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by id
					resource.TestCheckResourceAttrPair(
						dataSourceName, "configuration_id", cfgResourceName, "id"),
					resource.TestCheckResourceAttr(
						dataSourceName, "policy", "CEL"),
					resource.TestCheckResourceAttr(
						dataSourceName, "cel_expr", "target.signing_time >= timestamp('2025-05-31T00:00:00Z')"),
					resource.TestCheckResourceAttr(
						dataSourceName, "target_type", "CDHASH"),
					resource.TestCheckResourceAttr(
						dataSourceName, "target_identifier", "bff4a6a4d6b42e94e7d7f48e66b66c69b58fb409"),
					resource.TestCheckResourceAttr(
						dataSourceName, "description", "description"),
					resource.TestCheckResourceAttr(
						dataSourceName, "custom_message", "custom message"),
					resource.TestCheckNoResourceAttr(
						dataSourceName, "ruleset_id"),
					resource.TestCheckResourceAttr(
						dataSourceName, "primary_users.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						dataSourceName, "primary_users.*", "un"),
					resource.TestCheckTypeSetElemAttr(
						dataSourceName, "primary_users.*", "deux"),
					resource.TestCheckResourceAttr(
						dataSourceName, "excluded_primary_users.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						dataSourceName, "excluded_primary_users.*", "trois"),
					resource.TestCheckTypeSetElemAttr(
						dataSourceName, "excluded_primary_users.*", "quatre"),
					resource.TestCheckResourceAttr(
						dataSourceName, "serial_numbers.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						dataSourceName, "serial_numbers.*", "cinq"),
					resource.TestCheckTypeSetElemAttr(
						dataSourceName, "serial_numbers.*", "six"),
					resource.TestCheckResourceAttr(
						dataSourceName, "excluded_serial_numbers.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						dataSourceName, "excluded_serial_numbers.*", "sept"),
					resource.TestCheckTypeSetElemAttr(
						dataSourceName, "excluded_serial_numbers.*", "huit"),
					resource.TestCheckResourceAttr(
						dataSourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						dataSourceName, "tag_ids.*", tagResourceName, "id"),
					resource.TestCheckResourceAttr(
						dataSourceName, "excluded_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						dataSourceName, "excluded_tag_ids.*", tag2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						dataSourceName, "version", "1"),
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

func testAccSantaRuleDataSourceConfig(name string, tagName string, tag2Name string) string {
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
  target_type             = "CDHASH"
  target_identifier       = "bff4a6a4d6b42e94e7d7f48e66b66c69b58fb409"
  description             = "description"
  custom_message          = "custom message"
  primary_users           = ["un", "deux"]
  excluded_primary_users  = ["trois", "quatre"]
  serial_numbers          = ["cinq", "six"]
  excluded_serial_numbers = ["sept", "huit"]
  tag_ids                 = [zentral_tag.test.id]
  excluded_tag_ids        = [zentral_tag.test2.id]
}

data "zentral_santa_rule" "test" {
  id = zentral_santa_rule.test.id
}
`, name, tagName, tag2Name)
}
