package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOsqueryQueryDataSource(t *testing.T) {
	q1Name := acctest.RandString(12)
	q2Name := acctest.RandString(12)
	q1ResourceName := "zentral_osquery_query.check1"
	q2ResourceName := "zentral_osquery_query.check2"
	ds1ResourceName := "data.zentral_osquery_query.check1_by_name"
	ds2ResourceName := "data.zentral_osquery_query.check2_by_id"
	packResourceName := "zentral_osquery_pack.check2"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOsqueryQueryDataSourceConfig(q1Name, q2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", q1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", q1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "sql", "SELECT * FROM users;"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "platforms.#", "0"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "minimum_osquery_version"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "description", ""),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "value", ""),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "version", "1"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "compliance_check_enabled", "false"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "scheduling"),
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", q2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", q2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "sql", "SELECT 'FAILED' AS ztl_status, 'No reason!' AS why;"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "platforms.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "platforms.*", "darwin"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "minimum_osquery_version", "0.1.0"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "description", "A compliance check that always fails"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "value", "Not much"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "version", "1"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "compliance_check_enabled", "true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "scheduling.can_be_denylisted", "false"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "scheduling.interval", "161"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "scheduling.log_removed_actions", "false"),
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "scheduling.pack_id", packResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "scheduling.shard", "100"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "scheduling.snapshot_mode", "true"),
				),
			},
		},
	})
}

func testAccOsqueryQueryDataSourceConfig(q1Name string, q2Name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_query" "check1" {
  name = %[1]q
  sql  = "SELECT * FROM users;"
}

resource "zentral_osquery_pack" "check2" {
  name = %[1]q
}

resource "zentral_osquery_query" "check2" {
  name                     = %[2]q
  sql                      = "SELECT 'FAILED' AS ztl_status, 'No reason!' AS why;"
  platforms                = ["darwin"]
  minimum_osquery_version  = "0.1.0"
  description              = "A compliance check that always fails"
  value                    = "Not much"
  compliance_check_enabled = true
  scheduling = {
    can_be_denylisted   = false,
    interval            = 161,
    log_removed_actions = false,
    pack_id             = zentral_osquery_pack.check2.id
    shard               = 100
    snapshot_mode       = true
  }
}

data "zentral_osquery_query" "check1_by_name" {
  name = zentral_osquery_query.check1.name
}

data "zentral_osquery_query" "check2_by_id" {
  id = zentral_osquery_query.check2.id
}
`, q1Name, q2Name)
}
