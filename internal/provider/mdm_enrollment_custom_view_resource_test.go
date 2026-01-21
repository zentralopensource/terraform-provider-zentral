package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMEnrollmentCustomViewResource(t *testing.T) {
	name := acctest.RandString(12)
	description := acctest.RandString(12)
	resourceName := "zentral_mdm_enrollment_custom_view.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMEnrollmentCustomViewResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "html", "<html/>"),
					resource.TestCheckResourceAttr(
						resourceName, "requires_authentication", "false"),
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
				Config: testAccMDMEnrollmentCustomViewResourceConfigFull(name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName, "description", description),
					resource.TestCheckResourceAttr(
						resourceName, "html", "<html><body>updated</body></html>"),
					resource.TestCheckResourceAttr(
						resourceName, "requires_authentication", "true"),
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

func testAccMDMEnrollmentCustomViewResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_enrollment_custom_view" "test" {
  name = %[1]q
  html = "<html/>"
}
`, name)
}

func testAccMDMEnrollmentCustomViewResourceConfigFull(name string, description string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_enrollment_custom_view" "test" {
  name                    = %[1]q
  description             = %[2]q
  html                    = "<html><body>updated</body></html>"
  requires_authentication = true
}
`, name, description)
}
