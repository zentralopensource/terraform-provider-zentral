terraform {
  required_providers {
    zentral = {
      source = "zentralopensource/zentral"
    }
  }
}

// configure the provider
provider "zentral" {
  // URL where the API endpoints are mounted in the Zentral deployment.
  // The ZTL_API_BASE_URL environment variable can be used instead.
  base_url = "https://zentral.example.com/api/"

  // Zentral service account (better) or user API token.
  // This is a secret, it must be managed using a variable.
  // The ZTL_API_TOKEN environment variable can be used instead.
  token = var.api_token
}

// add a "Hamburg" business unit
resource "zentral_meta_business_unit" "hamburg" {
  name = "Hamburg"
}

// add a "Department" taxonomy
resource "zentral_taxonomy" "department" {
  name = "Department"
}

// add a "Sales" tag in the "Department" taxonomy
resource "zentral_tag" "sales_dept" {
  taxonomy_id = zentral_taxonomy.department.id
  name        = "sales"
  color       = "67d6af"
}
