package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMFileVaultConfigDataSource(t *testing.T) {
	c1Name := acctest.RandString(12)
	c2Name := acctest.RandString(12)
	c1ResourceName := "zentral_mdm_filevault_config.check1"
	c2ResourceName := "zentral_mdm_filevault_config.check2"
	ds1ResourceName := "data.zentral_mdm_filevault_config.check1_by_name"
	ds2ResourceName := "data.zentral_mdm_filevault_config.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMFileVaultConfigDataSourceConfig(c1Name, c2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", c1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "escrow_location_display_name", c1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "at_login_only", "false"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "bypass_attempts", "-1"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "show_recovery_key", "false"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "destroy_key_on_standby", "false"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "prk_rotation_interval_days", "0"),
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", c2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "escrow_location_display_name", c2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "at_login_only", "true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "bypass_attempts", "1"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "show_recovery_key", "true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "destroy_key_on_standby", "true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "prk_rotation_interval_days", "90"),
				),
			},
		},
	})
}

func testAccMDMFileVaultConfigDataSourceConfig(c1Name string, c2Name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_filevault_config" "check1" {
  name                         = %[1]q
  escrow_location_display_name = %[1]q
}

resource "zentral_mdm_filevault_config" "check2" {
  name                         = %[2]q
  escrow_location_display_name = %[2]q
  at_login_only                = true
  bypass_attempts              = 1
  show_recovery_key            = true
  destroy_key_on_standby       = true
  prk_rotation_interval_days   = 90
}

data "zentral_mdm_filevault_config" "check1_by_name" {
  name = zentral_mdm_filevault_config.check1.name
}

data "zentral_mdm_filevault_config" "check2_by_id" {
  id = zentral_mdm_filevault_config.check2.id
}
`, c1Name, c2Name)
}
