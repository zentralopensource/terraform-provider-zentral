package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMStoreAppResource(t *testing.T) {
	name := acctest.RandString(12)
	aResourceName := "zentral_mdm_artifact.test"
	etResourceName := "zentral_tag.excluded"
	laResourceName := "data.zentral_mdm_location_asset.test"
	resourceName := "zentral_mdm_store_app.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMStoreAppResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "artifact_id", aResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "location_asset_id", laResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "associated_domains.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "associated_domains_enable_direct_downloads", "false"),
					resource.TestCheckNoResourceAttr(
						resourceName, "content_filter_uuid"),
					resource.TestCheckNoResourceAttr(
						resourceName, "dns_proxy_uuid"),
					resource.TestCheckNoResourceAttr(
						resourceName, "vpn_uuid"),
					resource.TestCheckResourceAttr(
						resourceName, "prevent_backup", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "removable", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "remove_on_unenroll", "true"),
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
				Config: testAccMDMStoreAppResourceConfigFull(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "artifact_id", aResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "location_asset_id", laResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "associated_domains.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "associated_domains.0", "www.example.com"),
					resource.TestCheckResourceAttr(
						resourceName, "associated_domains_enable_direct_downloads", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "content_filter_uuid", "08198963-bd7b-4a40-80b5-1ebf8ae25c83"),
					resource.TestCheckResourceAttr(
						resourceName, "dns_proxy_uuid", "e1e8a2e5-2046-4eaf-acf8-4c878a372f52"),
					resource.TestCheckResourceAttr(
						resourceName, "vpn_uuid", "6e5395ff-e58d-441e-862c-826e6ad1dc33"),
					resource.TestCheckResourceAttr(
						resourceName, "prevent_backup", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "removable", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "remove_on_unenroll", "false"),
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

// TODO: hard coded values of a provisioned location asset
// on the server used for the integration tests
func testAccMDMStoreAppResourceConfigBare(name string) string {
	return fmt.Sprintf(`
data "zentral_mdm_location" "test" {
  name = "Terraform Provider CI/CD"
}

data "zentral_mdm_location_asset" "test" {
  location_id   = data.zentral_mdm_location.test.id
  adam_id       = "803453959"
  pricing_param = "STDQ"
}

resource "zentral_mdm_artifact" "test" {
  name      = %[1]q
  type      = "Store App"
  channel   = "Device"
  platforms = ["macOS"]
}

resource "zentral_mdm_store_app" "test" {
  artifact_id       = zentral_mdm_artifact.test.id
  location_asset_id = data.zentral_mdm_location_asset.test.id 
  version           = 1
  macos             = true
}
`, name)
}

func testAccMDMStoreAppResourceConfigFull(name string) string {
	return fmt.Sprintf(`
data "zentral_mdm_location" "test" {
  name = "Terraform Provider CI/CD"
}

data "zentral_mdm_location_asset" "test" {
  location_id   = data.zentral_mdm_location.test.id
  adam_id       = "803453959"
  pricing_param = "STDQ"
}

resource "zentral_mdm_artifact" "test" {
  name      = %[1]q
  type      = "Store App"
  channel   = "Device"
  platforms = ["macOS"]
}

resource "zentral_tag" "excluded" {
  name = "%[1]s excluded"
}

resource "zentral_tag" "shard" {
  name = "%[1]s shard"
}

resource "zentral_mdm_store_app" "test" {
  artifact_id                                = zentral_mdm_artifact.test.id
  location_asset_id                          = data.zentral_mdm_location_asset.test.id
  associated_domains                         = ["www.example.com"]
  associated_domains_enable_direct_downloads = true
  content_filter_uuid                        = "08198963-bd7b-4a40-80b5-1ebf8ae25c83"
  dns_proxy_uuid                             = "e1e8a2e5-2046-4eaf-acf8-4c878a372f52"
  vpn_uuid                                   = "6e5395ff-e58d-441e-862c-826e6ad1dc33"
  prevent_backup                             = true
  removable                                  = true
  remove_on_unenroll                         = false
  version                                    = 2
  macos                                      = true
  macos_min_version                          = "13.3.1"
  shard_modulo                               = 5
  default_shard                              = 1
  excluded_tag_ids                           = [zentral_tag.excluded.id]
  tag_shards                                 = [{ tag_id = zentral_tag.shard.id, shard = 4 }]
}
`, name)
}
