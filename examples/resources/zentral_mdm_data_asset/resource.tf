# A data asset carries a file the devices download over a
# `com.apple.asset.data` declaration: a PLIST, or a ZIP archive.
resource "zentral_mdm_artifact" "app-config" {
  name      = "App Config"
  type      = "Data Asset"
  channel   = "Device"
  platforms = ["iOS", "iPadOS"]
}

# Give the content in the request with `source`, base 64 encoded. Zentral
# computes the sha256, and the data asset has no filename.
resource "zentral_mdm_data_asset" "app-config-v1" {
  artifact_id = zentral_mdm_artifact.app-config.id
  type        = "PLIST"
  source      = filebase64("${path.module}/plists/app-config.v1.plist")
  ios         = true
  ipados      = true
  version     = 1
}

# Or point at an object in an S3 bucket the server can read. `file_sha256` is
# required with `file_uri`, and the server verifies the file against it.
resource "zentral_mdm_data_asset" "app-config-legacy" {
  artifact_id = zentral_mdm_artifact.app-config.id
  type        = "ZIP"
  file_uri    = "s3://acme-mdm-assets/app-config.v0.zip"
  file_sha256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  ios         = true
  ipados      = true
  version     = 2
}

# The declaration that references the data asset by artifact ID.
resource "zentral_mdm_artifact" "app" {
  name      = "App"
  type      = "Configuration"
  channel   = "Device"
  platforms = ["iOS", "iPadOS"]
}

resource "zentral_mdm_declaration" "app-v1" {
  artifact_id = zentral_mdm_artifact.app.id
  source = jsonencode({
    Type        = "com.apple.configuration.app.managed",
    Identifier  = "com.example.app",
    ServerToken = "8e0d4a3a-3d1f-4c1f-9a2e-6f4b1c5d7e90",
    Payload = {
      BundleID = "com.example.app"
      AppConfig = {
        DataAssetReference = "ztl:${zentral_mdm_artifact.app-config.id}"
      }
    }
  })
  ios     = true
  ipados  = true
  version = 1
}
