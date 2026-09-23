resource "zentral_santa_configuration" "default" {
  name        = "Default"
  client_mode = "MONITOR"

  # the button of the block notification, sent in the preflight response
  event_detail_source = "CUSTOM"
  event_detail_url    = "https://www.example.com/santa/blocked/?serial=%serial%&file=%file_identifier%"
  event_detail_text   = "Request an exception"
}
