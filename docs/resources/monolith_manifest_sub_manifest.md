---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zentral_monolith_manifest_sub_manifest Resource - terraform-provider-zentral"
subcategory: ""
description: |-
  The resource zentral_monolith_manifest_sub_manifest manages Monolith manifest sub manifests.
---

# zentral_monolith_manifest_sub_manifest (Resource)

The resource `zentral_monolith_manifest_sub_manifest` manages Monolith manifest sub manifests.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `manifest_id` (Number) ID of the manifest.
- `sub_manifest_id` (Number) ID of the sub manifest.

### Optional

- `tag_ids` (Set of Number) The `ID`s of the tags used to scope the sub manifest.

### Read-Only

- `id` (Number) `ID` of the manifest sub manifest.
