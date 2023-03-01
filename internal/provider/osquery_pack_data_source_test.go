package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOsqueryPackDataSource(t *testing.T) {
	p1Name := acctest.RandString(12)
	p2Name := acctest.RandString(12)
	p1ResourceName := "zentral_osquery_pack.check1"
	p2ResourceName := "zentral_osquery_pack.check2"
	ds1ResourceName := "data.zentral_osquery_pack.check1_by_name"
	ds2ResourceName := "data.zentral_osquery_pack.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOsqueryPackDataSourceConfig(p1Name, p2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", p1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", p1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "slug", strings.ToLower(p1Name)),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "description", ""),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "discovery_queries.#", "0"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "shard"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "shard"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "event_routing_key", ""),
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", p2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", p2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "slug", strings.ToLower(p2Name)),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "description", "description"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "discovery_queries.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "discovery_queries.*", "SELECT pid FROM processes WHERE name = 'ldap';"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "shard", "50"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "event_routing_key", "important"),
				),
			},
		},
	})
}

func testAccOsqueryPackDataSourceConfig(p1Name string, p2Name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_pack" "check1" {
  name = %[1]q
}

resource "zentral_osquery_pack" "check2" {
  name                     = %[2]q
  description       = "description"
  discovery_queries = ["SELECT pid FROM processes WHERE name = 'ldap';"]
  shard             = 50
  event_routing_key = "important"
}

data "zentral_osquery_pack" "check1_by_name" {
  name = zentral_osquery_pack.check1.name
}

data "zentral_osquery_pack" "check2_by_id" {
  id = zentral_osquery_pack.check2.id
}
`, p1Name, p2Name)
}
