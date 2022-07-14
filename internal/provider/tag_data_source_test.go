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
					resource.TestCheckResourceAttr("data.zentral_tag.test", "name", "default"),
					resource.TestCheckResourceAttr("data.zentral_tag.test", "color", "c377e0"),
				),
			},
			// Read by name testing
			{
				Config: testAccTagDataSourceByNameConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zentral_tag.test", "id", "3"),
					resource.TestCheckResourceAttr("data.zentral_tag.test", "color", "c377e0"),
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

const testAccTagDataSourceByNameConfig = `
data "zentral_tag" "test" {
  name = "default"
}
`
