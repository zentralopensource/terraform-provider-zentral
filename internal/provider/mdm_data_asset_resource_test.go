package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMDataAssetResource(t *testing.T) {
	name := acctest.RandString(12)
	aResourceName := "zentral_mdm_artifact.test"
	etResourceName := "zentral_tag.excluded"
	resourceName := "zentral_mdm_data_asset.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMDataAssetResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "artifact_id", aResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "type", "ZIP"),
					resource.TestCheckResourceAttr(
						resourceName, "file_uri", "s3://zentral-pro-services-artifacts-eu-central-1/terraform-provider-zentral/testdata/pam.d.v1.zip"),
					resource.TestCheckResourceAttr(
						resourceName, "file_sha256", "d2caf672d62f52467db52f82bb8103c1c8032a0a7a7dc1e0307a5bfc31afec30"),
					resource.TestCheckResourceAttr(
						resourceName, "filename", "pam.d.v1.zip"),
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
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"file_uri"},
			},
			// Update and Read
			{
				Config: testAccMDMDataAssetResourceConfigFull(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "artifact_id", aResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "type", "ZIP"),
					resource.TestCheckResourceAttr(
						resourceName, "file_uri", "s3://zentral-pro-services-artifacts-eu-central-1/terraform-provider-zentral/testdata/pam.d.v2.zip"),
					resource.TestCheckResourceAttr(
						resourceName, "file_sha256", "388e1a3ac8dcf2f8fa0bf269c4004c59221312d65d7e42ef93150e278542b6dc"),
					resource.TestCheckResourceAttr(
						resourceName, "filename", "pam.d.v2.zip"),
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
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"file_uri"},
			},
		},
	})
}

func testAccMDMDataAssetResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_artifact" "test" {
  name      = %[1]q
  type      = "Data Asset"
  channel   = "Device"
  platforms = ["macOS"]
}

resource "zentral_mdm_data_asset" "test" {
  artifact_id    = zentral_mdm_artifact.test.id
  type           = "ZIP"
  file_uri       = "s3://zentral-pro-services-artifacts-eu-central-1/terraform-provider-zentral/testdata/pam.d.v1.zip"
  file_sha256    = "d2caf672d62f52467db52f82bb8103c1c8032a0a7a7dc1e0307a5bfc31afec30"
  version        = 1
  macos          = true
}
`, name)
}

func testAccMDMDataAssetResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_artifact" "test" {
  name      = %[1]q
  type      = "Data Asset"
  channel   = "Device"
  platforms = ["macOS"]
}

resource "zentral_tag" "excluded" {
  name = "%[1]s excluded"
}

resource "zentral_tag" "shard" {
  name = "%[1]s shard"
}

resource "zentral_mdm_data_asset" "test" {
  artifact_id       = zentral_mdm_artifact.test.id
  type              = "ZIP"
  file_uri          = "s3://zentral-pro-services-artifacts-eu-central-1/terraform-provider-zentral/testdata/pam.d.v2.zip"
  file_sha256       = "388e1a3ac8dcf2f8fa0bf269c4004c59221312d65d7e42ef93150e278542b6dc"
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
