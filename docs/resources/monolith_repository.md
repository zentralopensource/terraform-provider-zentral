---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zentral_monolith_repository Resource - terraform-provider-zentral"
subcategory: ""
description: |-
  The resource zentral_monolith_repository manages Monolith repositories.
---

# zentral_monolith_repository (Resource)

The resource `zentral_monolith_repository` manages Monolith repositories.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `backend` (String) Repository backend.
- `name` (String) Name of the repository.

### Optional

- `meta_business_unit_id` (Number) The `ID` of the meta business unit this repository is restricted to.
- `s3` (Attributes) S3 backend parameters. (see [below for nested schema](#nestedatt--s3))

### Read-Only

- `id` (Number) `ID` of the repository.

<a id="nestedatt--s3"></a>
### Nested Schema for `s3`

Required:

- `bucket` (String) Name of the S3 bucket.

Optional:

- `access_key_id` (String) AWS access key ID.
- `assume_role_arn` (String) ARN of the IAM role to assume.
- `cloudfront_domain` (String) Cloudfront domain.
- `cloudfront_key_id` (String) Cloudfront key ID.
- `cloudfront_privkey_pem` (String) Cloudfront private key in PEM form.
- `endpoint_url` (String) S3 endpoint URL.
- `prefix` (String) Prefix of the Munki repository in the S3 bucket.
- `region_name` (String) Name of the S3 bucket region.
- `secret_access_key` (String, Sensitive) AWS secret access key.
- `signature_version` (String) Version of the AWS request signature to use.