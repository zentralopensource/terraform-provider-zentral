package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTagResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTagResourceConfig("one", "110000"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("zentral_tag.test", "taxonomy_id"),
					resource.TestCheckResourceAttr("zentral_tag.test", "name", "one"),
					resource.TestCheckResourceAttr("zentral_tag.test", "color", "110000"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "zentral_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccTagResourceConfig("two", "001100"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("zentral_tag.test", "taxonomy_id"),
					resource.TestCheckResourceAttr("zentral_tag.test", "name", "two"),
					resource.TestCheckResourceAttr("zentral_tag.test", "color", "001100"),
				),
			},
			// Update and Read with taxonomy testing
			{
				Config: testAccTagWithTaxonomyResourceConfig("two", "001100"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zentral_tag.test", "taxonomy_id", "1"),
					resource.TestCheckResourceAttr("zentral_tag.test", "name", "two"),
					resource.TestCheckResourceAttr("zentral_tag.test", "color", "001100"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "zentral_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read without taxonomy testing
			{
				Config: testAccTagResourceConfig("two", "001100"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("zentral_tag.test", "taxonomy_id"),
					resource.TestCheckResourceAttr("zentral_tag.test", "name", "two"),
					resource.TestCheckResourceAttr("zentral_tag.test", "color", "001100"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTagResourceConfig(name string, color string) string {
	return fmt.Sprintf(`
resource "zentral_tag" "test" {
  name = %[1]q
  color = %[2]q
}
`, name, color)
}

func testAccTagWithTaxonomyResourceConfig(name string, color string) string {
	return fmt.Sprintf(`
resource "zentral_tag" "test" {
  taxonomy_id = 1
  name = %[1]q
  color = %[2]q
}
`, name, color)
}
