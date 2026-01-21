package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMDEPEnrollmentCustomViewResource(t *testing.T) {
	name := acctest.RandString(12)
	resourceName := "zentral_mdm_dep_enrollment_custom_view.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMDEPEnrollmentCustomViewResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "dep_enrollment_id",
						"zentral_mdm_dep_enrollment.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						resourceName, "custom_view_id",
						"zentral_mdm_enrollment_custom_view.test", "id",
					),
					resource.TestCheckResourceAttr(resourceName, "weight", "0"),
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
				Config: testAccMDMDEPEnrollmentCustomViewResourceConfigFull(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "dep_enrollment_id",
						"zentral_mdm_dep_enrollment.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						resourceName, "custom_view_id",
						"zentral_mdm_enrollment_custom_view.test", "id",
					),
					resource.TestCheckResourceAttr(resourceName, "weight", "10"),
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

func testAccMDMDEPEnrollmentCustomViewResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_mdm_scep_issuer" "test" {
  name = %[1]q
  url = "https://www.example.com/scep"
  backend = "STATIC_CHALLENGE"
  static_challenge = {
    challenge = "yolo"
  }
}

# provisioned resource on the integration server
data "zentral_mdm_push_certificate" "test" {
  name = "TF provider GitHub"
}

# provisioned resource on the integration server
data "zentral_mdm_dep_virtual_server" "test" {
  name = "TF provider GitHub"
}

resource "zentral_mdm_dep_enrollment" "test" {
  display_name = %[1]q

  push_certificate_id = data.zentral_mdm_push_certificate.test.id

  scep_issuer_id = zentral_mdm_scep_issuer.test.id

  enrollment = {
    meta_business_unit_id = zentral_meta_business_unit.test.id
  }

  profile = {
    virtual_server_id = data.zentral_mdm_dep_virtual_server.test.id
    name              = %[1]q
  }
}

resource "zentral_mdm_enrollment_custom_view" "test" {
  name = %[1]q
  html = "<html/>"
}

resource "zentral_mdm_dep_enrollment_custom_view" "test" {
  dep_enrollment_id = zentral_mdm_dep_enrollment.test.id
  custom_view_id    = zentral_mdm_enrollment_custom_view.test.id
}
`, name)
}

func testAccMDMDEPEnrollmentCustomViewResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_mdm_scep_issuer" "test" {
  name = %[1]q
  url = "https://www.example.com/scep"
  backend = "STATIC_CHALLENGE"
  static_challenge = {
    challenge = "yolo"
  }
}

# provisioned resource on the integration server
data "zentral_mdm_push_certificate" "test" {
  name = "TF provider GitHub"
}

# provisioned resource on the integration server
data "zentral_mdm_dep_virtual_server" "test" {
  name = "TF provider GitHub"
}

resource "zentral_mdm_dep_enrollment" "test" {
  display_name = %[1]q

  push_certificate_id = data.zentral_mdm_push_certificate.test.id

  scep_issuer_id = zentral_mdm_scep_issuer.test.id

  enrollment = {
    meta_business_unit_id = zentral_meta_business_unit.test.id
  }

  profile = {
    virtual_server_id = data.zentral_mdm_dep_virtual_server.test.id
    name              = %[1]q
  }
}

resource "zentral_mdm_enrollment_custom_view" "test" {
  name = %[1]q
  html = "<html/>"
}

resource "zentral_mdm_dep_enrollment_custom_view" "test" {
  dep_enrollment_id = zentral_mdm_dep_enrollment.test.id
  custom_view_id    = zentral_mdm_enrollment_custom_view.test.id
  weight            = 10
}
`, name)
}
