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
	rResourceName := "zentral_monolith_repository.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMonolithCatalogResourceConfigBare(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "repository_id", rResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
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
					resource.TestCheckResourceAttrPair(
						resourceName, "repository_id", rResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
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
resource "zentral_monolith_repository" "test" {
  name    = %[1]q
  backend = "VIRTUAL"
}

resource "zentral_monolith_catalog" "test" {
  repository_id = zentral_monolith_repository.test.id
  name          = %[1]q
}
`, name)
}

func testAccMonolithCatalogResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_monolith_repository" "test" {
  name    = %[1]q
  backend = "VIRTUAL"
}

resource "zentral_monolith_catalog" "test" {
  repository_id = zentral_monolith_repository.test.id
  name     = %[1]q
}
`, name)
}
