package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMOTAEnrollmentResource(t *testing.T) {
	name := acctest.RandString(12)
	tagName := acctest.RandString(12)
	resourceName := "zentral_mdm_ota_enrollment.test"
	bpResourceName := "zentral_mdm_blueprint.test"
	pcResourceName := "data.zentral_mdm_push_certificate.test"
	rResourceName := "data.zentral_realm.test"
	scResouceName := "data.zentral_mdm_scep_config.test"
	mbuResourceName := "zentral_meta_business_unit.test"
	tagResourceName := "zentral_tag.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMOTAEnrollmentResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName, "display_name", "Zentral MDM"),
					resource.TestCheckNoResourceAttr(
						resourceName, "blueprint_id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "push_certificate_id", pcResourceName, "id"),
					resource.TestCheckNoResourceAttr(
						resourceName, "realm_uuid"),
					resource.TestCheckResourceAttrPair(
						resourceName, "scep_config_id", scResouceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "scep_verification", "false"),
					resource.TestCheckResourceAttrPair(
						resourceName, "meta_business_unit_id", mbuResourceName, "id"),
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
				Config: testAccMDMOTAEnrollmentResourceConfigFull(name, tagName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName, "display_name", name),
					resource.TestCheckResourceAttrPair(
						resourceName, "blueprint_id", bpResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "push_certificate_id", pcResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "realm_uuid", rResourceName, "uuid"),
					resource.TestCheckResourceAttrPair(
						resourceName, "scep_config_id", scResouceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "scep_verification", "true"),
					resource.TestCheckResourceAttrPair(
						resourceName, "meta_business_unit_id", mbuResourceName, "id"),
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

func testAccMDMOTAEnrollmentResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

# provisioned resource on the integration server
data "zentral_mdm_push_certificate" "test" {
  name = "TF provider GitHub"
}

# provisioned resource on the integration server
data "zentral_mdm_scep_config" "test" {
  name = "TF provider GitHub"
}

resource "zentral_mdm_ota_enrollment" "test" {
  name                  = %[1]q
  push_certificate_id   = data.zentral_mdm_push_certificate.test.id
  scep_config_id        = data.zentral_mdm_scep_config.test.id
  meta_business_unit_id = zentral_meta_business_unit.test.id
}
`, name)
}

func testAccMDMOTAEnrollmentResourceConfigFull(name string, tagName string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_mdm_blueprint" "test" {
  name = %[1]q
}

# provisioned resource on the integration server
data "zentral_mdm_push_certificate" "test" {
  name = "TF provider GitHub"
}

# provisioned resource on the integration server
data "zentral_realm" "test" {
  name = "TF provider GitHub"
}

# provisioned resource on the integration server
data "zentral_mdm_scep_config" "test" {
  name = "TF provider GitHub"
}

resource "zentral_taxonomy" "test" {
  name = %[1]q
}

resource "zentral_tag" "test" {
  taxonomy_id = zentral_taxonomy.test.id
  name        = %[2]q
}

resource "zentral_mdm_ota_enrollment" "test" {
  name                  = %[1]q
  display_name          = %[1]q
  blueprint_id          = zentral_mdm_blueprint.test.id
  push_certificate_id   = data.zentral_mdm_push_certificate.test.id
  realm_uuid            = data.zentral_realm.test.uuid
  scep_config_id        = data.zentral_mdm_scep_config.test.id
  scep_verification     = true
  meta_business_unit_id = zentral_meta_business_unit.test.id
  tag_ids               = [zentral_tag.test.id]
  serial_numbers        = ["un", "deux"]
  udids                 = ["trois", "quatre"]
  quota                 = 5
}
`, name, tagName)
}
