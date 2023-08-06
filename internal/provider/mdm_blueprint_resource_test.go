package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMBlueprintResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_mdm_blueprint.test"
	fcResourceName := "zentral_mdm_filevault_config.test"
	rpcResourceName := "zentral_mdm_recovery_password_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMBlueprintResourceConfigBare(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_interval", "86400"),
					resource.TestCheckResourceAttr(
						resourceName, "collect_apps", "NO"),
					resource.TestCheckResourceAttr(
						resourceName, "collect_certificates", "NO"),
					resource.TestCheckResourceAttr(
						resourceName, "collect_profiles", "NO"),
					resource.TestCheckNoResourceAttr(
						resourceName, "filevault_config_id"),
					resource.TestCheckNoResourceAttr(
						resourceName, "recovery_password_config_id"),
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
				Config: testAccMDMBlueprintResourceConfigFull(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_interval", "77777"),
					resource.TestCheckResourceAttr(
						resourceName, "collect_apps", "MANAGED_ONLY"),
					resource.TestCheckResourceAttr(
						resourceName, "collect_certificates", "ALL"),
					resource.TestCheckResourceAttr(
						resourceName, "collect_profiles", "MANAGED_ONLY"),
					resource.TestCheckResourceAttrPair(
						resourceName, "filevault_config_id", fcResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "recovery_password_config_id", rpcResourceName, "id"),
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

func testAccMDMBlueprintResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_blueprint" "test" {
  name = %[1]q
}
`, name)
}

func testAccMDMBlueprintResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_filevault_config" "test" {
  name                         = %[1]q
  escrow_location_display_name = %[1]q
}

resource "zentral_mdm_recovery_password_config" "test" {
  name = %[1]q
}

resource "zentral_mdm_blueprint" "test" {
  name                        = %[1]q
  inventory_interval          = 77777
  collect_apps                = "MANAGED_ONLY"
  collect_certificates        = "ALL"
  collect_profiles            = "MANAGED_ONLY"
  filevault_config_id         = zentral_mdm_filevault_config.test.id
  recovery_password_config_id = zentral_mdm_recovery_password_config.test.id
}
`, name)
}
