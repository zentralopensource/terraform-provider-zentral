resource "zentral_turbo_configuration" "default" {
  name = "Default"
}

resource "zentral_meta_business_unit" "default" {
  name = "Default"
}

resource "zentral_tag" "turbo" {
  name = "Turbo"
}

resource "zentral_turbo_enrollment" "default" {
  configuration_id      = zentral_turbo_configuration.default.id
  meta_business_unit_id = zentral_meta_business_unit.default.id
  tag_ids               = [zentral_tag.turbo.id]
}
