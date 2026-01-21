// The meta business unit for the enrollment
resource "zentral_meta_business_unit" "default" {
  name                   = "Default"
  api_enrollment_enabled = true
}

// Existing MDM server configured in ABM/ASM
data "zentral_mdm_dep_virtual_server" "this" {
  name = "tfdocs.zentral.com"
}

// Existing MDM push certificate
data "zentral_mdm_push_certificate" "default" {
  name = "Default Push certificate"
}

// Existing SCEP issuer
data "zentral_mdm_scep_issuer" "default" {
  name = "Default SCEP issuer"
}

// Simple ADE/DEP enrollment
resource "zentral_mdm_dep_enrollment" "simple" {
  display_name = "Zentral MDM"

  push_certificate_id = data.zentral_mdm_push_certificate.default.id

  // Certificate issuer
  scep_issuer_id = data.zentral_mdm_scep_issuer.default.id

  // Enrollment
  enrollment = {
    meta_business_unit_id = zentral_meta_business_unit.default.id
  }

  // https://developer.apple.com/documentation/devicemanagement/profile/
  profile = {
    virtual_server_id = data.zentral_mdm_dep_virtual_server.this
    name              = "Zentral MDM simple"
  }
}

// MDM enrollment custom view
// Can be shared by many MDM enrollments
resource "zentral_mdm_enrollment_custom_view" "welcome" {
  name                    = "Welcome screen"
  description             = "Welcome screen displayed pre realm authentication"
  html                    = file("${path.module}/mdm_custom_views/welcome.html")
  requires_authentication = false
}

// Link between the custom view and the simple DEP enrollment
resource "zentral_mdm_dep_enrollment_custom_view" "test" {
  dep_enrollment_id = zentral_mdm_dep_enrollment.simple.id
  custom_view_id    = zentral_mdm_enrollment_custom_view.welcome.id
  weight            = 1
}
