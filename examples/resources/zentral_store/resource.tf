resource "zentral_store" "http-minimal" {
  name    = "HTTP-minimal"
  backend = "Example of a minimal HTTP store backend"
  http = {
    endpoint_url = "https://www.example.com/post"
  }
}

resource "zentral_store" "http-full" {
  name          = "HTTP-Full"
  description   = "Example of a fully configured HTTP store backend"
  admin_console = false
  event_filters = {
    included_event_filters = [
      {
        # events with the `osquery` tag are included
        tags = ["osquery"]
      }
    ]
    excluded_event_filters = [
      {
        # `osquery_request` events are excluded
        # this filter takes precedence
        event_type = ["osquery_request"]
      }
    ]
  }
  backend = "HTTP"
  http = {
    endpoint_url = "https://www.example.com/post"
    username     = "yolo"
    password     = var.http_basic_auth_password
    headers = [
      { name = "X-Custom-Header"
      value = "Value" }
    ]
    concurrency     = 2
    request_timeout = 123
    max_retries     = 4
    verify_tls      = true
  }
}

resource "zentral_store" "panther" {
  name        = "Panther"
  description = "Example of a Panther store backend"
  backend     = "PANTHER"
  panther = {
    endpoint_url = "https://logs.acme.runpanther.net/http/b1a5141b-d494-452f-88ca-a04adb25b861"
    bearer_token = "19cd744d-fc4c-48eb-be9f-abd6dd17a575"
    batch_size   = 100
  }
}

resource "zentral_store" "splunk-minimal" {
  name        = "Splunk-minimal"
  description = "Example of a minimal Splunk store backend"
  backend     = "SPLUNK"
  splunk = {
    # HEC
    hec_url   = "https://www.example.com/services/collector/event"
    hec_token = var.splunk_hec_token
  }
}

resource "zentral_store" "splunk-full" {
  name                           = "Splunk-full"
  description                    = "Example of a fully configured Splunk store backend"
  admin_console                  = true
  events_url_authorized_role_ids = [6]
  backend                        = "SPLUNK"
  splunk = {
    # HEC
    hec_url   = "https://www.example.com/services/collector/event"
    hec_token = var.splunk_hec_token
    hec_extra_headers = [
      {
        name  = "X-Custom-HEC-Header"
        value = "HECHeaderValue"
      }
    ]
    hec_request_timeout           = 123
    hec_index                     = "CustomIndex"
    hec_source                    = "CustomSource"
    computer_name_as_host_sources = ["osquery", "munki"]
    custom_host_field             = "custom_host_field"
    serial_number_field           = "serial_number"
    batch_size                    = 50
    # Search links displayed in the Zentral admin console
    search_app_url = "https://www.example.com/search_app"
    # Search API used to fetch the events
    # and display them in the Zentral admin console
    search_url   = "https://www.example.com/search"
    search_token = var.splunk_search_token
    search_extra_headers = [
      {
        name  = "X-Custom-Search-Header"
        value = "SearchHeaderValue"
      }
    ]
    search_request_timeout = 234
    search_index           = "CustomIndex"
    search_source          = "CustomSource"
    # for HEC and Search API requests
    verify_tls = true
  }
}
