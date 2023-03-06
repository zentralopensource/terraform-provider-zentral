package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonolithCatalogResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_monolith_catalog.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMonolithCatalogResourceConfigBare(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "priority", "0"),
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
				Config: testAccMonolithCatalogResourceConfigFull(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "priority", "17"),
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

func testAccMonolithCatalogResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_monolith_catalog" "test" {
  name = %[1]q
}
`, name)
}

func testAccMonolithCatalogResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_monolith_catalog" "test" {
  name     = %[1]q
  priority = 17
}
`, name)
}
