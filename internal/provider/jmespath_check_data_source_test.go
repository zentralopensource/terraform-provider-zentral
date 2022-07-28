package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJMESPathCheckDataSource(t *testing.T) {
	tagName := acctest.RandString(12)
	c1Name := acctest.RandString(12)
	c2Name := acctest.RandString(12)
	tagResourceName := "zentral_tag.test"
	c1ResourceName := "zentral_jmespath_check.check1"
	c2ResourceName := "zentral_jmespath_check.check2"
	ds1ResourceName := "data.zentral_jmespath_check.check1_by_name"
	ds2ResourceName := "data.zentral_jmespath_check.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJMESPathCheckDataSourceConfig(tagName, c1Name, c2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name, platforms, tag
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", c1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", c1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "description", "desc1"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "source_name", "osquery1"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "platforms.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						ds1ResourceName, "platforms.*", "MACOS"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						ds1ResourceName, "tag_ids.*", tagResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "jmespath_expression", "ok1"),
					// Read by ID, no platforms, no tags
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", c2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", c2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "description", "desc2"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "source_name", "osquery2"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "platforms.#", "0"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "jmespath_expression", "ok2"),
				),
			},
		},
	})
}

func testAccJMESPathCheckDataSourceConfig(tagName string, c1Name string, c2Name string) string {
	return fmt.Sprintf(`
resource "zentral_tag" "test" {
  name = %q
}

resource "zentral_jmespath_check" "check1" {
  name = %q
  description = "desc1"
  source_name = "osquery1"
  platforms = ["MACOS"]
  tag_ids = [zentral_tag.test.id]
  jmespath_expression = "ok1"
}

resource "zentral_jmespath_check" "check2" {
  name = %q
  description = "desc2"
  source_name = "osquery2"
  jmespath_expression = "ok2"
}

data "zentral_jmespath_check" "check1_by_name" {
  name = zentral_jmespath_check.check1.name
}

data "zentral_jmespath_check" "check2_by_id" {
  id = zentral_jmespath_check.check2.id
}
`, tagName, c1Name, c2Name)
}
