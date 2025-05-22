package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMDeclarationResource(t *testing.T) {
	name := acctest.RandString(12)
	identifer := acctest.RandString(12)
	serverToken := acctest.RandString(12)
	aResourceName := "zentral_mdm_artifact.test"
	etResourceName := "zentral_tag.excluded"
	resourceName := "zentral_mdm_declaration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMDeclarationResourceConfigBare(name, identifer, serverToken),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "artifact_id", aResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "source",
						fmt.Sprintf(
							`{"Identifier":%[1]q,"Payload":{"MinimumLength":10},"ServerToken":%[2]q,"Type":"com.apple.configuration.passcode.settings"}`,
							identifer, serverToken,
						),
					),
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
				Config: testAccMDMDeclarationResourceConfigFull(name, identifer, serverToken),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "artifact_id", aResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "source",
						fmt.Sprintf(`{"Identifier":%[1]q,"Payload":{"MaximumFailedAttempts":3,"MinimumLength":11},"ServerToken":%[2]q,"Type":"com.apple.configuration.passcode.settings"}`,
							identifer, serverToken,
						),
					),
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

func testAccMDMDeclarationResourceConfigBare(name string, identifer string, serverToken string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_artifact" "test" {
  name      = %[1]q
  type      = "Configuration"
  channel   = "Device"
  platforms = ["macOS"]
}

resource "zentral_mdm_declaration" "test" {
  artifact_id = zentral_mdm_artifact.test.id
  source = jsonencode({
    Type        = "com.apple.configuration.passcode.settings"
    Identifier  = %[2]q
    ServerToken = %[3]q
    Payload = {
      MinimumLength = 10
    }
  })
  version = 1
  macos   = true
}
`, name, identifer, serverToken)
}

func testAccMDMDeclarationResourceConfigFull(name string, identifer string, serverToken string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_artifact" "test" {
  name      = %[1]q
  type      = "Configuration"
  channel   = "Device"
  platforms = ["macOS"]
}

resource "zentral_tag" "excluded" {
  name = "%[1]s excluded"
}

resource "zentral_tag" "shard" {
  name = "%[1]s shard"
}

resource "zentral_mdm_declaration" "test" {
  artifact_id = zentral_mdm_artifact.test.id
  source = jsonencode({
    Type        = "com.apple.configuration.passcode.settings"
    Identifier  = %[2]q
    ServerToken = %[3]q
    Payload = {
      MinimumLength = 11
      MaximumFailedAttempts = 3
    }
  })
  version           = 2
  macos             = true
  macos_min_version = "13.3.1"
  shard_modulo      = 5
  default_shard     = 1
  excluded_tag_ids  = [zentral_tag.excluded.id]
  tag_shards        = [{ tag_id = zentral_tag.shard.id, shard = 4 }]
}
`, name, identifer, serverToken)
}
