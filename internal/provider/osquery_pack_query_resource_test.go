package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOsqueryPackQueryResource(t *testing.T) {
	packName := acctest.RandString(12)
	resourceName := "zentral_osquery_pack_query.test"
	packResourceName := "zentral_osquery_pack.test"
	queryResourceName := "zentral_osquery_query.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccOsqueryPackQueryResourceConfigBare(packName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "pack_id", packResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "query_id", queryResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "slug", strings.ToLower(packName)),
					resource.TestCheckResourceAttr(
						resourceName, "interval", "60"),
					resource.TestCheckResourceAttr(
						resourceName, "log_removed_actions", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "snapshot_mode", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "shard", "50"),
					resource.TestCheckResourceAttr(
						resourceName, "can_be_denylisted", "true"),
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
				Config: testAccOsqueryPackQueryResourceConfigFull(packName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "pack_id", packResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "query_id", queryResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "slug", strings.ToLower(packName)),
					resource.TestCheckResourceAttr(
						resourceName, "interval", "1200"),
					resource.TestCheckResourceAttr(
						resourceName, "log_removed_actions", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "snapshot_mode", "true"),
					resource.TestCheckNoResourceAttr(
						resourceName, "shard"),
					resource.TestCheckResourceAttr(
						resourceName, "can_be_denylisted", "false"),
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

func testAccOsqueryPackQueryResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_pack" "test" {
  name = %[1]q
}

resource "zentral_osquery_query" "test" {
  name =  %[1]q
  sql  = "SELECT * FROM users;"
}

resource "zentral_osquery_pack_query" "test" {
  pack_id  = zentral_osquery_pack.test.id
  query_id = zentral_osquery_query.test.id
  interval = 60
  shard    = 50
}
`, name)
}

func testAccOsqueryPackQueryResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_pack" "test" {
  name = %[1]q
}

resource "zentral_osquery_query" "test" {
  name =  %[1]q
  sql  = "SELECT * FROM users;"
}

resource "zentral_osquery_pack_query" "test" {
  pack_id             = zentral_osquery_pack.test.id
  query_id            = zentral_osquery_query.test.id
  interval            = 1200
  log_removed_actions = false
  snapshot_mode       = true
  can_be_denylisted   = false
}
`, name)
}
