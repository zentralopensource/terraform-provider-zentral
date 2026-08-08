resource "zentral_turbo_configuration" "default" {
  name = "Default"
}

resource "zentral_turbo_script" "gatekeeper" {
  name   = "Gatekeeper enabled"
  source = "/usr/sbin/spctl --status | grep -q \"assessments enabled\""
}

resource "zentral_tag" "turbo" {
  name = "Turbo"
}

resource "zentral_tag" "exempt" {
  name = "Turbo exempt"
}

resource "zentral_turbo_recurring_job" "gatekeeper" {
  configuration_id = zentral_turbo_configuration.default.id
  job_id           = zentral_turbo_script.gatekeeper.job_id

  # run every hour; leave unset to use the configuration default
  interval = 3600

  # scope: machines tagged Turbo, except those tagged Turbo exempt
  tag_ids          = [zentral_tag.turbo.id]
  excluded_tag_ids = [zentral_tag.exempt.id]
}
