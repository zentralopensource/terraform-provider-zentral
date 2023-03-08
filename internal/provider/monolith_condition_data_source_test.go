package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonolithConditionDataSource(t *testing.T) {
	c1Name := acctest.RandString(12)
	c2Name := acctest.RandString(12)
	c1ResourceName := "zentral_monolith_condition.test1"
	c2ResourceName := "zentral_monolith_condition.test2"
	ds1ResourceName := "data.zentral_monolith_condition.test_by_id"
	ds2ResourceName := "data.zentral_monolith_condition.test_by_name"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonolithConditionDataSourceConfig(c1Name, c2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by id
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", c1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", c1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "predicate", "machine_type == \"desktop\""),
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", c2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", c2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "predicate", "machine_type == \"laptop\""),
				),
			},
		},
	})
}

func testAccMonolithConditionDataSourceConfig(c1Name string, c2Name string) string {
	return fmt.Sprintf(`
resource "zentral_monolith_condition" "test1" {
  name      = %[1]q
  predicate = "machine_type == \"desktop\""

}

resource "zentral_monolith_condition" "test2" {
  name      = %[2]q
  predicate = "machine_type == \"laptop\""
}

data "zentral_monolith_condition" "test_by_id" {
  id = zentral_monolith_condition.test1.id
}

data "zentral_monolith_condition" "test_by_name" {
  name = zentral_monolith_condition.test2.name
}
`, c1Name, c2Name)
}
