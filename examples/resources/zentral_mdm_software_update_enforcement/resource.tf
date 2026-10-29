# The macOS 26 devices get the latest macOS 26 update
resource "zentral_mdm_software_update_enforcement" "macos_26" {
  name                  = "macOS 26"
  platforms             = ["macOS"]
  min_device_os_version = "26"
  max_device_os_version = "27"
  max_os_version        = "27"
}

# The macOS 27 devices get the latest macOS 27 update
resource "zentral_mdm_software_update_enforcement" "macos_27" {
  name                  = "macOS 27"
  platforms             = ["macOS"]
  min_device_os_version = "27"
  max_device_os_version = "28"
  max_os_version        = "28"
}

# The tag that moves a device to macOS 27
resource "zentral_tag" "upgrade_27" {
  name = "upgrade-27"
}

# The tagged devices are upgraded to the latest macOS 27 update, whatever their OS version.
# The tags rank before the device OS version range, so this enforcement wins over "macOS 26"
# for a tagged macOS 26 device. After the upgrade, "macOS 27" takes over with the same target.
resource "zentral_mdm_software_update_enforcement" "upgrade_27" {
  name           = "Upgrade to macOS 27"
  platforms      = ["macOS"]
  tag_ids        = [zentral_tag.upgrade_27.id]
  max_os_version = "28"
}

# The three enforcements in one blueprint
resource "zentral_mdm_blueprint" "default" {
  name = "Default"
  software_update_enforcement_ids = [
    zentral_mdm_software_update_enforcement.macos_26.id,
    zentral_mdm_software_update_enforcement.macos_27.id,
    zentral_mdm_software_update_enforcement.upgrade_27.id,
  ]
}
