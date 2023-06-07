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
resource "zentral_mdm_blueprint" "test" {
  name                 = %[1]q
  inventory_interval   = 77777
  collect_apps         = "MANAGED_ONLY"
  collect_certificates = "ALL"
  collect_profiles     = "MANAGED_ONLY"
}
`, name)
}
