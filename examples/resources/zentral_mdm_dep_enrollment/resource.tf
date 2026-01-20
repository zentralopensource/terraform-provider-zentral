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

// New tag that devices will get at enrollment
resource "zentral_tag" "new-devices" {
  name = "New device"
}

// Existing realm used to authenticate during ADE/DEP
data "zentral_realm" "google-workspace-ade" {
  name = "Google Workspace ADE authentication"
}

// New blueprint
resource "zentral_mdm_blueprint" "new" {
  name = "The new blueprint"
}

// Existing ACME issuer
data "zentral_mdm_acme_issuer" "default" {
  name = "Default ACME issuer"
}

// Full ADE/DEP enrollment
resource "zentral_mdm_dep_enrollment" "full" {
  display_name = "Zentral MDM"

  push_certificate_id = data.zentral_mdm_push_certificate.default.id
  blueprint_id        = zentral_mdm_blueprint.new.id

  // Certificate issuers
  acme_issuer_id = data.zentral_mdm_acme_issuer.default.id
  scep_issuer_id = data.zentral_mdm_scep_issuer.default.id

  // Enrollment
  enrollment = {
    meta_business_unit_id = zentral_meta_business_unit.default.id
    tag_ids               = [zentral_tag.new-device.id]
  }

  // Authentication
  authentication = {
    realm_uuid                       = data.zentral_realm.google-workspace-ade.uuid
    use_for_setup_assistant_user     = true
    setup_assistant_user_is_admin    = true
    setup_assistant_username_pattern = "$REALM_USER.EMAIL_PREFIX"
  }

  // Extra admin user
  extra_admin = {
    hidden                  = true
    full_name               = "Org Admin"
    short_name              = "_ladm"
    password_complexity     = 3
    password_rotation_delay = 120
  }

  // https://developer.apple.com/documentation/devicemanagement/profile/
  profile = {
    virtual_server_id = data.zentral_mdm_dep_virtual_server.this
    name              = "Zentral MDM full"

    allow_pairing           = false
    auto_advance_setup      = false
    await_device_configured = true // required when authentication → use_for_user_creation

    is_mandatory     = true
    is_mdm_removable = false
    is_multi_user    = false
    is_supervised    = true

    skip_setup_items = [
      "ActionButton",
      "Android",
      "Appearance",
      "Welcome",
      "Zoom"
    ]

    language = "en"
    region   = "DE"

    department = "IT"
    org_magic  = "This is a secret"

    support_email_address = "support@example.com"
    support_phone_number  = "+1-555-438-3732"
  }

  // OS version enforcement
  os_version_enforcement = {
    macos_min_version = "26.2"
    // and rolling min version until… based on info sent by device
    auto_macos_min_version_until = "27"
    ios_min_version              = "26.2"
    // and rolling min version until… based on info sent by device
    auto_ios_min_version_until = "27"
  }
}
