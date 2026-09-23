resource "zentral_santa_configuration" "default" {
  name        = "Default"
  client_mode = "LOCKDOWN"
}

resource "zentral_tag" "developers" {
  name = "Developers"
}

resource "zentral_santa_scoped_path_regex" "developer_tools" {
  configuration_id = zentral_santa_configuration.default.id
  name             = "Developer tools"
  description      = "Allow the build output of the developers"
  policy           = "ALLOW"
  # anchored at the start of the path, no leading ^
  regex = "/Users/[^/]+/Developer/.+"

  tag_ids = [zentral_tag.developers.id]
}
