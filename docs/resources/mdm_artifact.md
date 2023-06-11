---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zentral_mdm_artifact Resource - terraform-provider-zentral"
subcategory: ""
description: |-
  The resource zentral_mdm_artifact manages MDM artifacts.
---

# zentral_mdm_artifact (Resource)

The resource `zentral_mdm_artifact` manages MDM artifacts.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `channel` (String) Channel of the artifact.
- `name` (String) Name of the artifact.
- `platforms` (Set of String) Platforms of the artifact.
- `type` (String) Type of the artifact.

### Optional

- `auto_update` (Boolean) If `true`, new version of this artifact will be automatically installed. Defaults to `true`.
- `install_during_setup_assistant` (Boolean) If `true`, this artifact will be installed during the setup assistant. Defaults to `false`.
- `reinstall_interval` (Number) In days, the time interval after which the artifact will be reinstalled. If `0`, the artifact will not be reinstalled. Defaults to `0`.
- `reinstall_on_os_update` (String) Possible values: `No`, `Major`, `Minor`, `Patch`. Defaults to `No`.
- `requires` (Set of String) `ID`s of the artifacts required by this artifact.

### Read-Only

- `id` (String) `ID` of the artifact.