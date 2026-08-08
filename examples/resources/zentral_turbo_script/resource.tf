resource "zentral_tag" "turbo" {
  name = "Turbo"
}

resource "zentral_turbo_script" "gatekeeper" {
  name        = "Gatekeeper enabled"
  description = "Fails if Gatekeeper is disabled"
  # trimspace: the API strips the surrounding whitespace, the heredoc trailing newline included
  source = trimspace(<<-EOT
    #!/bin/zsh
    /usr/sbin/spctl --status | grep -q "assessments enabled"
  EOT
  )

  # tags the machine on success, untags it on failure
  tag_id = zentral_tag.turbo.id

  # run as a compliance check
  compliance_check_enabled = true

  # only run on Apple Silicon, macOS 14 and later
  arch_amd64     = false
  arch_arm64     = true
  min_os_version = "14.0"
}
