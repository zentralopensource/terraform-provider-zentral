package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMDEPEnrollmentResource_basic(t *testing.T) {
	name := acctest.RandString(12)
	tagName := acctest.RandString(12)
	resourceName := "zentral_mdm_dep_enrollment.test"
	bpResourceName := "zentral_mdm_blueprint.test"
	pcResourceName := "data.zentral_mdm_push_certificate.test"
	rResourceName := "data.zentral_realm.test"
	siResouceName := "zentral_mdm_scep_issuer.test"
	mbuResourceName := "zentral_meta_business_unit.test"
	tagResourceName := "zentral_tag.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMDEPEnrollmentResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName, "display_name", name),
					resource.TestCheckResourceAttr(
						resourceName, "skip_setup_items.#", "1"),
					resource.TestCheckNoResourceAttr(
						resourceName, "blueprint_id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "push_certificate_id", pcResourceName, "id"),
					resource.TestCheckNoResourceAttr(
						resourceName, "realm_uuid"),
					resource.TestCheckNoResourceAttr(
						resourceName, "acme_issuer_id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "scep_issuer_id", siResouceName, "id"),
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
				Config: testAccMDMDEPEnrollmentResourceConfigFull(name, tagName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName, "display_name", name),

					resource.TestCheckResourceAttr(
						resourceName, "use_realm_user", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "realm_user_is_admin", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "hidden_admin", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "allow_pairing", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_advance_setup", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "await_device_configured", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "is_mandatory", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "is_mdm_removable", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "is_multi_user", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "is_supervised", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "include_tls_certificates", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "username_pattern", "$REALM_USER.EMAIL_PREFIX"),
					resource.TestCheckResourceAttr(
						resourceName, "admin_password_complexity", "3"),
					resource.TestCheckResourceAttr(
						resourceName, "admin_password_rotation_delay", "60"),
					resource.TestCheckResourceAttr(
						resourceName, "skip_setup_items.#", "3"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "skip_setup_items.*", "Accessibility"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "skip_setup_items.*", "ActionButton"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "skip_setup_items.*", "AgeBasedSafetySettings"),
					resource.TestCheckResourceAttr(
						resourceName, "support_email_address", "help@zentral.com"),
					resource.TestCheckResourceAttr(
						resourceName, "support_phone_number", "+1234567890"),
					resource.TestCheckResourceAttr(
						resourceName, "ios_max_version", "1.2.3"),
					resource.TestCheckResourceAttr(
						resourceName, "ios_min_version", "1.2.3"),
					resource.TestCheckResourceAttr(
						resourceName, "macos_max_version", "1.2.3"),
					resource.TestCheckResourceAttr(
						resourceName, "macos_min_version", "1.2.3"),
					resource.TestCheckResourceAttr(
						resourceName, "department", "department"),
					resource.TestCheckResourceAttr(
						resourceName, "org_magic", "abracadabra"),
					resource.TestCheckResourceAttr(
						resourceName, "language", "en"),
					resource.TestCheckResourceAttr(
						resourceName, "region", "DE"),
					resource.TestCheckResourceAttrPair(
						resourceName, "blueprint_id", bpResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "push_certificate_id", pcResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "realm_uuid", rResourceName, "uuid"),
					resource.TestCheckNoResourceAttr(
						resourceName, "acme_issuer_id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "scep_issuer_id", siResouceName, "id"),
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

func testAccMDMDEPEnrollmentResourceConfigBare(name string) string {
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
`, name)
}

func testAccMDMDEPEnrollmentResourceConfigFull(name string, tagName string) string {
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

resource "zentral_mdm_scep_issuer" "test" {
  name = %[1]q
  url = "https://www.example.com/scep"
  backend = "STATIC_CHALLENGE"
  static_challenge = {
    challenge = "yolo"
  }
}

resource "zentral_taxonomy" "test" {
  name = %[1]q
}

resource "zentral_tag" "test" {
  taxonomy_id = zentral_taxonomy.test.id
  name        = %[2]q
}

data "zentral_mdm_dep_virtual_server" "test" {
  name = "TF provider GitHub"
}

resource "zentral_mdm_dep_enrollment" "test" {
	name                  			= %[1]q
  	display_name          			= %[1]q

	meta_business_unit_id 			= zentral_meta_business_unit.test.id
	tag_ids               			= [zentral_tag.test.id]
	serial_numbers        			= ["un", "deux"]
	udids                 			= ["trois", "quatre"]
	quota                 			= 5

	use_realm_user			 		= true
	realm_user_is_admin		 		= false
	hidden_admin			 		= false
	allow_pairing			 		= false
	auto_advance_setup		 		= false
	await_device_configured	 		= true
	is_mandatory			 		= true
	is_mdm_removable		 		= true
	is_multi_user			 		= false
	is_supervised			 		= false
	include_tls_certificates 		= false

	username_pattern 				= "$REALM_USER.EMAIL_PREFIX"
	admin_password_complexity 		= 3
	admin_password_rotation_delay 	= 60

	skip_setup_items 				= ["Accessibility", "ActionButton", "AgeBasedSafetySettings"]

	support_email_address 			= "help@zentral.com"
	support_phone_number			= "+1234567890"

	ios_max_version					= "1.2.3"
	ios_min_version					= "1.2.3"
	macos_max_version				= "1.2.3"
	macos_min_version				= "1.2.3"

	department 						= "department"
	org_magic 						= "abracadabra"
	language 						= "en"
	region 							= "DE"

	blueprint_id          			= zentral_mdm_blueprint.test.id
  	push_certificate_id   			= data.zentral_mdm_push_certificate.test.id
  	realm_uuid            			= data.zentral_realm.test.uuid
  	scep_issuer_id        			= zentral_mdm_scep_issuer.test.id

	virtual_server_id				= data.zentral_mdm_dep_virtual_server.test.id
}
`, name, tagName)
}
