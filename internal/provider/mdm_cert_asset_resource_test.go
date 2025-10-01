package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMCertAssetResource(t *testing.T) {
	name := acctest.RandString(12)
	aResourceName := "zentral_mdm_artifact.test"
	aiResourceName := "zentral_mdm_acme_issuer.test"
	siResourceName := "zentral_mdm_scep_issuer.test"
	etResourceName := "zentral_tag.excluded"
	resourceName := "zentral_mdm_cert_asset.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMCertAssetResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "artifact_id", aResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "acme_issuer_id", aiResourceName, "id"),
					resource.TestCheckNoResourceAttr(
						resourceName, "scep_issuer_id"),
					resource.TestCheckResourceAttr(
						resourceName, "subject.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName, "subject.*", map[string]string{
							"type":  "CN",
							"value": "$ENROLLED_DEVICE.SERIAL_NUMBER",
						}),
					resource.TestCheckNoResourceAttr(
						resourceName, "subject_alt_name.dns_name"),
					resource.TestCheckNoResourceAttr(
						resourceName, "subject_alt_name.nt_principal_name"),
					resource.TestCheckNoResourceAttr(
						resourceName, "subject_alt_name.rfc822_name"),
					resource.TestCheckNoResourceAttr(
						resourceName, "subject_alt_name.uri"),
					resource.TestCheckResourceAttr(
						resourceName, "accessible", "Default"),
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
				Config: testAccMDMCertAssetResourceConfigFull(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "artifact_id", aResourceName, "id"),
					resource.TestCheckNoResourceAttr(
						resourceName, "acme_issuer_id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "scep_issuer_id", siResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "subject.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName, "subject.*", map[string]string{
							"type":  "CN",
							"value": "$ENROLLED_DEVICE.SERIAL_NUMBER",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName, "subject.*", map[string]string{
							"type":  "O",
							"value": "Zentral",
						}),
					resource.TestCheckNoResourceAttr(
						resourceName, "subject_alt_name.dns_name"),
					resource.TestCheckNoResourceAttr(
						resourceName, "subject_alt_name.nt_principal_name"),
					resource.TestCheckResourceAttr(
						resourceName, "subject_alt_name.rfc822_name", "support@zentral.com"),
					resource.TestCheckNoResourceAttr(
						resourceName, "subject_alt_name.uri"),
					resource.TestCheckResourceAttr(
						resourceName, "accessible", "AfterFirstUnlock"),
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
						resourceName, "macos_min_version", "15.7.1"),
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

func testAccMDMCertAssetResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource  "zentral_mdm_acme_issuer" "test" {
  name               = %[1]q
  directory_url      = "https://www.example.com/acme1"
  key_type           = "ECSECPrimeRandom"
  key_size           = 384
  usage_flags        = 1
  extended_key_usage = ["1.3.6.1.5.5.7.3.2"]
  backend          = "STATIC_CHALLENGE"
  static_challenge = {
    challenge = "Yolo"
  }
}

resource "zentral_mdm_artifact" "test" {
  name      = %[1]q
  type      = "Certificate Asset"
  channel   = "Device"
  platforms = ["macOS"]
}

resource "zentral_mdm_cert_asset" "test" {
  artifact_id    = zentral_mdm_artifact.test.id
  acme_issuer_id = zentral_mdm_acme_issuer.test.id
  subject = [
    {
      type  = "CN",
      value = "$ENROLLED_DEVICE.SERIAL_NUMBER"
    }
  ]
  version = 1
  macos   = true
}
`, name)
}

func testAccMDMCertAssetResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource  "zentral_mdm_acme_issuer" "test" {
  name               = %[1]q
  directory_url      = "https://www.example.com/acme1"
  key_type           = "ECSECPrimeRandom"
  key_size           = 384
  usage_flags        = 1
  extended_key_usage = ["1.3.6.1.5.5.7.3.2"]
  backend          = "STATIC_CHALLENGE"
  static_challenge = {
    challenge = "Yolo"
  }
}

resource  "zentral_mdm_scep_issuer" "test" {
  name             = %[1]q
  url              = "https://www.example.com/scep1"
  backend          = "STATIC_CHALLENGE"
  static_challenge = {
    challenge = "Yolo"
  }
}

resource "zentral_mdm_artifact" "test" {
  name      = %[1]q
  type      = "Certificate Asset"
  channel   = "Device"
  platforms = ["macOS"]
}

resource "zentral_tag" "excluded" {
  name = "%[1]s excluded"
}

resource "zentral_tag" "shard" {
  name = "%[1]s shard"
}

resource "zentral_mdm_cert_asset" "test" {
  artifact_id    = zentral_mdm_artifact.test.id
  scep_issuer_id = zentral_mdm_scep_issuer.test.id
  subject = [
    {
      type  = "CN"
      value = "$ENROLLED_DEVICE.SERIAL_NUMBER"
    },
    {
      type  = "O",
      value = "Zentral"
    },
  ]
  subject_alt_name = {
    rfc822_name = "support@zentral.com"
  }
  accessible        = "AfterFirstUnlock"
  version           = 2
  macos             = true
  macos_min_version = "15.7.1"
  shard_modulo      = 5
  default_shard     = 1
  excluded_tag_ids  = [zentral_tag.excluded.id]
  tag_shards = [
    {
      tag_id = zentral_tag.shard.id,
      shard  = 4
    }
  ]
}
`, name)
}
