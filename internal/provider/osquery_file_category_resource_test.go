package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOsqueryFileCategoryResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_osquery_file_category.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccOsqueryFileCategoryResourceConfigBare(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "file_paths.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "exclude_paths.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "file_paths_queries.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "access_monitoring", "false"),
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
				Config: testAccOsqueryFileCategoryResourceConfigFull(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "description", fmt.Sprintf("%s description", secondName)),
					resource.TestCheckResourceAttr(
						resourceName, "file_paths.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "file_paths.*", "/root/.ssh/%%"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "file_paths.*", "/home/%/.ssh/%%"),
					resource.TestCheckResourceAttr(
						resourceName, "exclude_paths.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "exclude_paths.*", "/home/not_to_monitor/.ssh/%%"),
					resource.TestCheckResourceAttr(
						resourceName, "file_paths_queries.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "access_monitoring", "true"),
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

func testAccOsqueryFileCategoryResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_file_category" "test" {
  name = %[1]q
}
`, name)
}

func testAccOsqueryFileCategoryResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_osquery_file_category" "test" {
  name               = %[1]q
  description        = "%[1]s description"
  file_paths         = ["/root/.ssh/%%%%", "/home/%%/.ssh/%%%%"]
  exclude_paths      = ["/home/not_to_monitor/.ssh/%%%%"]
  file_paths_queries = []
  access_monitoring  = true
}
`, name)
}
