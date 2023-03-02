package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOsqueryPackQueryDataSource(t *testing.T) {
	name := acctest.RandString(12)
	resourceName := "data.zentral_osquery_pack_query.check"
	packQueryResourceName := "zentral_osquery_pack_query.check"
	packResourceName := "zentral_osquery_pack.check"
	queryResourceName := "zentral_osquery_query.check"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOsqueryPackQueryDataSourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by ID
					resource.TestCheckResourceAttrPair(
						resourceName, "id", packQueryResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "pack_id", packResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "query_id", queryResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "slug", strings.ToLower(name)),
					resource.TestCheckResourceAttr(
						resourceName, "interval", "1200"),
					resource.TestCheckResourceAttr(
						resourceName, "log_removed_actions", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "snapshot_mode", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "shard", "50"),
					resource.TestCheckResourceAttr(
						resourceName, "can_be_denylisted", "false"),
				),
			},
		},
	})
}

func testAccOsqueryPackQueryDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_pack" "check" {
  name = %[1]q
}

resource "zentral_osquery_query" "check" {
  name = %[1]q
  sql  = "SELECT * FROM users;"
}
resource "zentral_osquery_pack_query" "check" {
  pack_id             = zentral_osquery_pack.check.id
  query_id            = zentral_osquery_query.check.id
  interval            = 1200
  log_removed_actions = false
  snapshot_mode       = true
  shard               = 50
  can_be_denylisted   = false
}

data "zentral_osquery_pack_query" "check" {
  id = zentral_osquery_pack_query.check.id
}
`, name)
}
