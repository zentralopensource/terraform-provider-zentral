package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOsqueryPackResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_osquery_pack.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccOsqueryPackResourceConfigBare(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "slug", strings.ToLower(firstName)),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "discovery_queries.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "discovery_queries.*", "SELECT pid FROM processes WHERE name = 'ldap';"),
					resource.TestCheckResourceAttr(
						resourceName, "shard", "50"),
					resource.TestCheckResourceAttr(
						resourceName, "event_routing_key", ""),
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
				Config: testAccOsqueryPackResourceConfigFull(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "slug", strings.ToLower(secondName)),
					resource.TestCheckResourceAttr(
						resourceName, "description", "description"),
					resource.TestCheckResourceAttr(
						resourceName, "discovery_queries.#", "0"),
					resource.TestCheckNoResourceAttr(
						resourceName, "shard"),
					resource.TestCheckResourceAttr(
						resourceName, "event_routing_key", "important"),
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

func testAccOsqueryPackResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_pack" "test" {
  name = %[1]q
  discovery_queries = ["SELECT pid FROM processes WHERE name = 'ldap';"]
  shard = 50
}
`, name)
}

func testAccOsqueryPackResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_pack" "test" {
  name              = %[1]q
  description       = "description"
  event_routing_key = "important"
}
`, name)
}
