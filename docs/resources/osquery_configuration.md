---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zentral_osquery_configuration Resource - terraform-provider-zentral"
subcategory: ""
description: |-
  The resource zentral_osquery_configuration manages Osquery configurations.
---

# zentral_osquery_configuration (Resource)

The resource `zentral_osquery_configuration` manages Osquery configurations.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the Osquery configuration.

### Optional

- `atc_ids` (Set of Number) List of the IDs of the ATCs to include in this configuration.
- `description` (String) Description of the Osquery configuration.
- `file_category_ids` (Set of Number) List of the IDs of the file categories to include in this configuration.
- `inventory` (Boolean) If `true`, Osquery is configured to collect inventory data. Defaults to `true`.
- `inventory_apps` (Boolean) If `true`, Osquery is configured to collect the applications. Defaults to `false`.
- `inventory_ec2` (Boolean) If `true`, Osquery is configured to collect the EC2 metadata. Defaults to `false`.
- `inventory_interval` (Number) Number of seconds to wait between collecting the inventory data.
- `options` (Map of String) A map of extra options to pass to Osquery in the flag file.

### Read-Only

- `id` (Number) `ID` of the Osquery configuration.
