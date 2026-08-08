package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTurboMSCPCheckResource(t *testing.T) {
	ruleID := fmt.Sprintf("rule_%s", acctest.RandString(8))
	resourceName := "zentral_turbo_mscp_check.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read with defaults
			{
				Config: testAccTurboMSCPCheckResourceConfigBare(ruleID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "rule_id", ruleID),
					resource.TestCheckResourceAttr(
						resourceName, "baseline", ""),
					resource.TestCheckNoResourceAttr(
						resourceName, "odv_int"),
					resource.TestCheckNoResourceAttr(
						resourceName, "odv_string"),
					resource.TestCheckNoResourceAttr(
						resourceName, "odv_bool"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "1"),
					resource.TestCheckResourceAttrSet(
						resourceName, "compliance_check_id"),
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
			// Update and Read with an ODV override
			{
				Config: testAccTurboMSCPCheckResourceConfigODV(ruleID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "rule_id", ruleID),
					resource.TestCheckResourceAttr(
						resourceName, "baseline", ""),
					resource.TestCheckResourceAttr(
						resourceName, "odv_int", "10"),
					// an identity-bearing change bumps the version
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
			// Update and Read with a boolean ODV override
			{
				Config: testAccTurboMSCPCheckResourceConfigODVBool(ruleID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "rule_id", ruleID),
					resource.TestCheckNoResourceAttr(
						resourceName, "odv_int"),
					resource.TestCheckResourceAttr(
						resourceName, "odv_bool", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "3"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read with a baseline
			{
				Config: testAccTurboMSCPCheckResourceConfigBaseline(ruleID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "rule_id", ruleID),
					resource.TestCheckResourceAttr(
						resourceName, "baseline", "cis_lvl1"),
					resource.TestCheckNoResourceAttr(
						resourceName, "odv_int"),
					resource.TestCheckNoResourceAttr(
						resourceName, "odv_bool"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "4"),
				),
			},
		},
	})
}

func testAccTurboMSCPCheckResourceConfigBare(ruleID string) string {
	return fmt.Sprintf(`
resource "zentral_turbo_mscp_check" "test" {
  rule_id = %[1]q
}
`, ruleID)
}

func testAccTurboMSCPCheckResourceConfigODV(ruleID string) string {
	return fmt.Sprintf(`
resource "zentral_turbo_mscp_check" "test" {
  rule_id = %[1]q
  odv_int = 10
}
`, ruleID)
}

func testAccTurboMSCPCheckResourceConfigODVBool(ruleID string) string {
	// false, not true: a zero value must round-trip as a value, not as an absent override
	return fmt.Sprintf(`
resource "zentral_turbo_mscp_check" "test" {
  rule_id  = %[1]q
  odv_bool = false
}
`, ruleID)
}

func testAccTurboMSCPCheckResourceConfigBaseline(ruleID string) string {
	return fmt.Sprintf(`
resource "zentral_turbo_mscp_check" "test" {
  rule_id  = %[1]q
  baseline = "cis_lvl1"
}
`, ruleID)
}

// Same contract as for the scripts: an mSCP check update bumps the version of its job, not the job
// itself, so the jobs referencing its job_id must not be replaced.
func TestAccTurboMSCPCheckJobIDStable(t *testing.T) {
	name := acctest.RandString(12)
	ruleID := fmt.Sprintf("rule_%s", acctest.RandString(8))
	resourceName := "zentral_turbo_mscp_check.test"
	jobResourceName := "zentral_turbo_recurring_job.test"

	var jobID string
	recordJobID := func(s *terraform.State) error {
		jobID = s.RootModule().Resources[jobResourceName].Primary.Attributes["id"]
		return nil
	}
	checkNotReplaced := func(s *terraform.State) error {
		if got := s.RootModule().Resources[jobResourceName].Primary.Attributes["id"]; got != jobID {
			return fmt.Errorf("the recurring job was replaced: %s -> %s", jobID, got)
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTurboMSCPCheckJobIDStableConfig(name, ruleID, 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version", "1"),
					recordJobID,
				),
			},
			// Update the ODV override and Read
			{
				Config: testAccTurboMSCPCheckJobIDStableConfig(name, ruleID, 20),
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

func testAccTurboMSCPCheckJobIDStableConfig(name string, ruleID string, odvInt int) string {
	return fmt.Sprintf(`
resource "zentral_turbo_configuration" "test" {
  name = %[1]q
}

resource "zentral_turbo_mscp_check" "test" {
  rule_id = %[2]q
  odv_int = %[3]d
}

resource "zentral_turbo_recurring_job" "test" {
  configuration_id = zentral_turbo_configuration.test.id
  job_id           = zentral_turbo_mscp_check.test.job_id
}
`, name, ruleID, odvInt)
}
