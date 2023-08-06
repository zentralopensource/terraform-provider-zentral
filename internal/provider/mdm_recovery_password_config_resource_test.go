package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMRecoveryPasswordConfigResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_mdm_recovery_password_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMRecoveryPasswordConfigResourceConfigDynamic(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "dynamic_password", "true"),
					resource.TestCheckNoResourceAttr(
						resourceName, "static_password"),
					resource.TestCheckResourceAttr(
						resourceName, "rotation_interval_days", "90"),
					resource.TestCheckResourceAttr(
						resourceName, "rotate_firmware_password", "true"),
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
				Config: testAccMDMRecoveryPasswordConfigResourceConfigStatic(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "dynamic_password", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "static_password", "12345678"),
					resource.TestCheckResourceAttr(
						resourceName, "rotation_interval_days", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "rotate_firmware_password", "false"),
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

func testAccMDMRecoveryPasswordConfigResourceConfigDynamic(name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_recovery_password_config" "test" {
  name                     = %[1]q
  rotation_interval_days   = 90
  rotate_firmware_password = true
}
`, name)
}

func testAccMDMRecoveryPasswordConfigResourceConfigStatic(name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_recovery_password_config" "test" {
  name                     = %[1]q
  dynamic_password         = false
  static_password          = "12345678"
}
`, name)
}
