---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zentral_mdm_filevault_config Data Source - terraform-provider-zentral"
subcategory: ""
description: |-
  The data source zentral_mdm_filevault_config allows details of a MDM FileVault configuration to be retrieved by its ID or its name.
---

# zentral_mdm_filevault_config (Data Source)

The data source `zentral_mdm_filevault_config` allows details of a MDM FileVault configuration to be retrieved by its `ID` or its `name`.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (Number) `ID` of the MDM FileVault configuration.
- `name` (String) Name of the FileVault configuration.

### Read-Only

- `at_login_only` (Boolean) If `true`, prevents requests for enabling FileVault at user logout time. Defaults to `false`.
- `bypass_attempts` (Number) The maximum number of times users can bypass enabling FileVault before being required to enable it to log in.
- `destroy_key_on_standby` (Boolean) Set to `true` to prevent storing the FileVault key across restarts. Defaults to `false`.
- `escrow_location_display_name` (String) Description of the location where the FDE PRK will be escrowed. This text will be inserted into the message the user sees when enabling FileVault.
- `prk_rotation_interval_days` (Number) The automatic PRK rotation interval in days. It has a maximum value of `365`. Defaults to `0` (no automatic rotation).
- `show_recovery_key` (Boolean) If `false`, prevents display of the personal recovery key to the user after FileVault is enabled. Defaults to `false`.