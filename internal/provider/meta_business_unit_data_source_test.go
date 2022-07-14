package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetaBusinessUnitDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read by ID testing
			{
				Config: testAccMetaBusinessUnitDataSourceByIDConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zentral_meta_business_unit.test", "name", "default"),
				),
			},
			// Read by name testing
			{
				Config: testAccMetaBusinessUnitDataSourceByNameConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zentral_meta_business_unit.test", "id", "1"),
				),
			},
		},
	})
}

const testAccMetaBusinessUnitDataSourceByIDConfig = `
data "zentral_meta_business_unit" "test" {
  id = 1
}
`

const testAccMetaBusinessUnitDataSourceByNameConfig = `
data "zentral_meta_business_unit" "test" {
  name = "default"
}
`
