package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonolithManifestCatalogResource(t *testing.T) {
	name := acctest.RandString(12)
	resourceName := "zentral_monolith_manifest_catalog.test"
	mResourceName := "zentral_monolith_manifest.test"
	cResourceName := "zentral_monolith_catalog.test"
	tResourceName := "zentral_tag.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMonolithManifestCatalogResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "manifest_id", mResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "catalog_id", cResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "0"),
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
				Config: testAccMonolithManifestCatalogResourceConfigFull(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "manifest_id", mResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "catalog_id", cResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "tag_ids.*", tResourceName, "id"),
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

func testAccMonolithManifestCatalogResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_monolith_manifest" "test" {
  name                  =  %[1]q
  meta_business_unit_id = zentral_meta_business_unit.test.id
}

resource "zentral_monolith_catalog" "test" {
  name = %[1]q
}

resource "zentral_monolith_manifest_catalog" "test" {
  manifest_id = zentral_monolith_manifest.test.id
  catalog_id  = zentral_monolith_catalog.test.id
}
`, name)
}

func testAccMonolithManifestCatalogResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_monolith_manifest" "test" {
  name                  =  %[1]q
  meta_business_unit_id = zentral_meta_business_unit.test.id
}

resource "zentral_monolith_catalog" "test" {
  name = %[1]q
}

resource "zentral_tag" "test" {
  name = %[1]q
}

resource "zentral_monolith_manifest_catalog" "test" {
  manifest_id = zentral_monolith_manifest.test.id
  catalog_id  = zentral_monolith_catalog.test.id
  tag_ids     = [zentral_tag.test.id]
}
`, name)
}
