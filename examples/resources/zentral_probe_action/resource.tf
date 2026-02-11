resource "zentral_probe_action" "revoke-certs" {
  name        = "Revoke certificates"
  description = "Post event serial number and date to HTTP certificate revocation endpoint."
  backend     = "HTTP_POST"
  http_post = {
    url = "https://www.example.com/api/certificates/revoke/"
    headers = [
      {
        name  = "Authorization"
        value = "Bearer ${var.bearer_token}"
      }
    ]
    cel_transformation = trimspace(<<-EOT
    {
      "serial_number": "machine_serial_number" in metadata ? metadata.machine_serial_number : null,
      "not_after": metadata.created_at
    }
    EOT
    )
  }
}

resource "zentral_probe" "dep-enrollment-revoke-certs" {
  name        = "DEP enrollment → revoke certificates"
  description = "Trigger older certificates revocation on DEP device enrollment."
  active      = true
  action_ids  = [zentral_probe_action.revoke-certs.id]
  metadata_filters = [
    {
      event_types = ["enrollment_secret_verification"]
    }
  ]
  payload_filters = [
    [
      {
        attribute = "type"
        operator  = "IN"
        values = [
          "dep_enrollment",
        ]
      }
    ]
  ]
}
