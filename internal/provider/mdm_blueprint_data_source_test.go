package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMBlueprintDataSource(t *testing.T) {
	c1Name := acctest.RandString(12)
	c2Name := acctest.RandString(12)
	c1ResourceName := "zentral_mdm_blueprint.check1"
	c2ResourceName := "zentral_mdm_blueprint.check2"
	fc2ResourceName := "zentral_mdm_filevault_config.check2"
	rpc2ResourceName := "zentral_mdm_recovery_password_config.check2"
	sue2ResourceName := "zentral_mdm_software_update_enforcement.check2"
	ds1ResourceName := "data.zentral_mdm_blueprint.check1_by_name"
	ds2ResourceName := "data.zentral_mdm_blueprint.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMBlueprintDataSourceConfig(c1Name, c2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", c1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", c1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "inventory_interval", "86400"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "collect_apps", "NO"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "collect_certificates", "NO"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "collect_profiles", "NO"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "filevault_config_id"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "recovery_password_config_id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "software_update_enforcement_ids.#", "0"),
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", c2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", c2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "inventory_interval", "77777"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "collect_apps", "MANAGED_ONLY"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "collect_certificates", "ALL"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "collect_profiles", "MANAGED_ONLY"),
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "filevault_config_id", fc2ResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "recovery_password_config_id", rpc2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "software_update_enforcement_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						ds2ResourceName, "software_update_enforcement_ids.*", sue2ResourceName, "id"),
				),
			},
		},
	})
}

func testAccMDMBlueprintDataSourceConfig(c1Name string, c2Name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_blueprint" "check1" {
  name = %[1]q
}

resource "zentral_mdm_filevault_config" "check2" {
  name                         = %[2]q
  escrow_location_display_name = %[2]q
}

resource "zentral_mdm_recovery_password_config" "check2" {
  name = %[2]q
}

resource "zentral_mdm_software_update_enforcement" "check2" {
  name           = %[1]q
  platforms      = ["macOS"]
  max_os_version = "15"
}

resource "zentral_mdm_blueprint" "check2" {
  name                            = %[2]q
  inventory_interval              = 77777
  collect_apps                    = "MANAGED_ONLY"
  collect_certificates            = "ALL"
  collect_profiles                = "MANAGED_ONLY"
  filevault_config_id             = zentral_mdm_filevault_config.check2.id
  recovery_password_config_id     = zentral_mdm_recovery_password_config.check2.id
  software_update_enforcement_ids = [zentral_mdm_software_update_enforcement.check2.id]
}

data "zentral_mdm_blueprint" "check1_by_name" {
  name = zentral_mdm_blueprint.check1.name
}

data "zentral_mdm_blueprint" "check2_by_id" {
  id = zentral_mdm_blueprint.check2.id
}
`, c1Name, c2Name)
}
