# Tags

resource "zentral_taxonomy" "office" {
  name = "Office"
}

resource "zentral_tag" "office-ham" {
  taxonomy_id = zentral_taxonomy.office.id
  name        = "HAM"
}

resource "zentral_tag" "office-ldn" {
  taxonomy_id = zentral_taxonomy.office.id
  name        = "LDN"
}

# Existing Google Workspace connection

data "zentral_gws_connection" "this" {
  name = "ACME Google Workspace"
}

# Group Tag Mappings

resource "zentral_gws_group_tag_mapping" "office_ham" {
  connection_id = data.zentral_gws_connection.this.id
  group_email   = "acme-office-ham@example.com"
  tag_ids       = [zentral_tag.office-ham.id]
}

resource "zentral_gws_group_tag_mapping" "office_ldn" {
  connection_id = data.zentral_gws_connection.this.id
  group_email   = "acme-office-ldn@example.com"
  tag_ids       = [zentral_tag.office-ldn.id]
}
