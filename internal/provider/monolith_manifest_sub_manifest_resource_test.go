package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonolithManifestSubManifestResource(t *testing.T) {
	name := acctest.RandString(12)
	resourceName := "zentral_monolith_manifest_sub_manifest.test"
	mResourceName := "zentral_monolith_manifest.test"
	smResourceName := "zentral_monolith_sub_manifest.test"
	tResourceName := "zentral_tag.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMonolithManifestSubManifestResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "manifest_id", mResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "sub_manifest_id", smResourceName, "id"),
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
				Config: testAccMonolithManifestSubManifestResourceConfigFull(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "manifest_id", mResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "sub_manifest_id", smResourceName, "id"),
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

func testAccMonolithManifestSubManifestResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_monolith_manifest" "test" {
  name                  =  %[1]q
  meta_business_unit_id = zentral_meta_business_unit.test.id
}

resource "zentral_monolith_sub_manifest" "test" {
  name = %[1]q
}

resource "zentral_monolith_manifest_sub_manifest" "test" {
  manifest_id      = zentral_monolith_manifest.test.id
  sub_manifest_id  = zentral_monolith_sub_manifest.test.id
}
`, name)
}

func testAccMonolithManifestSubManifestResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_monolith_manifest" "test" {
  name                  =  %[1]q
  meta_business_unit_id = zentral_meta_business_unit.test.id
}

resource "zentral_monolith_sub_manifest" "test" {
  name = %[1]q
}

resource "zentral_tag" "test" {
  name = %[1]q
}

resource "zentral_monolith_manifest_sub_manifest" "test" {
  manifest_id      = zentral_monolith_manifest.test.id
  sub_manifest_id  = zentral_monolith_sub_manifest.test.id
  tag_ids          = [zentral_tag.test.id]
}
`, name)
}
