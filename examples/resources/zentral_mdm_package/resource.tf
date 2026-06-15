# An MDM package: a .pkg the server downloads from `source_uri`, verifies
# against `sha256`, and serves to devices via the Apple ManifestURL contract.
# The file is immutable post-create — changing `source_uri` or `sha256`
# forces replacement.
resource "zentral_mdm_package" "internal-tool" {
  name        = "InternalTool"
  description = "Internal IT tooling, distributed via MDM."
  source_uri  = "s3://acme-mdm-packages/internal-tool-1.2.3.pkg"
  sha256      = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
}

# Distribute the package via a `com.apple.configuration.package` DDM
# declaration. The declaration carries a `Configuration` artifact on the
# `Device` channel, scoped to macOS (this declaration type is macOS 26+ only).
resource "zentral_mdm_artifact" "internal-tool" {
  name      = "Internal Tool"
  type      = "Configuration"
  channel   = "Device"
  platforms = ["macOS"]
}

resource "zentral_mdm_declaration" "internal-tool-v1" {
  artifact_id = zentral_mdm_artifact.internal-tool.id
  source = jsonencode({
    Type        = "com.apple.configuration.package"
    Identifier  = "com.example.internal-tool"
    ServerToken = "a3f4c7e2-1d8b-4f9c-8e7a-2b6d5c1f0e93"
    Payload = {
      # `ztl:<package-uuid>` is resolved server-side: at declaration-build
      # time Zentral mints a signed ManifestURL bound to (MDM enrollment
      # session, package). The signature is deterministic, so the URL stays
      # stable across regenerations and is valid for the lifetime of the
      # MDM enrollment session.
      ManifestURL = "ztl:${zentral_mdm_package.internal-tool.id}"
      InstallBehavior = {
        Install = "Required"
      }
    }
  })
  macos   = true
  version = 1
}
