package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Acceptance test for zentral_mdm_enterprise_app. The server downloads the
// package from package_uri and verifies package_sha256, so the test reads
// these from environment variables and skips when they are absent. Shared
// with TestAccMDMPackageResource:
//
//   ZTL_TEST_PACKAGE_URI       — a .pkg or .ipa URL reachable from the test Zentral
//   ZTL_TEST_PACKAGE_SHA256    — its SHA-256 hex digest
//
// The Full step swaps to a second package + bumps the artifact version, so it
// requires a second pair; when only the first pair is set, that step is
// skipped:
//
//   ZTL_TEST_PACKAGE_URI_2
//   ZTL_TEST_PACKAGE_SHA256_2

func TestAccMDMEnterpriseAppResource(t *testing.T) {
	pkgURI := os.Getenv("ZTL_TEST_PACKAGE_URI")
	pkgSHA256 := os.Getenv("ZTL_TEST_PACKAGE_SHA256")
	if pkgURI == "" || pkgSHA256 == "" {
		t.Skip("ZTL_TEST_PACKAGE_URI and ZTL_TEST_PACKAGE_SHA256 must be set to run this test")
	}
	pkgURI2 := os.Getenv("ZTL_TEST_PACKAGE_URI_2")
	pkgSHA2562 := os.Getenv("ZTL_TEST_PACKAGE_SHA256_2")

	name := acctest.RandString(12)
	aResourceName := "zentral_mdm_artifact.test"
	etResourceName := "zentral_tag.excluded"
	resourceName := "zentral_mdm_enterprise_app.test"

	steps := []resource.TestStep{
		// Create and Read
		{
			Config: testAccMDMEnterpriseAppResourceConfigBare(name, pkgURI, pkgSHA256),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPair(
					resourceName, "artifact_id", aResourceName, "id"),
				resource.TestCheckResourceAttr(
					resourceName, "package_uri", pkgURI),
				resource.TestCheckResourceAttr(
					resourceName, "package_sha256", pkgSHA256),
				resource.TestCheckResourceAttr(
					resourceName, "ios", "false"),
				resource.TestCheckResourceAttr(
					resourceName, "ios_max_version", ""),
				resource.TestCheckResourceAttr(
					resourceName, "ios_min_version", ""),
				resource.TestCheckResourceAttr(
					resourceName, "ipados", "false"),
				resource.TestCheckResourceAttr(
					resourceName, "ipados_max_version", ""),
				resource.TestCheckResourceAttr(
					resourceName, "ipados_min_version", ""),
				resource.TestCheckResourceAttr(
					resourceName, "macos", "true"),
				resource.TestCheckResourceAttr(
					resourceName, "macos_max_version", ""),
				resource.TestCheckResourceAttr(
					resourceName, "macos_min_version", ""),
				resource.TestCheckResourceAttr(
					resourceName, "tvos", "false"),
				resource.TestCheckResourceAttr(
					resourceName, "tvos_max_version", ""),
				resource.TestCheckResourceAttr(
					resourceName, "tvos_min_version", ""),
				resource.TestCheckResourceAttr(
					resourceName, "shard_modulo", "100"),
				resource.TestCheckResourceAttr(
					resourceName, "default_shard", "100"),
				resource.TestCheckResourceAttr(
					resourceName, "excluded_tag_ids.#", "0"),
				resource.TestCheckResourceAttr(
					resourceName, "tag_shards.#", "0"),
				resource.TestCheckResourceAttr(
					resourceName, "version", "1"),
			),
		},
		// ImportState
		{
			ResourceName:      resourceName,
			ImportState:       true,
			ImportStateVerify: true,
		},
	}

	// The Full step swaps the package and bumps the version. It needs a
	// distinct second package; skip silently when the _2 pair isn't set.
	if pkgURI2 != "" && pkgSHA2562 != "" {
		steps = append(steps,
			// Update and Read
			resource.TestStep{
				Config: testAccMDMEnterpriseAppResourceConfigFull(name, pkgURI2, pkgSHA2562),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "artifact_id", aResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "package_uri", pkgURI2),
					resource.TestCheckResourceAttr(
						resourceName, "package_sha256", pkgSHA2562),
					resource.TestCheckResourceAttr(
						resourceName, "ios", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "ios_max_version", ""),
					resource.TestCheckResourceAttr(
						resourceName, "ios_min_version", ""),
					resource.TestCheckResourceAttr(
						resourceName, "ipados", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "ipados_max_version", ""),
					resource.TestCheckResourceAttr(
						resourceName, "ipados_min_version", ""),
					resource.TestCheckResourceAttr(
						resourceName, "macos", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "macos_max_version", ""),
					resource.TestCheckResourceAttr(
						resourceName, "macos_min_version", "13.3.1"),
					resource.TestCheckResourceAttr(
						resourceName, "tvos", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "tvos_max_version", ""),
					resource.TestCheckResourceAttr(
						resourceName, "tvos_min_version", ""),
					resource.TestCheckResourceAttr(
						resourceName, "shard_modulo", "5"),
					resource.TestCheckResourceAttr(
						resourceName, "default_shard", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "excluded_tag_ids.*", etResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_shards.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"tag_shards.*",
						map[string]string{
							"shard": "4",
						},
					),
					resource.TestCheckResourceAttr(
						resourceName, "version", "2"),
				),
			},
			// ImportState
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

func testAccMDMEnterpriseAppResourceConfigBare(name, pkgURI, pkgSHA256 string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_artifact" "test" {
  name      = %[1]q
  type      = "Enterprise App"
  channel   = "Device"
  platforms = ["macOS"]
}

resource "zentral_mdm_enterprise_app" "test" {
  artifact_id    = zentral_mdm_artifact.test.id
  package_uri    = %[2]q
  package_sha256 = %[3]q
  version        = 1
  macos          = true
}
`, name, pkgURI, pkgSHA256)
}

func testAccMDMEnterpriseAppResourceConfigFull(name, pkgURI, pkgSHA256 string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_artifact" "test" {
  name      = %[1]q
  type      = "Enterprise App"
  channel   = "Device"
  platforms = ["macOS"]
}

resource "zentral_tag" "excluded" {
  name = "%[1]s excluded"
}

resource "zentral_tag" "shard" {
  name = "%[1]s shard"
}

resource "zentral_mdm_enterprise_app" "test" {
  artifact_id       = zentral_mdm_artifact.test.id
  package_uri       = %[2]q
  package_sha256    = %[3]q
  version           = 2
  macos             = true
  macos_min_version = "13.3.1"
  shard_modulo      = 5
  default_shard     = 1
  excluded_tag_ids  = [zentral_tag.excluded.id]
  tag_shards        = [{tag_id = zentral_tag.shard.id, shard=4}]
}
`, name, pkgURI, pkgSHA256)
}
