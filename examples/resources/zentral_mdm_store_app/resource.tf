# First we fetch the ID of the apps & books location

data "zentral_mdm_location" "default" {
  name = "My Apps & Books Location"
}


# Then we fetch the ID of the location asset (the store app).

data "zentral_mdm_location_asset" "macos-slack" {
  location_id   = data.zentral_mdm_location.default.id
  adam_id       = "803453959"
  pricing_param = "STDQ"
}

# We add the artifact to manage the store app.

resource "zentral_mdm_artifact" "macos-slack" {
  name      = "Slack for macOS"
  type      = "Store App"
  channel   = "Device"
  platforms = ["macOS"]
}

# We add the store app, which is an artifact version

resource "zentral_mdm_store_app" "macos-slack-1" {
  artifact_id       = zentral_mdm_artifact.macos-slack.id
  location_asset_id = data.zentral_mdm_location_asset.macos-slack.id
  version           = 1
  macos             = true
}
