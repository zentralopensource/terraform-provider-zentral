package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOsqueryConfigurationPackResource(t *testing.T) {
	name := acctest.RandString(12)
	resourceName := "zentral_osquery_configuration_pack.test"
	cfgResourceName := "zentral_osquery_configuration.test"
	packResourceName := "zentral_osquery_pack.test"
	inclTagResourceName := "zentral_tag.included"
	exclTagResourceName := "zentral_tag.excluded"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccOsqueryConfigurationPackResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "configuration_id", cfgResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "pack_id", packResourceName, "id"),
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
				Config: testAccOsqueryConfigurationPackResourceConfigFull(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "configuration_id", cfgResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "pack_id", packResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "tag_ids.*", inclTagResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "excluded_tag_ids.*", exclTagResourceName, "id"),
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

func testAccOsqueryConfigurationPackResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_configuration" "test" {
  name =  %[1]q
}

resource "zentral_osquery_pack" "test" {
  name = %[1]q
}

resource "zentral_osquery_configuration_pack" "test" {
  configuration_id = zentral_osquery_configuration.test.id
  pack_id          = zentral_osquery_pack.test.id
}
`, name)
}

func testAccOsqueryConfigurationPackResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_configuration" "test" {
  name =  %[1]q
}

resource "zentral_osquery_pack" "test" {
  name = %[1]q
}

resource "zentral_tag" "included" {
  name = "%[1]s-included"
}

resource "zentral_tag" "excluded" {
  name = "%[1]s-excluded"
}

resource "zentral_osquery_configuration_pack" "test" {
  configuration_id = zentral_osquery_configuration.test.id
  pack_id          = zentral_osquery_pack.test.id
  tag_ids          = [zentral_tag.included.id]
  excluded_tag_ids = [zentral_tag.excluded.id]
}
`, name)
}
