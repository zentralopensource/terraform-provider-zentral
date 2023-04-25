package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSantaConfigurationDataSource(t *testing.T) {
	c1Name := acctest.RandString(12)
	c2Name := acctest.RandString(12)
	c1ResourceName := "zentral_santa_configuration.check1"
	c2ResourceName := "zentral_santa_configuration.check2"
	ds1ResourceName := "data.zentral_santa_configuration.check1_by_name"
	ds2ResourceName := "data.zentral_santa_configuration.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSantaConfigurationDataSourceConfig(c1Name, c2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name, platforms, tag
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", c1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", c1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "client_mode", "MONITOR"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "client_certificate_auth", "false"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "batch_size", "50"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "full_sync_interval", "600"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "enable_bundles", "false"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "enable_transitive_rules", "false"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "allowed_path_regex", ""),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "blocked_path_regex", ""),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "block_usb_mount", "false"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "remount_usb_mode.#", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "allow_unknown_shard", "100"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "enable_all_event_upload_shard", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "sync_incident_severity", "0"),
					// Read by ID, no platforms, no tags
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", c2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", c2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "client_mode", "LOCKDOWN"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "client_certificate_auth", "false"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "batch_size", "50"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "full_sync_interval", "600"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "enable_bundles", "false"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "enable_transitive_rules", "false"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "allowed_path_regex", ""),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "blocked_path_regex", ""),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "block_usb_mount", "true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "remount_usb_mode.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "remount_usb_mode.*", "rdonly"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "remount_usb_mode.*", "noexec"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "allow_unknown_shard", "100"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "enable_all_event_upload_shard", "0"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "sync_incident_severity", "0"),
				),
			},
		},
	})
}

func testAccSantaConfigurationDataSourceConfig(c1Name string, c2Name string) string {
	return fmt.Sprintf(`
resource "zentral_santa_configuration" "check1" {
  name = %q
}

resource "zentral_santa_configuration" "check2" {
  name             = %q
  client_mode      = "LOCKDOWN"
  block_usb_mount  = true
  remount_usb_mode = ["noexec", "rdonly"]
}

data "zentral_santa_configuration" "check1_by_name" {
  name = zentral_santa_configuration.check1.name
}

data "zentral_santa_configuration" "check2_by_id" {
  id = zentral_santa_configuration.check2.id
}
`, c1Name, c2Name)
}
