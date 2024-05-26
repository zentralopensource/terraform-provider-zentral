package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonolithManifestEnrollmentPackageResource(t *testing.T) {
	name := acctest.RandString(12)
	resourceName := "zentral_monolith_manifest_enrollment_package.test"
	mResourceName := "zentral_monolith_manifest.test"
	eResourceName := "zentral_munki_enrollment.test"
	tResourceName := "zentral_tag.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMonolithManifestEnrollmentPackageResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "manifest_id", mResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "builder", "zentral.contrib.munki.osx_package.builder.MunkiZentralEnrollPkgBuilder"),
					resource.TestCheckResourceAttrPair(
						resourceName, "enrollment_id", eResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "1"),
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
				Config: testAccMonolithManifestEnrollmentPackageResourceConfigFull(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "manifest_id", mResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "builder", "zentral.contrib.munki.osx_package.builder.MunkiZentralEnrollPkgBuilder"),
					resource.TestCheckResourceAttrPair(
						resourceName, "enrollment_id", eResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "2"),
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

func testAccMonolithManifestEnrollmentPackageResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_monolith_manifest" "test" {
  name                  =  %[1]q
  meta_business_unit_id = zentral_meta_business_unit.test.id
}

resource "zentral_munki_configuration" "test" {
  name = %[1]q
}

resource "zentral_munki_enrollment" "test" {
  configuration_id      = zentral_munki_configuration.test.id
  meta_business_unit_id = zentral_meta_business_unit.test.id
}

resource "zentral_monolith_manifest_enrollment_package" "test" {
  manifest_id   = zentral_monolith_manifest.test.id
  builder       = "zentral.contrib.munki.osx_package.builder.MunkiZentralEnrollPkgBuilder"
  enrollment_id = zentral_munki_enrollment.test.id
}
`, name)
}

func testAccMonolithManifestEnrollmentPackageResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_monolith_manifest" "test" {
  name                  =  %[1]q
  meta_business_unit_id = zentral_meta_business_unit.test.id
}

resource "zentral_tag" "test" {
  name = %[1]q
}

resource "zentral_munki_configuration" "test" {
  name = %[1]q
}

resource "zentral_munki_enrollment" "test" {
  configuration_id      = zentral_munki_configuration.test.id
  meta_business_unit_id = zentral_meta_business_unit.test.id
}

resource "zentral_monolith_manifest_enrollment_package" "test" {
  manifest_id   = zentral_monolith_manifest.test.id
  builder       = "zentral.contrib.munki.osx_package.builder.MunkiZentralEnrollPkgBuilder"
  enrollment_id = zentral_munki_enrollment.test.id
  tag_ids       = [zentral_tag.test.id]
}
`, name)
}
