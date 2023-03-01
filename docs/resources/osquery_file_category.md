---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zentral_osquery_file_category Resource - terraform-provider-zentral"
subcategory: ""
description: |-
  The resource zentral_osquery_file_category manages Osquery file categories.
---

# zentral_osquery_file_category (Resource)

The resource `zentral_osquery_file_category` manages Osquery file categories.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the Osquery file category.

### Optional

- `access_monitoring` (Boolean) If `true`, FIM will include file access for this file category. Defaults to `false`.
- `description` (String) Description of the Osquery file category.
- `exclude_paths` (Set of String) Set of paths to exclude from the Osquery file category.
- `file_paths` (Set of String) Set of paths to include in the Osquery file category.
- `file_paths_queries` (Set of String) Set of queries returning paths to monitor as path columns in the results.

### Read-Only

- `id` (Number) `ID` of the Osquery file category.

