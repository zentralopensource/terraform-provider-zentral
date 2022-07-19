package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTagResource(t *testing.T) {
	firstName := acctest.RandString(12)
	firstColor := acctest.RandStringFromCharSet(6, "abcdef0123456789")
	txName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	secondColor := acctest.RandStringFromCharSet(6, "abcdef0123456789")
	txResourceName := "zentral_taxonomy.test"
	resourceName := "zentral_tag.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTagResourceConfig(txName, firstName, firstColor),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(
						resourceName, "taxonomy_id"),
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "color", firstColor),
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
				Config: testAccTagResourceConfig(txName, secondName, secondColor),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(
						resourceName, "taxonomy_id"),
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "color", secondColor),
				),
			},
			// Update and Read with taxonomy
			{
				Config: testAccTagWithTaxonomyResourceConfig(txName, secondName, secondColor),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resourceName, "taxonomy_id", txResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "color", secondColor),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read without taxonomy
			{
				Config: testAccTagResourceConfig(txName, secondName, secondColor),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(
						resourceName, "taxonomy_id"),
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "color", secondColor),
				),
			},
		},
	})
}

func testAccTagResourceConfig(txName string, name string, color string) string {
	return fmt.Sprintf(`
resource "zentral_taxonomy" "test" {
  name = %[1]q
}

resource "zentral_tag" "test" {
  name = %[2]q
  color = %[3]q
}
`, txName, name, color)
}

func testAccTagWithTaxonomyResourceConfig(txName string, name string, color string) string {
	return fmt.Sprintf(`
resource "zentral_taxonomy" "test" {
  name = %[1]q
}

resource "zentral_tag" "test" {
  taxonomy_id = zentral_taxonomy.test.id
  name = %[2]q
  color = %[3]q
}
`, txName, name, color)
}
