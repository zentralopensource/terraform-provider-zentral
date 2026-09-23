package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// The attributes the server would clear or normalise are refused in the plan, before any request.
func TestSantaValidation(t *testing.T) {
	server := newOlderZentralServer(t)

	step := func(config string, expectError string) resource.TestStep {
		return resource.TestStep{
			Config:      testOlderZentralProviderConfig(server) + config,
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(expectError),
		}
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			step(`
resource "zentral_santa_configuration" "test" {
  name                = "test"
  event_detail_source = "CUSTOM"
}
`, `event_detail_url\s+is\s+required\s+when\s+event_detail_source\s+is\s+CUSTOM`),
			step(`
resource "zentral_santa_configuration" "test" {
  name             = "test"
  event_detail_url = "https://www.example.com/"
}
`, `event_detail_url\s+can\s+only\s+be\s+set\s+when\s+event_detail_source\s+is\s+CUSTOM`),
			step(`
resource "zentral_santa_configuration" "test" {
  name                = "test"
  event_detail_source = "NONE"
  event_detail_text   = "More info"
}
`, `event_detail_text\s+can\s+only\s+be\s+set\s+when\s+event_detail_source\s+is\s+CUSTOM\s+or\s+VOTING_PORTAL`),
			step(`
resource "zentral_santa_configuration" "test" {
  name                = "test"
  event_detail_source = "VOTING_PORTAL"
  event_detail_text   = "More info "
}
`, `must\s+not\s+start\s+or\s+end\s+with\s+whitespace`),
			step(`
resource "zentral_santa_scoped_client_mode" "test" {
  configuration_id  = 1
  name              = "test"
  client_mode       = "LOCKDOWN"
  event_detail_text = "More info"
}
`, `event_detail_text\s+can\s+only\s+be\s+set\s+when\s+event_detail_source\s+is\s+CUSTOM\s+or\s+VOTING_PORTAL`),
			step(`
resource "zentral_santa_scoped_client_mode" "test" {
  configuration_id = 1
  name             = "test"
  client_mode      = "LOCKDOWN"
  serial_numbers   = [" 012345678 "]
}
`, `must\s+not\s+start\s+or\s+end\s+with\s+whitespace`),
		},
	})
}
