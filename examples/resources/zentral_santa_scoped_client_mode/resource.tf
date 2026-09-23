resource "zentral_santa_configuration" "default" {
  name        = "Default"
  client_mode = "MONITOR"
}

resource "zentral_tag" "santa_lockdown" {
  name = "Santa lockdown"
}

resource "zentral_santa_scoped_client_mode" "lockdown" {
  configuration_id = zentral_santa_configuration.default.id
  name             = "Lockdown"
  description      = "Machines tagged for lockdown, except the IT team test machine"
  client_mode      = "LOCKDOWN"

  # the machines in scope get their own block notification button
  event_detail_source = "CUSTOM"
  event_detail_url    = "https://www.example.com/santa/blocked/?serial=%serial%"
  event_detail_text   = "Request an exception"

  tag_ids                 = [zentral_tag.santa_lockdown.id]
  excluded_serial_numbers = ["0123456789"]
}
