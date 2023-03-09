package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonolithEnrollmentResource(t *testing.T) {
	name := acctest.RandString(12)
	tagName := acctest.RandString(12)
	resourceName := "zentral_monolith_enrollment.test"
	mResourceName := "zentral_monolith_manifest.test"
	mbuResourceName := "zentral_meta_business_unit.test"
	tagResourceName := "zentral_tag.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMonolithEnrollmentResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "manifest_id", mResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "meta_business_unit_id", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "serial_numbers.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "udids.#", "0"),
					resource.TestCheckNoResourceAttr(
						resourceName, "quota"),
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
				Config: testAccMonolithEnrollmentResourceConfigFull(name, tagName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "manifest_id", mResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "meta_business_unit_id", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "2"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "tag_ids.*", tagResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "serial_numbers.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "serial_numbers.*", "un"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "serial_numbers.*", "deux"),
					resource.TestCheckResourceAttr(
						resourceName, "udids.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "udids.*", "trois"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "udids.*", "quatre"),
					resource.TestCheckResourceAttr(
						resourceName, "quota", "5"),
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

func testAccMonolithEnrollmentResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_monolith_manifest" "test" {
  name                  = %[1]q
  meta_business_unit_id = zentral_meta_business_unit.test.id
}

resource "zentral_monolith_enrollment" "test" {
  manifest_id           = zentral_monolith_manifest.test.id
  meta_business_unit_id = zentral_meta_business_unit.test.id
}
`, name)
}

func testAccMonolithEnrollmentResourceConfigFull(name string, tagName string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_monolith_manifest" "test" {
  name                  = %[1]q
  meta_business_unit_id = zentral_meta_business_unit.test.id
}

resource "zentral_taxonomy" "test" {
  name = %[1]q
}

resource "zentral_tag" "test" {
  taxonomy_id = zentral_taxonomy.test.id
  name        = %[2]q
}

resource "zentral_monolith_enrollment" "test" {
  manifest_id           = zentral_monolith_manifest.test.id
  meta_business_unit_id = zentral_meta_business_unit.test.id
  tag_ids               = [zentral_tag.test.id]
  serial_numbers        = ["un", "deux"]
  udids                 = ["trois", "quatre"]
  quota                 = 5
}
`, name, tagName)
}
