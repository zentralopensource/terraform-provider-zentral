package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSantaConfigurationResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_santa_configuration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccSantaConfigurationResourceConfigBare(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "client_mode", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "client_certificate_auth", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "batch_size", "50"),
					resource.TestCheckResourceAttr(
						resourceName, "full_sync_interval", "600"),
					resource.TestCheckResourceAttr(
						resourceName, "enable_bundles", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "enable_transitive_rules", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "allowed_path_regex", ""),
					resource.TestCheckResourceAttr(
						resourceName, "blocked_path_regex", ""),
					resource.TestCheckResourceAttr(
						resourceName, "block_usb_mount", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "remount_usb_mode.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "allow_unknown_shard", "100"),
					resource.TestCheckResourceAttr(
						resourceName, "enable_all_event_upload_shard", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "sync_incident_severity", "0"),
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
				Config: testAccSantaConfigurationResourceConfigFull(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "client_mode", "2"),
					resource.TestCheckResourceAttr(
						resourceName, "client_certificate_auth", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "batch_size", "51"),
					resource.TestCheckResourceAttr(
						resourceName, "full_sync_interval", "601"),
					resource.TestCheckResourceAttr(
						resourceName, "enable_bundles", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "enable_transitive_rules", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "allowed_path_regex", "un"),
					resource.TestCheckResourceAttr(
						resourceName, "blocked_path_regex", "deux"),
					resource.TestCheckResourceAttr(
						resourceName, "block_usb_mount", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "remount_usb_mode.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "remount_usb_mode.*", "rdonly"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "remount_usb_mode.*", "noexec"),
					resource.TestCheckResourceAttr(
						resourceName, "allow_unknown_shard", "99"),
					resource.TestCheckResourceAttr(
						resourceName, "enable_all_event_upload_shard", "10"),
					resource.TestCheckResourceAttr(
						resourceName, "sync_incident_severity", "100"),
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

func testAccSantaConfigurationResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_santa_configuration" "test" {
  name = %[1]q
}
`, name)
}

func testAccSantaConfigurationResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_santa_configuration" "test" {
  name                          = %[1]q
  client_mode                   = 2
  client_certificate_auth       = true
  batch_size                    = 51
  full_sync_interval            = 601
  enable_bundles                = true
  enable_transitive_rules       = true
  allowed_path_regex            = "un"
  blocked_path_regex            = "deux"
  block_usb_mount               = true
  remount_usb_mode              = ["noexec", "rdonly"]
  allow_unknown_shard           = 99
  enable_all_event_upload_shard = 10
  sync_incident_severity        = 100
}
`, name)
}
