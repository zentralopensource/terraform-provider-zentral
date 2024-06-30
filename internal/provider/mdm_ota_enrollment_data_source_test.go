package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMOTAEnrollmentDataSource(t *testing.T) {
	bidName := acctest.RandString(12)
	bnName := acctest.RandString(12)
	tagName := acctest.RandString(12)
	bidResourceName := "zentral_mdm_ota_enrollment.by_id"
	bnResourceName := "zentral_mdm_ota_enrollment.by_name"
	ds1ResourceName := "data.zentral_mdm_ota_enrollment.by_id"
	ds2ResourceName := "data.zentral_mdm_ota_enrollment.by_name"
	bpResourceName := "zentral_mdm_blueprint.test"
	pcResourceName := "data.zentral_mdm_push_certificate.test"
	scResouceName := "data.zentral_mdm_scep_config.test"
	mbuResourceName := "zentral_meta_business_unit.test"
	tagResourceName := "zentral_tag.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMOTAEnrollmentDataSourceConfig(bidName, bnName, tagName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by id
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", bidResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", bidName),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "display_name", "Zentral MDM"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "blueprint_id"),
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "push_certificate_id", pcResourceName, "id"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "realm_id"),
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "scep_config_id", scResouceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "scep_verification", "false"),
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "meta_business_unit_id", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "serial_numbers.#", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "udids.#", "0"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "quota"),
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", bnResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", bnName),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "display_name", bnName),
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "blueprint_id", bpResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "push_certificate_id", pcResourceName, "id"),
					resource.TestCheckNoResourceAttr(
						ds2ResourceName, "realm_id"),
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "scep_config_id", scResouceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "scep_verification", "true"),
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "meta_business_unit_id", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						ds2ResourceName, "tag_ids.*", tagResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "serial_numbers.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "serial_numbers.*", "un"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "serial_numbers.*", "deux"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "udids.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "udids.*", "trois"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "udids.*", "quatre"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "quota", "5"),
				),
			},
		},
	})
}

func testAccMDMOTAEnrollmentDataSourceConfig(bidName string, bnName string, tagName string) string {
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
data "zentral_mdm_scep_config" "test" {
  name = "TF provider GitHub"
}

resource "zentral_taxonomy" "test" {
  name = %[1]q
}

resource "zentral_tag" "test" {
  taxonomy_id = zentral_taxonomy.test.id
  name        = %[3]q
}

resource "zentral_mdm_ota_enrollment" "by_id" {
  name                  = %[1]q
  push_certificate_id   = data.zentral_mdm_push_certificate.test.id
  scep_config_id        = data.zentral_mdm_scep_config.test.id
  meta_business_unit_id = zentral_meta_business_unit.test.id
}

resource "zentral_mdm_ota_enrollment" "by_name" {
  name                  = %[2]q
  display_name          = %[2]q
  blueprint_id          = zentral_mdm_blueprint.test.id
  push_certificate_id   = data.zentral_mdm_push_certificate.test.id
  # TODO Add realm_id when implemented
  scep_config_id        = data.zentral_mdm_scep_config.test.id
  scep_verification     = true
  meta_business_unit_id = zentral_meta_business_unit.test.id
  tag_ids               = [zentral_tag.test.id]
  serial_numbers        = ["un", "deux"]
  udids                 = ["trois", "quatre"]
  quota                 = 5
}

data "zentral_mdm_ota_enrollment" "by_id" {
  id = zentral_mdm_ota_enrollment.by_id.id
}

data "zentral_mdm_ota_enrollment" "by_name" {
  name = zentral_mdm_ota_enrollment.by_name.name
}
`, bidName, bnName, tagName)
}
