package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStoreResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_store.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccStoreResourceConfigHTTPBase(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "admin_console", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "events_url_authorized_role_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.included_event_filters.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.excluded_event_filters.#", "0"),
					// HTTP
					resource.TestCheckResourceAttr(
						resourceName, "backend", "HTTP"),
					resource.TestCheckResourceAttr(
						resourceName, "http.endpoint_url", "https://www.example.com/post"),
					resource.TestCheckNoResourceAttr(
						resourceName, "http.username"),
					resource.TestCheckNoResourceAttr(
						resourceName, "http.password"),
					resource.TestCheckResourceAttr(
						resourceName, "http.headers.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "http.concurrency", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "http.request_timeout", "120"),
					resource.TestCheckResourceAttr(
						resourceName, "http.max_retries", "3"),
					resource.TestCheckResourceAttr(
						resourceName, "http.verify_tls", "true"),
					// SPLUNK
					resource.TestCheckNoResourceAttr(
						resourceName, "splunk"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccStoreResourceConfigHTTPFull(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "First description"),
					resource.TestCheckResourceAttr(
						resourceName, "admin_console", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "events_url_authorized_role_ids.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "events_url_authorized_role_ids.0", "6"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.included_event_filters.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.included_event_filters.0.tags.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.included_event_filters.0.tags.0", "zentral"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.included_event_filters.0.event_type.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.included_event_filters.0.routing_key.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.excluded_event_filters.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.excluded_event_filters.0.tags.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.excluded_event_filters.0.event_type.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.excluded_event_filters.0.event_type.0", "zentral_logout"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.excluded_event_filters.0.routing_key.#", "0"),
					// HTTP
					resource.TestCheckResourceAttr(
						resourceName, "backend", "HTTP"),
					resource.TestCheckResourceAttr(
						resourceName, "http.endpoint_url", "https://www.example.com/post"),
					resource.TestCheckResourceAttr(
						resourceName, "http.username", "yolo"),
					resource.TestCheckResourceAttr(
						resourceName, "http.password", "fomo"),
					resource.TestCheckResourceAttr(
						resourceName, "http.headers.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "http.headers.0.name", "X-Custom-Header"),
					resource.TestCheckResourceAttr(
						resourceName, "http.headers.0.value", "Value"),
					resource.TestCheckResourceAttr(
						resourceName, "http.concurrency", "2"),
					resource.TestCheckResourceAttr(
						resourceName, "http.request_timeout", "123"),
					resource.TestCheckResourceAttr(
						resourceName, "http.max_retries", "4"),
					resource.TestCheckResourceAttr(
						resourceName, "http.verify_tls", "false"),
					// SPLUNK
					resource.TestCheckNoResourceAttr(
						resourceName, "splunk"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccStoreResourceConfigSplunkBase(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "Second description"),
					resource.TestCheckResourceAttr(
						resourceName, "admin_console", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "events_url_authorized_role_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.included_event_filters.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.excluded_event_filters.#", "0"),
					// SPLUNK
					resource.TestCheckResourceAttr(
						resourceName, "backend", "SPLUNK"),
					// HEC
					resource.TestCheckResourceAttr(
						resourceName, "splunk.hec_url", "https://www.example.com/hec"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.hec_token", "HECToken"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.hec_extra_headers.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.hec_request_timeout", "300"),
					resource.TestCheckNoResourceAttr(
						resourceName, "splunk.hec_index"),
					resource.TestCheckNoResourceAttr(
						resourceName, "splunk.hec_source"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.computer_name_as_host_sources.#", "0"),
					resource.TestCheckNoResourceAttr(
						resourceName, "splunk.custom_host_field"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.serial_number_field", "machine_serial_number"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.batch_size", "1"),
					// Events URLs
					resource.TestCheckNoResourceAttr(
						resourceName, "splunk.search_app_url"),
					// Events search
					resource.TestCheckNoResourceAttr(
						resourceName, "splunk.search_url"),
					resource.TestCheckNoResourceAttr(
						resourceName, "splunk.search_token"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.search_extra_headers.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.search_request_timeout", "300"),
					resource.TestCheckNoResourceAttr(
						resourceName, "splunk.search_index"),
					resource.TestCheckNoResourceAttr(
						resourceName, "splunk.search_source"),
					// Common
					resource.TestCheckResourceAttr(
						resourceName, "splunk.verify_tls", "true"),
					// HTTP
					resource.TestCheckNoResourceAttr(
						resourceName, "http"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccStoreResourceConfigSplunkFull(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "First description"),
					resource.TestCheckResourceAttr(
						resourceName, "admin_console", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "events_url_authorized_role_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.included_event_filters.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "event_filters.excluded_event_filters.#", "0"),
					// SPLUNK
					resource.TestCheckResourceAttr(
						resourceName, "backend", "SPLUNK"),
					// HEC
					resource.TestCheckResourceAttr(
						resourceName, "splunk.hec_url", "https://www.example.com/hec"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.hec_token", "HECToken"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.hec_extra_headers.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.hec_extra_headers.0.name", "X-Custom-HEC-Header"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.hec_extra_headers.0.value", "HECHeaderValue"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.hec_request_timeout", "123"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.hec_index", "HECIndex"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.hec_source", "HECSource"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.computer_name_as_host_sources.#", "2"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.computer_name_as_host_sources.0", "osquery"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.computer_name_as_host_sources.1", "munki"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.custom_host_field", "custom_host_field"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.serial_number_field", "serial_number"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.batch_size", "50"),
					// Events URLs
					resource.TestCheckResourceAttr(
						resourceName, "splunk.search_app_url", "https://www.example.com/search_app"),
					// Events search
					resource.TestCheckResourceAttr(
						resourceName, "splunk.search_url", "https://www.example.com/search"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.search_token", "SearchToken"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.search_extra_headers.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.search_extra_headers.0.name", "X-Custom-Search-Header"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.search_extra_headers.0.value", "SearchHeaderValue"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.search_request_timeout", "234"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.search_index", "SearchIndex"),
					resource.TestCheckResourceAttr(
						resourceName, "splunk.search_source", "SearchSource"),
					// Common
					resource.TestCheckResourceAttr(
						resourceName, "splunk.verify_tls", "false"),
					// HTTP
					resource.TestCheckNoResourceAttr(
						resourceName, "http"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccStoreResourceConfigHTTPBase(name string) string {
	return fmt.Sprintf(`
resource "zentral_store" "test" {
  name    = %[1]q
  backend = "HTTP"
  http = {
    endpoint_url = "https://www.example.com/post"
  }
}
`, name)
}

func testAccStoreResourceConfigHTTPFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_store" "test" {
  name                           = %[1]q
  description                    = "First description"
  admin_console                  = true
  # TODO: replace with zentral_role when available
  events_url_authorized_role_ids = [6]
  event_filters = {
    included_event_filters = [
      {
        tags = ["zentral"]
      }
    ]
    excluded_event_filters = [
      {
        event_type = ["zentral_logout"]
      }
    ]
  }
  backend = "HTTP"
  http = {
    endpoint_url = "https://www.example.com/post"
    username     = "yolo"
    password     = "fomo"
    headers = [
      { name = "X-Custom-Header"
      value = "Value" }
    ]
    concurrency     = 2
    request_timeout = 123
    max_retries     = 4
    verify_tls      = false
  }
}
`, name)
}

func testAccStoreResourceConfigSplunkBase(name string) string {
	return fmt.Sprintf(`
resource "zentral_store" "test" {
  name        = %[1]q
  description = "Second description"
  backend     = "SPLUNK"
  splunk = {
    hec_url   = "https://www.example.com/hec"
    hec_token = "HECToken"
  }
}
`, name)
}

func testAccStoreResourceConfigSplunkFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_store" "test" {
  name        = %[1]q
  description = "First description"
  backend     = "SPLUNK"
  splunk = {
    hec_url   = "https://www.example.com/hec"
    hec_token = "HECToken"
    hec_extra_headers = [
      {
        name  = "X-Custom-HEC-Header"
        value = "HECHeaderValue"
      }
    ]
    hec_request_timeout           = 123
    hec_index                     = "HECIndex"
    hec_source                    = "HECSource"
    computer_name_as_host_sources = ["osquery", "munki"]
    custom_host_field             = "custom_host_field"
    serial_number_field           = "serial_number"
    batch_size                    = 50
    search_app_url                = "https://www.example.com/search_app"
    search_url                    = "https://www.example.com/search"
    search_token                  = "SearchToken"
    search_extra_headers = [
      {
        name  = "X-Custom-Search-Header"
        value = "SearchHeaderValue"
      }
    ]
    search_request_timeout = 234
    search_index           = "SearchIndex"
    search_source          = "SearchSource"
    verify_tls             = false
  }
}
`, name)
}
