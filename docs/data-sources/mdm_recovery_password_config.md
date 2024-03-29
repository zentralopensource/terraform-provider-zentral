---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zentral_mdm_recovery_password_config Data Source - terraform-provider-zentral"
subcategory: ""
description: |-
  The data source zentral_mdm_recovery_password_config allows details of a MDM recovery password configuration to be retrieved by its ID or its name.
---

# zentral_mdm_recovery_password_config (Data Source)

The data source `zentral_mdm_recovery_password_config` allows details of a MDM recovery password configuration to be retrieved by its `ID` or its `name`.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (Number) `ID` of the MDM recovery password configuration.
- `name` (String) Name of the recovery password configuration.

### Read-Only

- `dynamic_password` (Boolean) If `true`, a unique password is generated for each device. Defaults to `true`.
- `rotate_firmware_password` (Boolean) Set to `true` to rotate the firmware passwords. Defaults to `false`.
- `rotation_interval_days` (Number) The automatic recovery password rotation interval in days. It has a maximum value of `365`. Defaults to `0` (no automatic rotation).
- `static_password` (String) The  static password to set for all devices.
