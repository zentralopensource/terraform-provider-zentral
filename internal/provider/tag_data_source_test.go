package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTagDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read by ID testing
			{
				Config: testAccTagDataSourceByIDConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("data.zentral_tag.test", "taxonomy_id"),
					resource.TestCheckResourceAttr("data.zentral_tag.test", "name", "default"),
					resource.TestCheckResourceAttr("data.zentral_tag.test", "color", "c377e0"),
				),
			},
			// Read by name testing
			{
				Config: testAccTagDataSourceByNameConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("data.zentral_tag.test", "taxonomy_id"),
					resource.TestCheckResourceAttr("data.zentral_tag.test", "id", "3"),
					resource.TestCheckResourceAttr("data.zentral_tag.test", "color", "c377e0"),
				),
			},
			// Read by ID with taxonomy testing
			{
				Config: testAccTagDataSourceByIDWithTaxonomyConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zentral_tag.test", "taxonomy_id", "1"),
					resource.TestCheckResourceAttr("data.zentral_tag.test", "name", "default-with-taxonomy"),
					resource.TestCheckResourceAttr("data.zentral_tag.test", "color", "eb5a46"),
				),
			},
			// Read by name with taxonomy testing
			{
				Config: testAccTagDataSourceByNameWithTaxonomyConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zentral_tag.test", "id", "44"),
					resource.TestCheckResourceAttr("data.zentral_tag.test", "taxonomy_id", "1"),
					resource.TestCheckResourceAttr("data.zentral_tag.test", "color", "eb5a46"),
				),
			},
		},
	})
}

const testAccTagDataSourceByIDConfig = `
data "zentral_tag" "test" {
  id = 3
}
`

const testAccTagDataSourceByIDWithTaxonomyConfig = `
data "zentral_tag" "test" {
	id = 44
  }
`

const testAccTagDataSourceByNameConfig = `
data "zentral_tag" "test" {
  name = "default"
}
`

const testAccTagDataSourceByNameWithTaxonomyConfig = `
data "zentral_tag" "test" {
  name = "default-with-taxonomy"
}
`
