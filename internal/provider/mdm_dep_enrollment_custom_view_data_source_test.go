package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMDEPEnrollmentCustomViewDataSource(t *testing.T) {
	name := acctest.RandString(12)
	linkResourceName := "zentral_mdm_dep_enrollment_custom_view.test"
	ds1ResourceName := "data.zentral_mdm_dep_enrollment_custom_view.by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMDEPEnrollmentCustomViewDataSourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by id
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", linkResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "dep_enrollment", "zentral_mdm_dep_enrollment.test", "id"),
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "custom_view", "zentral_mdm_enrollment_custom_view.test", "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "weight", "0"),
				),
			},
		},
	})
}

func testAccMDMDEPEnrollmentCustomViewDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_taxonomy" "test" {
  name = %[1]q
}

# provisioned resource on the integration server
data "zentral_mdm_push_certificate" "test" {
  name = "TF provider GitHub"
}

resource "zentral_mdm_scep_issuer" "test" {
  name = %[1]q
  url = "https://www.example.com/scep"
  backend = "STATIC_CHALLENGE"
  static_challenge = {
    challenge = "yolo"
  }
}

data "zentral_mdm_dep_virtual_server" "test" {
  name = "TF provider GitHub"
}

resource "zentral_mdm_dep_enrollment" "test" {
	name                  			= %[1]q
  	display_name          			= %[1]q

	meta_business_unit_id 			= zentral_meta_business_unit.test.id
	tag_ids               			= []
	serial_numbers        			= []
	udids                 			= []
	
	skip_setup_items 				= ["Accessibility"]

	push_certificate_id   			= data.zentral_mdm_push_certificate.test.id
  	scep_issuer_id        			= zentral_mdm_scep_issuer.test.id

	virtual_server_id				= data.zentral_mdm_dep_virtual_server.test.id
}

resource "zentral_mdm_enrollment_custom_view" "test" {
	name                  			= %[1]q
	description						= %[1]q
	html							= "<html/>"
}

resource "zentral_mdm_dep_enrollment_custom_view" "test" {
	dep_enrollment  = zentral_mdm_dep_enrollment.test.id
	custom_view 	= zentral_mdm_enrollment_custom_view.test.id
}

data "zentral_mdm_dep_enrollment_custom_view" "by_id" {
  id = zentral_mdm_dep_enrollment_custom_view.test.id
}

`, name)
}
