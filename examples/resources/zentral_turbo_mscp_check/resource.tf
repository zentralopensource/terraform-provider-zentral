resource "zentral_turbo_mscp_check" "firewall" {
  rule_id = "os_firewall_enable"
}

# track a baseline's default ODV for the rule
resource "zentral_turbo_mscp_check" "firewall_cis" {
  rule_id  = "os_firewall_enable"
  baseline = "cis_lvl1"
}

# pin an explicit ODV override (mutually exclusive with baseline)
resource "zentral_turbo_mscp_check" "screensaver_timeout" {
  rule_id = "os_screensaver_idle_time_configure"
  odv_int = 1200
}
