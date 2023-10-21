package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMunkiConfigurationDataSource(t *testing.T) {
	c1Name := acctest.RandString(12)
	c2Name := acctest.RandString(12)
	c1ResourceName := "zentral_munki_configuration.check1"
	c2ResourceName := "zentral_munki_configuration.check2"
	ds1ResourceName := "data.zentral_munki_configuration.check1_by_name"
	ds2ResourceName := "data.zentral_munki_configuration.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMunkiConfigurationDataSourceConfig(c1Name, c2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name, platforms, tag
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", c1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", c1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "description", ""),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "inventory_apps_full_info_shard", "100"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "principal_user_detection_sources.#", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "principal_user_detection_domains.#", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "collected_condition_keys.#", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "managed_installs_sync_interval_days", "7"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "script_checks_run_interval_seconds", "86400"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "auto_reinstall_incidents", "false"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "auto_failed_install_incidents", "false"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "version", "0"),
					// Read by ID, no platforms, no tags
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", c2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", c2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "description", "Description"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "inventory_apps_full_info_shard", "50"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "principal_user_detection_sources.#", "2"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "principal_user_detection_sources.0", "google_chrome"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "principal_user_detection_sources.1", "company_portal"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "principal_user_detection_domains.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "principal_user_detection_domains.*", "zentral.io"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "collected_condition_keys.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "collected_condition_keys.*", "arch"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "collected_condition_keys.*", "machine_type"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "managed_installs_sync_interval_days", "3"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "script_checks_run_interval_seconds", "12345"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "auto_reinstall_incidents", "true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "auto_failed_install_incidents", "true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "version", "0"),
				),
			},
		},
	})
}

func testAccMunkiConfigurationDataSourceConfig(c1Name string, c2Name string) string {
	return fmt.Sprintf(`
resource "zentral_munki_configuration" "check1" {
  name = %[1]q
}

resource "zentral_munki_configuration" "check2" {
  name                                = %[2]q
  description                         = "Description"
  inventory_apps_full_info_shard      = 50
  principal_user_detection_sources    = ["google_chrome", "company_portal"]
  principal_user_detection_domains    = ["zentral.io"]
  collected_condition_keys            = ["arch", "machine_type"]
  managed_installs_sync_interval_days = 3
  script_checks_run_interval_seconds  = 12345
  auto_reinstall_incidents            = true
  auto_failed_install_incidents       = true
}

data "zentral_munki_configuration" "check1_by_name" {
  name = zentral_munki_configuration.check1.name
}

data "zentral_munki_configuration" "check2_by_id" {
  id = zentral_munki_configuration.check2.id
}
`, c1Name, c2Name)
}
