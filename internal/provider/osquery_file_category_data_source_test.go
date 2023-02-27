package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOsqueryFileCategoryDataSource(t *testing.T) {
	fc1Name := acctest.RandString(12)
	fc2Name := acctest.RandString(12)
	fc1ResourceName := "zentral_osquery_file_category.check1"
	fc2ResourceName := "zentral_osquery_file_category.check2"
	ds1ResourceName := "data.zentral_osquery_file_category.check1_by_name"
	ds2ResourceName := "data.zentral_osquery_file_category.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOsqueryFileCategoryDataSourceConfig(fc1Name, fc2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", fc1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", fc1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "description", ""),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "file_paths.#", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "exclude_paths.#", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "file_paths_queries.#", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "access_monitoring", "false"),
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", fc2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", fc2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "description", fmt.Sprintf("%s description", fc2Name)),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "file_paths.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "file_paths.*", "/root/.ssh/%%"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "file_paths.*", "/home/%/.ssh/%%"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "exclude_paths.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "exclude_paths.*", "/home/not_to_monitor/.ssh/%%"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "file_paths_queries.#", "0"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "access_monitoring", "true"),
				),
			},
		},
	})
}

func testAccOsqueryFileCategoryDataSourceConfig(fc1Name string, fc2Name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_file_category" "check1" {
  name = %[1]q
}

resource "zentral_osquery_file_category" "check2" {
  name               = %[2]q
  description        = "%[2]s description"
  file_paths         = ["/root/.ssh/%%%%", "/home/%%/.ssh/%%%%"]
  exclude_paths      = ["/home/not_to_monitor/.ssh/%%%%"]
  file_paths_queries = []
  access_monitoring  = true
}

data "zentral_osquery_file_category" "check1_by_name" {
  name = zentral_osquery_file_category.check1.name
}

data "zentral_osquery_file_category" "check2_by_id" {
  id = zentral_osquery_file_category.check2.id
}
`, fc1Name, fc2Name)
}
