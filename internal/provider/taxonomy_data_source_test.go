package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTaxonomyDataSource(t *testing.T) {
	rName := acctest.RandString(12)
	txResourceName := "zentral_taxonomy.test"
	ds1ResourceName := "data.zentral_taxonomy.by_id"
	ds2ResourceName := "data.zentral_taxonomy.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaxonomyDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", txResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", rName),
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", txResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", rName),
				),
			},
		},
	})
}

func testAccTaxonomyDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
resource "zentral_taxonomy" "test" {
  name = %q
}

data "zentral_taxonomy" "by_id" {
  id = zentral_taxonomy.test.id
}

data "zentral_taxonomy" "by_name" {
  name = zentral_taxonomy.test.name
}
`, rName)
}
