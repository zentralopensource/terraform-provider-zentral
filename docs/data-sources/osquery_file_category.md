---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zentral_osquery_file_category Data Source - terraform-provider-zentral"
subcategory: ""
description: |-
  The data source zentral_osquery_file_category allows details of a Osquery file category to be retrieved by its ID or name.
---

# zentral_osquery_file_category (Data Source)

The data source `zentral_osquery_file_category` allows details of a Osquery file category to be retrieved by its `ID` or name.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (Number) `ID` of the Osquery file category.
- `name` (String) Name of the Osquery file category.

### Read-Only

- `access_monitoring` (Boolean) If `true`, FIM will include file access for this file category. Defaults to `false`.
- `description` (String) Description of the Osquery file category.
- `exclude_paths` (Set of String) Set of paths to exclude from the Osquery file category.
- `file_paths` (Set of String) Set of paths to include in the Osquery file category.
- `file_paths_queries` (Set of String) Set of queries returning paths to monitor as path columns in the results.
