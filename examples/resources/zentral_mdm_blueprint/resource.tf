# The default AxM apps & books location to use for the licenses
data "zentral_mdm_location" "default" {
  name = "ACME corp"
}

# A default FileVault configuration
resource "zentral_mdm_filevault_config" "default" {
  name                         = "Default"
  escrow_location_display_name = "ACME corp"
}

# A firmware password / recovery lock configuration
resource "zentral_mdm_recovery_password_config" "default" {
  name = "Default"
}

# Enforce OS updates until v27 (excluded) on some platforms
resource "zentral_mdm_software_update_enforcement" "default" {
  name           = "Default"
  platforms      = ["iOS", "iPadOS", "macOS"]
  max_os_version = "27"
}

# The MDM blueprint
resource "zentral_mdm_blueprint" "default" {
  name                            = "Default"
  inventory_interval              = 86400
  collect_apps                    = "MANAGED_ONLY"
  collect_certificates            = "ALL"
  collect_profiles                = "MANAGED_ONLY"
  legacy_profiles_via_ddm         = true
  default_location_id             = data.zentral_mdm_location.default.id
  filevault_config_id             = zentral_mdm_filevault_config.default.id
  recovery_password_config_id     = zentral_mdm_recovery_password_config.default.id
  software_update_enforcement_ids = [zentral_mdm_software_update_enforcement.default.id]
}
