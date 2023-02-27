package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOsqueryATCDataSource(t *testing.T) {
	a1Name := acctest.RandString(12)
	a2Name := acctest.RandString(12)
	a1ResourceName := "zentral_osquery_atc.check1"
	a2ResourceName := "zentral_osquery_atc.check2"
	ds1ResourceName := "data.zentral_osquery_atc.check1_by_name"
	ds2ResourceName := "data.zentral_osquery_atc.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOsqueryATCDataSourceConfig(a1Name, a2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", a1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", a1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "description", "Access the Google Santa rules.db"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "table_name", fmt.Sprintf("%s_santa_rules", a1Name)),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "columns.#", "1"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "columns.0", "identifier"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "platforms.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						ds1ResourceName, "platforms.*", "darwin"),
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", a2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", a2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "description", "Access the Google Santa rules.db"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "table_name", fmt.Sprintf("%s_santa_rules", a2Name)),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "query", "SELECT * FROM rules;"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "path", "/var/db/santa/rules.db"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "columns.#", "5"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "columns.0", "identifier"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "columns.1", "state"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "columns.2", "type"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "columns.3", "custommsg"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "columns.4", "timestamp"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "platforms.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "platforms.*", "darwin"),
				),
			},
		},
	})
}

func testAccOsqueryATCDataSourceConfig(a1Name string, a2Name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_atc" "check1" {
  name        = %[1]q
  description = "Access the Google Santa rules.db"
  table_name  = "%[1]s_santa_rules"
  query       = "SELECT identifier FROM rules;"
  path        = "/var/db/santa/rules.db"
  columns     = ["identifier"]
  platforms   = ["darwin"]
}

resource "zentral_osquery_atc" "check2" {
  name        = %[2]q
  description = "Access the Google Santa rules.db"
  table_name  = "%[2]s_santa_rules"
  query       = "SELECT * FROM rules;"
  path        = "/var/db/santa/rules.db"
  columns     = ["identifier", "state", "type", "custommsg", "timestamp"]
  platforms   = ["darwin"]
}

data "zentral_osquery_atc" "check1_by_name" {
  name = zentral_osquery_atc.check1.name
}

data "zentral_osquery_atc" "check2_by_id" {
  id = zentral_osquery_atc.check2.id
}
`, a1Name, a2Name)
}
