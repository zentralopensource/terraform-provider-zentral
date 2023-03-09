package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMunkiConfigurationResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_munki_configuration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMunkiConfigurationResourceConfigBare(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_apps_full_info_shard", "100"),
					resource.TestCheckResourceAttr(
						resourceName, "principal_user_detection_sources.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "principal_user_detection_domains.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "collected_condition_keys.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "managed_installs_sync_interval_days", "7"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_reinstall_incidents", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_failed_install_incidents", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "0"),
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
				Config: testAccMunkiConfigurationResourceConfigFull(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "Description"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_apps_full_info_shard", "50"),
					resource.TestCheckResourceAttr(
						resourceName, "principal_user_detection_sources.#", "2"),
					resource.TestCheckResourceAttr(
						resourceName, "principal_user_detection_sources.0", "google_chrome"),
					resource.TestCheckResourceAttr(
						resourceName, "principal_user_detection_sources.1", "company_portal"),
					resource.TestCheckResourceAttr(
						resourceName, "principal_user_detection_domains.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "principal_user_detection_domains.*", "zentral.io"),
					resource.TestCheckResourceAttr(
						resourceName, "collected_condition_keys.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "collected_condition_keys.*", "arch"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "collected_condition_keys.*", "machine_type"),
					resource.TestCheckResourceAttr(
						resourceName, "managed_installs_sync_interval_days", "3"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_reinstall_incidents", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_failed_install_incidents", "true"),
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
		},
	})
}

func testAccMunkiConfigurationResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_munki_configuration" "test" {
  name = %[1]q
}
`, name)
}

func testAccMunkiConfigurationResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_munki_configuration" "test" {
  name                                = %[1]q
  description                         = "Description"
  inventory_apps_full_info_shard      = 50
  principal_user_detection_sources    = ["google_chrome", "company_portal"]
  principal_user_detection_domains    = ["zentral.io"]
  collected_condition_keys            = ["arch", "machine_type"]
  managed_installs_sync_interval_days = 3
  auto_reinstall_incidents            = true
  auto_failed_install_incidents       = true
}
`, name)
}
