package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMDEPEnrollmentResource(t *testing.T) {
	name := acctest.RandString(12)
	tagName := acctest.RandString(12)
	resourceName := "zentral_mdm_dep_enrollment.test"
	bpResourceName := "zentral_mdm_blueprint.test"
	pcResourceName := "data.zentral_mdm_push_certificate.test"
	rResourceName := "data.zentral_realm.test"
	aiResourceName := "zentral_mdm_acme_issuer.test"
	siResouceName := "zentral_mdm_scep_issuer.test"
	vsResourceName := "data.zentral_mdm_dep_virtual_server.test"
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
						resourceName, "display_name", name),
					resource.TestCheckResourceAttrPair(
						resourceName, "push_certificate_id", pcResourceName, "id"),
					resource.TestCheckNoResourceAttr(
						resourceName, "blueprint_id"),
					resource.TestCheckNoResourceAttr(
						resourceName, "acme_issuer_id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "scep_issuer_id", siResouceName, "id"),
					// EnrollmentSecret
					resource.TestCheckResourceAttrPair(
						resourceName, "enrollment.meta_business_unit_id", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "enrollment.tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "enrollment.serial_numbers.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "enrollment.uuids.#", "0"),
					resource.TestCheckNoResourceAttr(
						resourceName, "enrollment.quota"),
					// Authentication
					resource.TestCheckNoResourceAttr(
						resourceName, "authentication.realm_uuid"),
					resource.TestCheckResourceAttr(
						resourceName, "authentication.use_for_setup_assistant_user", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "authentication.setup_assistant_user_is_admin", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "authentication.setup_assistant_username_pattern", ""),
					// Extra admin
					resource.TestCheckResourceAttr(
						resourceName, "extra_admin.hidden", "false"),
					resource.TestCheckNoResourceAttr(
						resourceName, "extra_admin.full_name"),
					resource.TestCheckNoResourceAttr(
						resourceName, "extra_admin.short_name"),
					resource.TestCheckResourceAttr(
						resourceName, "extra_admin.password_complexity", "3"),
					resource.TestCheckResourceAttr(
						resourceName, "extra_admin.password_rotation_delay", "60"),
					// Profile
					resource.TestCheckResourceAttrPair(
						resourceName, "profile.virtual_server_id", vsResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.name", name),
					resource.TestCheckResourceAttr(
						resourceName, "profile.allow_pairing", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.auto_advance_setup", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.await_device_configured", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.is_mandatory", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.is_mdm_removable", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.is_multi_user", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.is_supervised", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.include_anchor_certs", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.skip_setup_items.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.language", ""),
					resource.TestCheckResourceAttr(
						resourceName, "profile.region", ""),
					resource.TestCheckResourceAttr(
						resourceName, "profile.department", ""),
					resource.TestCheckResourceAttr(
						resourceName, "profile.org_magic", ""),
					resource.TestCheckResourceAttr(
						resourceName, "profile.support_email_address", ""),
					resource.TestCheckResourceAttr(
						resourceName, "profile.support_phone_number", ""),
					// OS version enforcement
					resource.TestCheckResourceAttr(
						resourceName, "os_version_enforcement.ios_min_version", ""),
					resource.TestCheckResourceAttr(
						resourceName, "os_version_enforcement.auto_ios_min_version_until", ""),
					resource.TestCheckResourceAttr(
						resourceName, "os_version_enforcement.macos_min_version", ""),
					resource.TestCheckResourceAttr(
						resourceName, "os_version_enforcement.auto_macos_min_version_until", ""),
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
						resourceName, "display_name", name),
					resource.TestCheckResourceAttrPair(
						resourceName, "push_certificate_id", pcResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "blueprint_id", bpResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "acme_issuer_id", aiResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "scep_issuer_id", siResouceName, "id"),
					// EnrollmentSecret
					resource.TestCheckResourceAttrPair(
						resourceName, "enrollment.meta_business_unit_id", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "enrollment.tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "enrollment.tag_ids.*", tagResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "enrollment.serial_numbers.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "enrollment.serial_numbers.*", "un"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "enrollment.serial_numbers.*", "deux"),
					resource.TestCheckResourceAttr(
						resourceName, "enrollment.udids.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "enrollment.udids.*", "trois"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "enrollment.udids.*", "quatre"),
					resource.TestCheckResourceAttr(
						resourceName, "enrollment.quota", "5"),
					// Authentication
					resource.TestCheckResourceAttrPair(
						resourceName, "authentication.realm_uuid", rResourceName, "uuid"),
					resource.TestCheckResourceAttr(
						resourceName, "authentication.use_for_setup_assistant_user", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "authentication.setup_assistant_user_is_admin", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "authentication.setup_assistant_username_pattern", "$REALM_USER.EMAIL_PREFIX"),
					// Extra admin
					resource.TestCheckResourceAttr(
						resourceName, "extra_admin.hidden", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "extra_admin.full_name", "Our Admin"),
					resource.TestCheckResourceAttr(
						resourceName, "extra_admin.short_name", "ouradmin"),
					resource.TestCheckResourceAttr(
						resourceName, "extra_admin.password_complexity", "2"),
					resource.TestCheckResourceAttr(
						resourceName, "extra_admin.password_rotation_delay", "57"),
					// Profile
					resource.TestCheckResourceAttrPair(
						resourceName, "profile.virtual_server_id", vsResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.name", name),
					resource.TestCheckResourceAttr(
						resourceName, "profile.allow_pairing", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.auto_advance_setup", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.await_device_configured", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.is_mandatory", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.is_mdm_removable", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.is_multi_user", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.is_supervised", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.include_anchor_certs", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.skip_setup_items.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "profile.skip_setup_items.*", "Accessibility"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "profile.skip_setup_items.*", "ActionButton"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.language", "en"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.region", "DE"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.department", "IT"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.org_magic", "abracadabra"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.support_email_address", "support@example.com"),
					resource.TestCheckResourceAttr(
						resourceName, "profile.support_phone_number", "+1-555-438-3732"),
					// OS version enforcement
					resource.TestCheckResourceAttr(
						resourceName, "os_version_enforcement.ios_min_version", "26.2"),
					resource.TestCheckResourceAttr(
						resourceName, "os_version_enforcement.auto_ios_min_version_until", "27"),
					resource.TestCheckResourceAttr(
						resourceName, "os_version_enforcement.macos_min_version", "26.3"),
					resource.TestCheckResourceAttr(
						resourceName, "os_version_enforcement.auto_macos_min_version_until", "28"),
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

resource "zentral_mdm_scep_issuer" "test" {
  name    = %[1]q
  url     = "https://www.example.com/scep"
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

resource  "zentral_mdm_acme_issuer" "test" {
  name               = %[1]q
  description        = "Description"
  directory_url      = "https://www.example.com/acme"
  key_type           = "ECSECPrimeRandom"
  key_size           = 256
  usage_flags        = 1
  extended_key_usage = ["1.3.6.1.5.5.7.3.2"]
  hardware_bound     = true
  attest             = true
  backend            = "MICROSOFT_CA"
  microsoft_ca       = {
    "url"      = "https://www.example.com/ndes"
    "username" = "yolo"
    "password" = "fomo"
  }
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

# provisioned resource on the integration server
data "zentral_mdm_push_certificate" "test" {
  name = "TF provider GitHub"
}

# provisioned resource on the integration server
data "zentral_realm" "test" {
  name = "TF provider GitHub"
}

# provisioned resource on the integration server
data "zentral_mdm_dep_virtual_server" "test" {
  name = "TF provider GitHub"
}

resource "zentral_mdm_dep_enrollment" "test" {
  display_name = %[1]q

  push_certificate_id = data.zentral_mdm_push_certificate.test.id
  blueprint_id        = zentral_mdm_blueprint.test.id

  acme_issuer_id = zentral_mdm_acme_issuer.test.id
  scep_issuer_id = zentral_mdm_scep_issuer.test.id

  enrollment = {
    meta_business_unit_id = zentral_meta_business_unit.test.id
    tag_ids               = [zentral_tag.test.id]
    serial_numbers        = ["un", "deux"]
    udids                 = ["trois", "quatre"]
    quota                 = 5
  }

  authentication = {
    realm_uuid                       = data.zentral_realm.test.uuid
    use_for_setup_assistant_user     = true
    setup_assistant_user_is_admin    = true
    setup_assistant_username_pattern = "$REALM_USER.EMAIL_PREFIX"
  }

  extra_admin = {
    hidden                  = true
    full_name               = "Our Admin"
    short_name              = "ouradmin"
    password_complexity     = 2
    password_rotation_delay = 57
  }

  profile = {
    virtual_server_id = data.zentral_mdm_dep_virtual_server.test.id
    name              = %[1]q

    allow_pairing           = false
    auto_advance_setup      = true
    await_device_configured = true

    is_mandatory     = true
    is_mdm_removable = true
    is_multi_user    = true
    is_supervised    = true

    include_anchor_certs = true

    skip_setup_items = [
      "Accessibility",
      "ActionButton",
    ]

    language = "en"
    region   = "DE"

    department = "IT"
    org_magic  = "abracadabra"

    support_email_address = "support@example.com"
    support_phone_number  = "+1-555-438-3732"
  }

  os_version_enforcement = {
    ios_min_version              = "26.2"
    auto_ios_min_version_until   = "27"
    macos_min_version            = "26.3"
    auto_macos_min_version_until = "28"
  }
}
`, name, tagName)
}
