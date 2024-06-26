---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zentral_mdm_push_certificate Data Source - terraform-provider-zentral"
subcategory: ""
description: |-
  The data source zentral_mdm_push_certificate allows details of a MDM push certificate to be retrieved by its ID or its name.
---

# zentral_mdm_push_certificate (Data Source)

The data source `zentral_mdm_push_certificate` allows details of a MDM push certificate to be retrieved by its `ID` or its `name`.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (Number) `ID` of the MDM push certificate.
- `name` (String) Name of the push certificate.

### Read-Only

- `certificate` (String) Push certificate in `PEM` form.
- `provisioning_uid` (String) Provisioning `UID` of the push certificate.
- `topic` (String) APNS topic of the push certificate.
