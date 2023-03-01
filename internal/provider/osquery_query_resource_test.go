package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOsqueryQueryResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_osquery_query.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccOsqueryQueryResourceConfigBare(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "sql", "SELECT * FROM users;"),
					resource.TestCheckResourceAttr(
						resourceName, "platforms.#", "0"),
					resource.TestCheckNoResourceAttr(
						resourceName, "minimum_osquery_version"),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "value", ""),
					resource.TestCheckResourceAttr(
						resourceName, "version", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "compliance_check_enabled", "false"),
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
				Config: testAccOsqueryQueryResourceConfigFull(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "sql", "SELECT 'FAILED' AS ztl_status, 'No reason!' AS why;"),
					resource.TestCheckResourceAttr(
						resourceName, "platforms.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "platforms.*", "darwin"),
					resource.TestCheckResourceAttr(
						resourceName, "minimum_osquery_version", "0.1.0"),
					resource.TestCheckResourceAttr(
						resourceName, "description", "A compliance check that always fails"),
					resource.TestCheckResourceAttr(
						resourceName, "value", "Not much"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "2"),
					resource.TestCheckResourceAttr(
						resourceName, "compliance_check_enabled", "true"),
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

func testAccOsqueryQueryResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_query" "test" {
  name = %[1]q
  sql  = "SELECT * FROM users;"
}
`, name)
}

func testAccOsqueryQueryResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_query" "test" {
  name                     = %[1]q
  sql                      = "SELECT 'FAILED' AS ztl_status, 'No reason!' AS why;"
  platforms                = ["darwin"]
  minimum_osquery_version  = "0.1.0"
  description              = "A compliance check that always fails"
  value                    = "Not much"
  compliance_check_enabled = true
}
`, name)
}
