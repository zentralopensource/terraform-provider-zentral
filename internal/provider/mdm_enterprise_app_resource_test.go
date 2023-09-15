package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMEnterpriseAppResource(t *testing.T) {
	name := acctest.RandString(12)
	aResourceName := "zentral_mdm_artifact.test"
	etResourceName := "zentral_tag.excluded"
	resourceName := "zentral_mdm_enterprise_app.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMEnterpriseAppResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "artifact_id", aResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "package_uri", "s3://zentral-pro-services-artifacts-eu-central-1/terraform-provider-zentral/testdata/Outset-4.0.21900.pkg"),
					resource.TestCheckResourceAttr(
						resourceName, "package_sha256", "6e364e261e15ffdae89e6dfbd1b9e583319cf897d2dfe624521f5d0ef43b0be7"),
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
			// Update and Read
			{
				Config: testAccMDMEnterpriseAppResourceConfigFull(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "artifact_id", aResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "package_uri", "s3://zentral-pro-services-artifacts-eu-central-1/terraform-provider-zentral/testdata/Outset-4.1.0.21912.pkg"),
					resource.TestCheckResourceAttr(
						resourceName, "package_sha256", "d26590460ccee1895c999d3275475a1a9006366500cea6b4500b731e39d3cc2a"),
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMDMEnterpriseAppResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_artifact" "test" {
  name      = %[1]q
  type      = "Enterprise App"
  channel   = "Device"
  platforms = ["macOS"]
}

resource "zentral_mdm_enterprise_app" "test" {
  artifact_id    = zentral_mdm_artifact.test.id
  package_uri    = "s3://zentral-pro-services-artifacts-eu-central-1/terraform-provider-zentral/testdata/Outset-4.0.21900.pkg"
  package_sha256 = "6e364e261e15ffdae89e6dfbd1b9e583319cf897d2dfe624521f5d0ef43b0be7"
  version        = 1
  macos          = true
}
`, name)
}

func testAccMDMEnterpriseAppResourceConfigFull(name string) string {
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
  package_uri       = "s3://zentral-pro-services-artifacts-eu-central-1/terraform-provider-zentral/testdata/Outset-4.1.0.21912.pkg"
  package_sha256    = "d26590460ccee1895c999d3275475a1a9006366500cea6b4500b731e39d3cc2a"
  version           = 2
  macos             = true
  macos_min_version = "13.3.1"
  shard_modulo      = 5
  default_shard     = 1
  excluded_tag_ids  = [zentral_tag.excluded.id]
  tag_shards        = [{tag_id = zentral_tag.shard.id, shard=4}]
}
`, name)
}
