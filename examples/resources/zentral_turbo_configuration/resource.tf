resource "zentral_turbo_configuration" "default" {
  name = "Default"
}

resource "zentral_turbo_configuration" "tuned" {
  name                    = "Tuned"
  description             = "Faster refresh, smaller result batches"
  collect_inventory       = true
  inventory_interval      = 3600
  default_check_interval  = 7200
  config_refresh_interval = 300
  results_batch_size      = 50
}
