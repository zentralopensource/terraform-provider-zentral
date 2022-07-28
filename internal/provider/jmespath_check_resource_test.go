package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJMESPathCheckResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	t1Name := acctest.RandString(12)
	t2Name := acctest.RandString(12)
	resourceName := "zentral_jmespath_check.test"
	t1Resource := "zentral_tag.test1"
	t2Resource := "zentral_tag.test2"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccJMESPathCheckResourceConfigBare(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "source_name", "osquery"),
					resource.TestCheckResourceAttr(
						resourceName, "platforms.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "jmespath_expression", "ok"),
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
				Config: testAccJMESPathCheckResourceConfigFull(secondName, t1Name, t2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "desc"),
					resource.TestCheckResourceAttr(
						resourceName, "source_name", "osquery"),
					resource.TestCheckResourceAttr(
						resourceName, "platforms.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "platforms.*", "MACOS"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "platforms.*", "LINUX"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "tag_ids.*", t1Resource, "id"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "tag_ids.*", t2Resource, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "jmespath_expression", "ok"),
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

func testAccJMESPathCheckResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_jmespath_check" "test" {
  name = %[1]q
  source_name = "osquery"
  jmespath_expression = "ok"
}
`, name)
}

func testAccJMESPathCheckResourceConfigFull(name string, t1Name string, t2Name string) string {
	return fmt.Sprintf(`
resource "zentral_tag" "test1" {
  name = %[2]q
}

resource "zentral_tag" "test2" {
  name = %[3]q
}

resource "zentral_jmespath_check" "test" {
  name = %[1]q
  description = "desc"
  source_name = "osquery"
  platforms = ["MACOS", "LINUX"]
  tag_ids = [zentral_tag.test1.id, zentral_tag.test2.id]
  jmespath_expression = "ok"
}
`, name, t1Name, t2Name)
}
