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
	thirdName := acctest.RandString(12)
	packResourceName := "zentral_osquery_pack.test"
	resourceName := "zentral_osquery_query.test"
	tagResourceName := "zentral_tag.test"

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
					resource.TestCheckNoResourceAttr(
						resourceName, "tag_id"),
					resource.TestCheckNoResourceAttr(
						resourceName, "scheduling"),
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
				Config: testAccOsqueryQueryResourceConfigCC(secondName),
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
					resource.TestCheckNoResourceAttr(
						resourceName, "tag_id"),
					resource.TestCheckResourceAttr(
						resourceName, "scheduling.can_be_denylisted", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "scheduling.interval", "161"),
					resource.TestCheckResourceAttr(
						resourceName, "scheduling.log_removed_actions", "false"),
					resource.TestCheckResourceAttrPair(
						resourceName, "scheduling.pack_id", packResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "scheduling.shard", "10"),
					resource.TestCheckResourceAttr(
						resourceName, "scheduling.snapshot_mode", "true"),
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
				Config: testAccOsqueryQueryResourceConfigTag(thirdName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", thirdName),
					resource.TestCheckResourceAttr(
						resourceName, "sql", "SELECT 1;"),
					resource.TestCheckResourceAttr(
						resourceName, "platforms.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "platforms.*", "darwin"),
					resource.TestCheckResourceAttr(
						resourceName, "minimum_osquery_version", "0.1.0"),
					resource.TestCheckResourceAttr(
						resourceName, "description", "A query that always adds a tag"),
					resource.TestCheckResourceAttr(
						resourceName, "value", "Not much"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "3"),
					resource.TestCheckResourceAttr(
						resourceName, "compliance_check_enabled", "false"),
					resource.TestCheckResourceAttrPair(
						resourceName, "tag_id", tagResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "scheduling.can_be_denylisted", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "scheduling.interval", "161"),
					resource.TestCheckResourceAttr(
						resourceName, "scheduling.log_removed_actions", "false"),
					resource.TestCheckResourceAttrPair(
						resourceName, "scheduling.pack_id", packResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "scheduling.shard", "10"),
					resource.TestCheckResourceAttr(
						resourceName, "scheduling.snapshot_mode", "true"),
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

func testAccOsqueryQueryResourceConfigCC(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_pack" "test" {
  name = %[1]q
}

resource "zentral_osquery_query" "test" {
  name                     = %[1]q
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
    pack_id             = zentral_osquery_pack.test.id
    shard               = 10
    snapshot_mode       = true
  }
}
`, name)
}

func testAccOsqueryQueryResourceConfigTag(name string) string {
	return fmt.Sprintf(`
resource "zentral_tag" "test" {
  name = %[1]q
}

resource "zentral_osquery_pack" "test" {
  name = %[1]q
}

resource "zentral_osquery_query" "test" {
  name                     = %[1]q
  sql                      = "SELECT 1;"
  platforms                = ["darwin"]
  minimum_osquery_version  = "0.1.0"
  description              = "A query that always adds a tag"
  value                    = "Not much"
  tag_id                   = zentral_tag.test.id
  scheduling = {
    can_be_denylisted   = false,
    interval            = 161,
    log_removed_actions = false,
    pack_id             = zentral_osquery_pack.test.id
    shard               = 10
    snapshot_mode       = true
  }
}
`, name)
}
