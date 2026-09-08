package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTurboRecurringJobResource(t *testing.T) {
	name := acctest.RandString(12)
	tagName := acctest.RandString(12)
	excludedTagName := acctest.RandString(12)
	resourceName := "zentral_turbo_recurring_job.test"
	cfgResourceName := "zentral_turbo_configuration.test"
	scriptResourceName := "zentral_turbo_script.test"
	tagResourceName := "zentral_tag.test"
	excludedTagResourceName := "zentral_tag.excluded"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTurboRecurringJobResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "configuration_id", cfgResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "job_id", scriptResourceName, "job_id"),
					resource.TestCheckNoResourceAttr(
						resourceName, "interval"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "serial_numbers.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_serial_numbers.#", "0"),
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
				Config: testAccTurboRecurringJobResourceConfigFull(name, tagName, excludedTagName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "configuration_id", cfgResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "job_id", scriptResourceName, "job_id"),
					resource.TestCheckResourceAttr(
						resourceName, "interval", "3600"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "tag_ids.*", tagResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "excluded_tag_ids.*", excludedTagResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "serial_numbers.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "serial_numbers.*", "un"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "serial_numbers.*", "deux"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_serial_numbers.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "excluded_serial_numbers.*", "trois"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Remove the tags and the serial numbers from the config
			{
				Config: testAccTurboRecurringJobResourceConfigFull(name, tagName, excludedTagName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "serial_numbers.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_serial_numbers.#", "0"),
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

func testAccTurboRecurringJobResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_turbo_configuration" "test" {
  name = %[1]q
}

resource "zentral_turbo_script" "test" {
  name   = %[1]q
  source = "echo ok"
}

resource "zentral_turbo_recurring_job" "test" {
  configuration_id = zentral_turbo_configuration.test.id
  job_id           = zentral_turbo_script.test.job_id
}
`, name)
}

func testAccTurboRecurringJobResourceConfigFull(name string, tagName string, excludedTagName string, withSets bool) string {
	sets := ""
	if withSets {
		sets = `tag_ids                 = [zentral_tag.test.id]
  excluded_tag_ids        = [zentral_tag.excluded.id]
  serial_numbers          = ["un", "deux"]
  excluded_serial_numbers = ["trois"]`
	}
	return fmt.Sprintf(`
resource "zentral_turbo_configuration" "test" {
  name = %[1]q
}

resource "zentral_turbo_script" "test" {
  name   = %[1]q
  source = "echo ok"
}

resource "zentral_taxonomy" "test" {
  name = %[1]q
}

resource "zentral_tag" "test" {
  taxonomy_id = zentral_taxonomy.test.id
  name        = %[2]q
}

resource "zentral_tag" "excluded" {
  taxonomy_id = zentral_taxonomy.test.id
  name        = %[3]q
}

resource "zentral_turbo_recurring_job" "test" {
  configuration_id        = zentral_turbo_configuration.test.id
  job_id                  = zentral_turbo_script.test.job_id
  interval                = 3600
  %[4]s
}
`, name, tagName, excludedTagName, sets)
}
