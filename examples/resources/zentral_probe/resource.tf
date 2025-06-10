data "zentral_meta_business_unit" "hh" {
  name = "Hamburg"
}

data "zentral_tag" "vip" {
  name = "VIP"
}

resource "zentral_probe_action" "slack" {
  name    = "Post to #alerts Slack channel"
  backend = "HTTP_POST"
  http_post = {
    url = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
  }
}

resource "zentral_probe" "santa-block-slack" {
  name              = "Santa block events → Slack #alerts"
  description       = "Post Santa block events in the Hamburg business unit for the VIP tagged machines to Slack."
  incident_severity = 100
  active            = true
  action_ids        = [zentral_probe_action.slack.id]
  inventory_filters = [
    {
      meta_business_unit_ids = [zentral_meta_business_unit.hh.id]
      tag_ids                = [zentral_tag.vip.id]
    }
  ]
  metadata_filters = [
    {
      event_types = ["santa_event"]
    }
  ]
  payload_filters = [
    [
      {
        attribute = "decision"
        operator  = "IN"
        values = [
          "BLOCK_BINARY",
          "BLOCK_CERTIFICATE",
          "BLOCK_UNKNOWN",
          "BLOCK_CDHASH",
          "BLOCK_SIGNINGID",
          "BLOCK_TEAMID",
          "BLOCK_SCOPE",
        ]
      }
    ]
  ]
}
