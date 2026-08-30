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
  source = base64encode(<<-EOT
    <?xml version="1.0" encoding="UTF-8"?>
    <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
    <plist version="1.0">
    <dict>
      <key>ApiURL</key>
      <string>https://api.example.com</string>
    </dict>
    </plist>
  EOT
  )
  ios     = true
  ipados  = true
  version = 1
}

# Or point at an object in an S3 bucket the server can read.
resource "zentral_mdm_artifact" "app-resources" {
  name      = "App Resources"
  type      = "Data Asset"
  channel   = "Device"
  platforms = ["iOS", "iPadOS"]
}

# `file_sha256` is required with `file_uri`, and the server downloads the file
# and verifies it against the digest, so replace both with your own.
resource "zentral_mdm_data_asset" "app-resources-v1" {
  artifact_id = zentral_mdm_artifact.app-resources.id
  type        = "ZIP"
  file_uri    = "s3://acme-mdm-assets/app-resources.v1.zip"
  file_sha256 = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
  ios         = true
  ipados      = true
  version     = 1
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
