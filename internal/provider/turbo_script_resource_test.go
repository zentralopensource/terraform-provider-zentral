package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTurboScriptResource(t *testing.T) {
	name := acctest.RandString(12)
	tagName := acctest.RandString(12)
	resourceName := "zentral_turbo_script.test"
	tagResourceName := "zentral_tag.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read with defaults
			{
				Config: testAccTurboScriptResourceConfigBare(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "source", "echo ok"),
					resource.TestCheckNoResourceAttr(
						resourceName, "tag_id"),
					resource.TestCheckResourceAttr(
						resourceName, "arch_amd64", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "arch_arm64", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "min_os_version", ""),
					resource.TestCheckResourceAttr(
						resourceName, "max_os_version", ""),
					resource.TestCheckResourceAttr(
						resourceName, "compliance_check_enabled", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "1"),
					resource.TestCheckResourceAttrSet(
						resourceName, "job_id"),
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
				Config: testAccTurboScriptResourceConfigFull(name, tagName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", fmt.Sprintf("%s-updated", name)),
					resource.TestCheckResourceAttr(
						resourceName, "description", "a description"),
					resource.TestCheckResourceAttr(
						resourceName, "source", "echo ko"),
					resource.TestCheckResourceAttrPair(
						resourceName, "tag_id", tagResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "arch_amd64", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "arch_arm64", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "min_os_version", "14.0"),
					resource.TestCheckResourceAttr(
						resourceName, "max_os_version", "16.0"),
					resource.TestCheckResourceAttr(
						resourceName, "compliance_check_enabled", "true"),
					resource.TestCheckResourceAttrSet(
						resourceName, "compliance_check_id"),
					// the source changed, so the version was bumped
					resource.TestCheckResourceAttr(
						resourceName, "version", "2"),
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

func testAccTurboScriptResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_turbo_script" "test" {
  name   = %[1]q
  source = "echo ok"
}
`, name)
}

func testAccTurboScriptResourceConfigFull(name string, tagName string) string {
	return fmt.Sprintf(`
resource "zentral_taxonomy" "test" {
  name = %[1]q
}

resource "zentral_tag" "test" {
  taxonomy_id = zentral_taxonomy.test.id
  name        = %[2]q
}

resource "zentral_turbo_script" "test" {
  name                     = "%[1]s-updated"
  description              = "a description"
  source                   = "echo ko"
  tag_id                   = zentral_tag.test.id
  arch_amd64               = true
  arch_arm64               = false
  min_os_version           = "14.0"
  max_os_version           = "16.0"
  compliance_check_enabled = true
}
`, name, tagName)
}

// A script update bumps the version of its job, but not the job itself. If job_id were not pinned to
// the state value it would be unknown in the plan, and the jobs referencing it would be replaced —
// re-running a one-time job on every machine that had already completed it.
func TestAccTurboScriptJobIDStable(t *testing.T) {
	name := acctest.RandString(12)
	resourceName := "zentral_turbo_script.test"
	jobResourceName := "zentral_turbo_one_time_job.test"

	var jobID string
	recordJobID := func(s *terraform.State) error {
		jobID = s.RootModule().Resources[jobResourceName].Primary.Attributes["id"]
		return nil
	}
	checkNotReplaced := func(s *terraform.State) error {
		if got := s.RootModule().Resources[jobResourceName].Primary.Attributes["id"]; got != jobID {
			return fmt.Errorf("the one-time job was replaced: %s -> %s", jobID, got)
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTurboScriptJobIDStableConfig(name, "echo ok"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version", "1"),
					recordJobID,
				),
			},
			// Update the script source and Read
			{
				Config: testAccTurboScriptJobIDStableConfig(name, "echo ko"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version", "2"),
					resource.TestCheckResourceAttrPair(
						jobResourceName, "job_id", resourceName, "job_id"),
					checkNotReplaced,
				),
			},
		},
	})
}

func testAccTurboScriptJobIDStableConfig(name string, source string) string {
	return fmt.Sprintf(`
resource "zentral_turbo_configuration" "test" {
  name = %[1]q
}

resource "zentral_turbo_script" "test" {
  name   = %[1]q
  source = %[2]q
}

resource "zentral_turbo_one_time_job" "test" {
  configuration_id = zentral_turbo_configuration.test.id
  job_id           = zentral_turbo_script.test.job_id
}
`, name, source)
}
