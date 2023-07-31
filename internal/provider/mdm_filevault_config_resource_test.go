package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMFileVaultConfigResource(t *testing.T) {
	firstName := acctest.RandString(12)
	firstEscName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	secondEscName := acctest.RandString(12)
	resourceName := "zentral_mdm_filevault_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMFileVaultConfigResourceConfigBare(firstName, firstEscName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "escrow_location_display_name", firstEscName),
					resource.TestCheckResourceAttr(
						resourceName, "at_login_only", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "bypass_attempts", "-1"),
					resource.TestCheckResourceAttr(
						resourceName, "show_recovery_key", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "destroy_key_on_standby", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "prk_rotation_interval_days", "0"),
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
				Config: testAccMDMFileVaultConfigResourceConfigFull(secondName, secondEscName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "escrow_location_display_name", secondEscName),
					resource.TestCheckResourceAttr(
						resourceName, "at_login_only", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "bypass_attempts", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "show_recovery_key", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "destroy_key_on_standby", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "prk_rotation_interval_days", "90"),
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

func testAccMDMFileVaultConfigResourceConfigBare(name string, escName string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_filevault_config" "test" {
  name                         = %[1]q
  escrow_location_display_name = %[2]q
}
`, name, escName)
}

func testAccMDMFileVaultConfigResourceConfigFull(name string, escName string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_filevault_config" "test" {
  name                         = %[1]q
  escrow_location_display_name = %[2]q
  at_login_only                = true
  bypass_attempts              = 1
  show_recovery_key            = true
  destroy_key_on_standby       = true
  prk_rotation_interval_days   = 90
}
`, name, escName)
}
