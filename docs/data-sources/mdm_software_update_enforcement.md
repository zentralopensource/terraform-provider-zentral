---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zentral_mdm_software_update_enforcement Data Source - terraform-provider-zentral"
subcategory: ""
description: |-
  The data source zentral_mdm_software_update_enforcement allows details of a MDM software update enforcement to be retrieved by its ID or its name.
---

# zentral_mdm_software_update_enforcement (Data Source)

The data source `zentral_mdm_software_update_enforcement` allows details of a MDM software update enforcement to be retrieved by its `ID` or its `name`.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (Number) `ID` of the MDM software update enforcement.
- `name` (String) Name of the software update enforcement.

### Read-Only

- `build_version` (String) The target build version to update the device to by the appropriate time.
- `delay_days` (Number) Number of days after a software update release before the device force installs it.
- `details_url` (String) The URL of a web page that shows details that the organization provides about the enforced update.
- `local_datetime` (String) The local date time value that specifies when to force install the software update.
- `local_time` (String) The local time value that specifies when to force install the software update.
- `max_os_version` (String) The maximum (excluded) target OS version to update the device to by the appropriate time.
- `os_version` (String) The target OS version to update the device to by the appropriate time.
- `platforms` (Set of String) The platforms this software update enforcement is scoped to.
- `tag_ids` (Set of Number) The `ID`s of the tags used to scope the software update enforcement.