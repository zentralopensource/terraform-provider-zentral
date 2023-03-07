package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonolithSubManifestPkgInfoResource(t *testing.T) {
	name := acctest.RandString(12)
	resourceName := "zentral_monolith_sub_manifest_pkg_info.test"
	smResourceName := "zentral_monolith_sub_manifest.test"
	etResourceName := "zentral_tag.excluded"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMonolithSubManifestPkgInfoResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "sub_manifest_id", smResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "key", "optional_installs"),
					resource.TestCheckResourceAttr(
						resourceName, "pkg_info_name", "Firefox"),
					resource.TestCheckResourceAttr(
						resourceName, "featured_item", "false"),
					resource.TestCheckNoResourceAttr(
						resourceName, "condition_id"),
					resource.TestCheckResourceAttr(
						resourceName, "shard_modulo", "100"),
					resource.TestCheckResourceAttr(
						resourceName, "default_shard", "100"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_shards.#", "0"),
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
				Config: testAccMonolithSubManifestPkgInfoResourceConfigFull(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "sub_manifest_id", smResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "key", "optional_installs"),
					resource.TestCheckResourceAttr(
						resourceName, "pkg_info_name", "Firefox"),
					resource.TestCheckResourceAttr(
						resourceName, "featured_item", "true"),
					resource.TestCheckNoResourceAttr(
						resourceName, "condition_id"),
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

func testAccMonolithSubManifestPkgInfoResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_monolith_sub_manifest" "test" {
  name = %[1]q
}

resource "zentral_monolith_sub_manifest_pkg_info" "test" {
  sub_manifest_id = zentral_monolith_sub_manifest.test.id
  key             = "optional_installs"
  pkg_info_name   = "Firefox"
}
`, name)
}

func testAccMonolithSubManifestPkgInfoResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_monolith_sub_manifest" "test" {
  name = %[1]q
}

resource "zentral_tag" "excluded" {
  name = "%[1]s excluded"
}

resource "zentral_tag" "shard" {
  name = "%[1]s shard"
}

resource "zentral_monolith_sub_manifest_pkg_info" "test" {
  sub_manifest_id  = zentral_monolith_sub_manifest.test.id
  key              = "optional_installs"
  pkg_info_name    = "Firefox"
  featured_item    = true
  shard_modulo     = 5
  default_shard    = 1
  excluded_tag_ids = [zentral_tag.excluded.id]
  tag_shards = [
    {
      tag_id = zentral_tag.shard.id,
      shard  = 4
    }
  ]
}
`, name)
}
