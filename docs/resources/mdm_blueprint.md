---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zentral_mdm_blueprint Resource - terraform-provider-zentral"
subcategory: ""
description: |-
  The resource zentral_mdm_blueprint manages MDM blueprints.
---

# zentral_mdm_blueprint (Resource)

The resource `zentral_mdm_blueprint` manages MDM blueprints.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the blueprint.

### Optional

- `collect_apps` (String) Possible values: `NO`, `MANAGED_ONLY`, `ALL`.
- `collect_certificates` (String) Possible values: `NO`, `MANAGED_ONLY`, `ALL`.
- `collect_profiles` (String) Possible values: `NO`, `MANAGED_ONLY`, `ALL`.
- `inventory_interval` (Number) In seconds, the minimum interval between two inventory collection. Minimum 4h, maximum 7d, default 1d.

### Read-Only

- `id` (Number) `ID` of the blueprint.