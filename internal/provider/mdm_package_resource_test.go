package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Acceptance test for zentral_mdm_package. The server downloads the package
// from source_uri and verifies sha256, so the test cannot fabricate these — it
// reads them from environment variables and skips when they are absent.
//
// Required for the basic flow:
//   ZTL_TEST_PACKAGE_URI       — a .pkg or .ipa URL reachable from the test Zentral
//   ZTL_TEST_PACKAGE_SHA256    — its SHA-256 hex digest
//
// Optional, to exercise the RequiresReplace path on source_uri/sha256:
//   ZTL_TEST_PACKAGE_URI_2
//   ZTL_TEST_PACKAGE_SHA256_2

func TestAccMDMPackageResource(t *testing.T) {
	pkgURI := os.Getenv("ZTL_TEST_PACKAGE_URI")
	pkgSHA256 := os.Getenv("ZTL_TEST_PACKAGE_SHA256")
	if pkgURI == "" || pkgSHA256 == "" {
		t.Skip("ZTL_TEST_PACKAGE_URI and ZTL_TEST_PACKAGE_SHA256 must be set to run this test")
	}
	pkgURI2 := os.Getenv("ZTL_TEST_PACKAGE_URI_2")
	pkgSHA2562 := os.Getenv("ZTL_TEST_PACKAGE_SHA256_2")

	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	firstDescription := acctest.RandString(24)
	secondDescription := acctest.RandString(24)
	resourceName := "zentral_mdm_package.test"

	steps := []resource.TestStep{
		// Create and Read
		{
			Config: testAccMDMPackageResourceConfig(firstName, firstDescription, pkgURI, pkgSHA256),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(
					resourceName, "name", firstName),
				resource.TestCheckResourceAttr(
					resourceName, "description", firstDescription),
				resource.TestCheckResourceAttr(
					resourceName, "source_uri", pkgURI),
				resource.TestCheckResourceAttr(
					resourceName, "sha256", pkgSHA256),
				resource.TestCheckResourceAttrSet(
					resourceName, "type"),
				resource.TestCheckResourceAttrSet(
					resourceName, "size"),
				resource.TestCheckResourceAttrSet(
					resourceName, "filename"),
				resource.TestCheckResourceAttrSet(
					resourceName, "product_id"),
				resource.TestCheckResourceAttrSet(
					resourceName, "product_version"),
				resource.TestCheckResourceAttrSet(
					resourceName, "bundles"),
				resource.TestCheckResourceAttrSet(
					resourceName, "manifest"),
			),
		},
		// ImportState
		{
			ResourceName:      resourceName,
			ImportState:       true,
			ImportStateVerify: true,
		},
		// Update name + description (source_uri / sha256 unchanged)
		{
			Config: testAccMDMPackageResourceConfig(secondName, secondDescription, pkgURI, pkgSHA256),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(
					resourceName, "name", secondName),
				resource.TestCheckResourceAttr(
					resourceName, "description", secondDescription),
				resource.TestCheckResourceAttr(
					resourceName, "source_uri", pkgURI),
				resource.TestCheckResourceAttr(
					resourceName, "sha256", pkgSHA256),
			),
		},
		// ImportState
		{
			ResourceName:      resourceName,
			ImportState:       true,
			ImportStateVerify: true,
		},
	}

	// Optional RequiresReplace step: switching source_uri + sha256 destroys
	// the resource and recreates it (new id), because the server enforces
	// file immutability and the schema mirrors that with RequiresReplace.
	if pkgURI2 != "" && pkgSHA2562 != "" {
		steps = append(steps,
			resource.TestStep{
				Config: testAccMDMPackageResourceConfig(secondName, secondDescription, pkgURI2, pkgSHA2562),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "source_uri", pkgURI2),
					resource.TestCheckResourceAttr(
						resourceName, "sha256", pkgSHA2562),
				),
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

func testAccMDMPackageResourceConfig(name, description, sourceURI, sha256 string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_package" "test" {
  name        = %[1]q
  description = %[2]q
  source_uri  = %[3]q
  sha256      = %[4]q
}
`, name, description, sourceURI, sha256)
}
