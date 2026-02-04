resource "zentral_mdm_artifact" "pp-internal-app" {
  name      = "Provisioning Profile for internal app"
  type      = "Provisioning Profile"
  channel   = "Device"
  platforms = ["macOS"]
}

resource "zentral_mdm_provisioning_profile" "pp-internal-app-v1" {
  artifact_id = zentral_mdm_artifact.pp-internal-app.id
  source      = filebase64("${path.module}/provisionprofiles/internal_app.provisionprofile")
  version     = 1
  macos       = true
}
