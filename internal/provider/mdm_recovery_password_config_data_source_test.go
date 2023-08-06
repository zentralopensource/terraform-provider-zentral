package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMRecoveryPasswordConfigDataSource(t *testing.T) {
	c1Name := acctest.RandString(12)
	c2Name := acctest.RandString(12)
	c1ResourceName := "zentral_mdm_recovery_password_config.check1"
	c2ResourceName := "zentral_mdm_recovery_password_config.check2"
	ds1ResourceName := "data.zentral_mdm_recovery_password_config.check1_by_name"
	ds2ResourceName := "data.zentral_mdm_recovery_password_config.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMRecoveryPasswordConfigDataSourceConfig(c1Name, c2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", c1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "dynamic_password", "true"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "static_password"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "rotation_interval_days", "90"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "rotate_firmware_password", "true"),
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", c2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "dynamic_password", "false"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "static_password", "12345678"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "rotation_interval_days", "0"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "rotate_firmware_password", "false"),
				),
			},
		},
	})
}

func testAccMDMRecoveryPasswordConfigDataSourceConfig(c1Name string, c2Name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_recovery_password_config" "check1" {
  name                     = %[1]q
  rotation_interval_days   = 90
  rotate_firmware_password = true
}

resource "zentral_mdm_recovery_password_config" "check2" {
  name                     = %[2]q
  dynamic_password         = false
  static_password          = "12345678"
}

data "zentral_mdm_recovery_password_config" "check1_by_name" {
  name = zentral_mdm_recovery_password_config.check1.name
}

data "zentral_mdm_recovery_password_config" "check2_by_id" {
  id = zentral_mdm_recovery_password_config.check2.id
}
`, c1Name, c2Name)
}
