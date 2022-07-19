package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetaBusinessUnitDataSource(t *testing.T) {
	rName := acctest.RandString(12)
	mbuResourceName := "zentral_meta_business_unit.test"
	ds1ResourceName := "data.zentral_meta_business_unit.by_id"
	ds2ResourceName := "data.zentral_meta_business_unit.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMetaBusinessUnitDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", rName),
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", rName),
				),
			},
		},
	})
}

func testAccMetaBusinessUnitDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %q
}

data "zentral_meta_business_unit" "by_id" {
  id = zentral_meta_business_unit.test.id
}

data "zentral_meta_business_unit" "by_name" {
  name = zentral_meta_business_unit.test.name
}
`, rName)
}
