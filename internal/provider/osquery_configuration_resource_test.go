package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOsqueryConfigurationResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_osquery_configuration.test"
	atcResourceName := "zentral_osquery_atc.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccOsqueryConfigurationResourceConfigBare(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "inventory", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_apps", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_ec2", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_interval", "86400"),
					resource.TestCheckResourceAttr(
						resourceName, "options.%", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "automatic_table_constructions.#", "0"),
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
				Config: testAccOsqueryConfigurationResourceConfigFull(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "description"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_apps", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_ec2", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_interval", "600"),
					resource.TestCheckResourceAttr(
						resourceName, "options.%", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "options.config_refresh", "120"),
					resource.TestCheckResourceAttr(
						resourceName, "automatic_table_constructions.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "automatic_table_constructions.*", atcResourceName, "id"),
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

func testAccOsqueryConfigurationResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_configuration" "test" {
  name = %[1]q
}
`, name)
}

func testAccOsqueryConfigurationResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_atc" "test" {
  name        = "%[1]s_atc"
  description = "Access the Google Santa rules.db"
  table_name  = "%[1]s_test_tf_santa_rules"
  query       = "SELECT * FROM rules;"
  path        = "/var/db/santa/rules.db"
  columns     = ["identifier", "state", "type", "custommsg", "timestamp"]
  platforms   = ["darwin"]
}

resource "zentral_osquery_configuration" "test" {
  name                          = %[1]q
  description                   = "description"
  inventory                     = true
  inventory_apps                = true
  inventory_ec2                 = true
  inventory_interval            = 600
  options                       = { config_refresh = "120" }
  automatic_table_constructions = [zentral_osquery_atc.test.id]
}
`, name)
}
