package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// The bounds of the device OS version range are refused in the plan, before any request.
func TestMDMSoftwareUpdateEnforcementValidation(t *testing.T) {
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
resource "zentral_mdm_software_update_enforcement" "test" {
  name                  = "test"
  platforms             = ["macOS"]
  max_os_version        = "27"
  min_device_os_version = "26.1.0.1"
}
`, `min_device_os_version\s+must\s+be\s+a\s+version\s+with\s+at\s+most\s+three\s+components`),
			step(`
resource "zentral_mdm_software_update_enforcement" "test" {
  name                  = "test"
  platforms             = ["macOS"]
  max_os_version        = "27"
  max_device_os_version = "26.1 (a)"
}
`, `max_device_os_version\s+must\s+be\s+a\s+version\s+with\s+at\s+most\s+three\s+components`),
			step(`
resource "zentral_mdm_software_update_enforcement" "test" {
  name                  = "test"
  platforms             = ["macOS"]
  max_os_version        = "27"
  min_device_os_version = "0.0"
}
`, `min_device_os_version\s+must\s+not\s+be\s+0`),
			step(`
resource "zentral_mdm_software_update_enforcement" "test" {
  name                  = "test"
  platforms             = ["macOS"]
  max_os_version        = "27"
  min_device_os_version = "27"
  max_device_os_version = "26.4"
}
`, `max_device_os_version\s+must\s+be\s+greater\s+than\s+min_device_os_version`),
			step(`
resource "zentral_mdm_software_update_enforcement" "test" {
  name                  = "test"
  platforms             = ["macOS"]
  max_os_version        = "27"
  min_device_os_version = "26"
  max_device_os_version = "26.0"
}
`, `max_device_os_version\s+must\s+be\s+greater\s+than\s+min_device_os_version`),
		},
	})
}
