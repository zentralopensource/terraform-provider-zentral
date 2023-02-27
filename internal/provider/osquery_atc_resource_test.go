package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOsqueryATCResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_osquery_atc.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccOsqueryATCResourceCreation(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "table_name", fmt.Sprintf("%s_santa_rul", firstName)),
					resource.TestCheckResourceAttr(
						resourceName, "query", "SELECT identifier, state, type FROM rules;"),
					resource.TestCheckResourceAttr(
						resourceName, "path", "/var/db/santa/rules.db"),
					resource.TestCheckResourceAttr(
						resourceName, "columns.#", "3"),
					resource.TestCheckResourceAttr(
						resourceName, "columns.0", "identifier"),
					resource.TestCheckResourceAttr(
						resourceName, "columns.1", "state"),
					resource.TestCheckResourceAttr(
						resourceName, "columns.2", "type"),
					resource.TestCheckResourceAttr(
						resourceName, "platforms.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "platforms.*", "darwin"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "platforms.*", "linux"),
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
				Config: testAccOsqueryATCResourceUpdate(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "Access the Google Santa rules.db"),
					resource.TestCheckResourceAttr(
						resourceName, "table_name", fmt.Sprintf("%s_santa_rules", secondName)),
					resource.TestCheckResourceAttr(
						resourceName, "query", "SELECT * FROM rules;"),
					resource.TestCheckResourceAttr(
						resourceName, "path", "/var/db/santa/rules.db"),
					resource.TestCheckResourceAttr(
						resourceName, "columns.#", "5"),
					resource.TestCheckResourceAttr(
						resourceName, "columns.0", "identifier"),
					resource.TestCheckResourceAttr(
						resourceName, "columns.1", "state"),
					resource.TestCheckResourceAttr(
						resourceName, "columns.2", "type"),
					resource.TestCheckResourceAttr(
						resourceName, "columns.3", "custommsg"),
					resource.TestCheckResourceAttr(
						resourceName, "columns.4", "timestamp"),
					resource.TestCheckResourceAttr(
						resourceName, "platforms.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "platforms.*", "darwin"),
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

func testAccOsqueryATCResourceCreation(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_atc" "test" {
  name       = %[1]q
  table_name = "%[1]s_santa_rul"
  query      = "SELECT identifier, state, type FROM rules;"
  path       = "/var/db/santa/rules.db"
  columns    = ["identifier", "state", "type"]
  platforms  = ["darwin", "linux"]
}
`, name)
}

func testAccOsqueryATCResourceUpdate(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_atc" "test" {
  name        = %[1]q
  description = "Access the Google Santa rules.db"
  table_name  = "%[1]s_santa_rules"
  query       = "SELECT * FROM rules;"
  path        = "/var/db/santa/rules.db"
  columns     = ["identifier", "state", "type", "custommsg", "timestamp"]
  platforms   = ["darwin"]
}
`, name)
}
