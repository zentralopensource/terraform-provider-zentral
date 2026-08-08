resource "zentral_turbo_configuration" "default" {
  name = "Default"
}

resource "zentral_turbo_script" "gatekeeper" {
  name   = "Gatekeeper enabled"
  source = "/usr/sbin/spctl --status | grep -q \"assessments enabled\""
}

resource "zentral_turbo_one_time_job" "gatekeeper_audit" {
  configuration_id = zentral_turbo_configuration.default.id
  job_id           = zentral_turbo_script.gatekeeper.job_id

  # deliver only within this window (UTC, no timezone suffix)
  not_before = "2026-09-01T09:00:00"
  not_after  = "2026-09-30T09:00:00"

  # scope: a specific fleet of serial numbers
  serial_numbers = ["C02XXYYZZ001", "C02XXYYZZ002"]
}
