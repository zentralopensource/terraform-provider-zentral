package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTaxonomyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read by ID testing
			{
				Config: testAccTaxonomyDataSourceByIDConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zentral_taxonomy.test", "name", "default"),
				),
			},
			// Read by name testing
			{
				Config: testAccTaxonomyDataSourceByNameConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zentral_taxonomy.test", "id", "1"),
				),
			},
		},
	})
}

const testAccTaxonomyDataSourceByIDConfig = `
data "zentral_taxonomy" "test" {
  id = 1
}
`

const testAccTaxonomyDataSourceByNameConfig = `
data "zentral_taxonomy" "test" {
  name = "default"
}
`
