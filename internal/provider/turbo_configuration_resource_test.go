package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTurboConfigurationResource(t *testing.T) {
	name := acctest.RandString(12)
	resourceName := "zentral_turbo_configuration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read with defaults
			{
				Config: testAccTurboConfigurationResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "collect_inventory", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_interval", "86400"),
					resource.TestCheckResourceAttr(
						resourceName, "default_check_interval", "86400"),
					resource.TestCheckResourceAttr(
						resourceName, "config_refresh_interval", "600"),
					resource.TestCheckResourceAttr(
						resourceName, "results_batch_size", "100"),
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
				Config: testAccTurboConfigurationResourceConfigFull(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", fmt.Sprintf("%s-updated", name)),
					resource.TestCheckResourceAttr(
						resourceName, "description", "a description"),
					resource.TestCheckResourceAttr(
						resourceName, "collect_inventory", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_interval", "3600"),
					resource.TestCheckResourceAttr(
						resourceName, "default_check_interval", "7200"),
					resource.TestCheckResourceAttr(
						resourceName, "config_refresh_interval", "300"),
					resource.TestCheckResourceAttr(
						resourceName, "results_batch_size", "50"),
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

func testAccTurboConfigurationResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_turbo_configuration" "test" {
  name = %[1]q
}
`, name)
}

func testAccTurboConfigurationResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_turbo_configuration" "test" {
  name                    = "%[1]s-updated"
  description             = "a description"
  collect_inventory       = false
  inventory_interval      = 3600
  default_check_interval  = 7200
  config_refresh_interval = 300
  results_batch_size      = 50
}
`, name)
}
