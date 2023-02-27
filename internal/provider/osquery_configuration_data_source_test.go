package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOsqueryConfigurationDataSource(t *testing.T) {
	c1Name := acctest.RandString(12)
	c2Name := acctest.RandString(12)
	c1ResourceName := "zentral_osquery_configuration.check1"
	c2ResourceName := "zentral_osquery_configuration.check2"
	ds1ResourceName := "data.zentral_osquery_configuration.check1_by_name"
	ds2ResourceName := "data.zentral_osquery_configuration.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOsqueryConfigurationDataSourceConfig(c1Name, c2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", c1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", c1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "description", ""),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "inventory", "true"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "inventory_apps", "false"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "inventory_ec2", "false"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "inventory_interval", "86400"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "options.%", "0"),
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", c2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", c2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "description", "description"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "inventory", "true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "inventory_apps", "true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "inventory_ec2", "true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "inventory_interval", "600"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "options.%", "1"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "options.config_refresh", "120"),
				),
			},
		},
	})
}

func testAccOsqueryConfigurationDataSourceConfig(c1Name string, c2Name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_configuration" "check1" {
  name = %q
}

resource "zentral_osquery_configuration" "check2" {
  name               = %q
  description        = "description"
  inventory          = true
  inventory_apps     = true
  inventory_ec2      = true
  inventory_interval = 600
  options            = { config_refresh = "120" }
}

data "zentral_osquery_configuration" "check1_by_name" {
  name = zentral_osquery_configuration.check1.name
}

data "zentral_osquery_configuration" "check2_by_id" {
  id = zentral_osquery_configuration.check2.id
}
`, c1Name, c2Name)
}
