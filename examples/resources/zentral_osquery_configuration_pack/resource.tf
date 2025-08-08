resource "zentral_osquery_configuration" "default" {
  name = "Default"
}

resource "zentral_osquery_pack" "compliance-checks" {
  name = "Compliance checks"
}

resource "zentral_tag" "office" {
  name = "Office"
}

resource "zentral_tag" "server" {
  name = "Server"
}

resource "zentral_osquery_configuration_pack" "default-compliance-checks" {
  configuration_id = zentral_osquery_configuration.default.id
  pack_id          = zentral_osquery_pack.compliance-checks.id
  tag_ids          = [zentral_tag.office.id]
  excluded_tag_ids = [zentral_tag.server.id]
}
