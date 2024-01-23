package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonolithCatalogDataSource(t *testing.T) {
	c1Name := acctest.RandString(12)
	c2Name := acctest.RandString(12)
	c1ResourceName := "zentral_monolith_catalog.test1"
	c2ResourceName := "zentral_monolith_catalog.test2"
	ds1ResourceName := "data.zentral_monolith_catalog.test_by_id"
	ds2ResourceName := "data.zentral_monolith_catalog.test_by_name"
	r1ResourceName := "zentral_monolith_repository.test1"
	r2ResourceName := "zentral_monolith_repository.test2"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonolithCatalogDataSourceConfig(c1Name, c2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by id
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", c1ResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "repository_id", r1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", c1Name),
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", c2ResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "repository_id", r2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", c2Name),
				),
			},
		},
	})
}

func testAccMonolithCatalogDataSourceConfig(c1Name string, c2Name string) string {
	return fmt.Sprintf(`
resource "zentral_monolith_repository" "test1" {
  name    = %[1]q
  backend = "VIRTUAL"
}

resource "zentral_monolith_catalog" "test1" {
  repository_id = zentral_monolith_repository.test1.id
  name          = %[1]q
}

resource "zentral_monolith_repository" "test2" {
  name    = %[2]q
  backend = "VIRTUAL"
}

resource "zentral_monolith_catalog" "test2" {
  repository_id = zentral_monolith_repository.test2.id
  name          = %[2]q
}

data "zentral_monolith_catalog" "test_by_id" {
  id = zentral_monolith_catalog.test1.id
}

data "zentral_monolith_catalog" "test_by_name" {
  name = zentral_monolith_catalog.test2.name
}
`, c1Name, c2Name)
}
