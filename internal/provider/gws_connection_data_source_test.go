package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGWSConnectionDataSource(t *testing.T) {
	ds1ResourceName := "data.zentral_gws_connection.check1_by_name"
	ds2ResourceName := "data.zentral_gws_connection.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGWSConnectionDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttr(
						ds1ResourceName, "id", "47aa8fc5-cb13-42c7-aed1-bb7950eec6f0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", "TF provider GitHub"),
					// Read by ID
					resource.TestCheckResourceAttr(
						ds2ResourceName, "id", "47aa8fc5-cb13-42c7-aed1-bb7950eec6f0"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", "TF provider GitHub"),
				),
			},
		},
	})
}

// TODO: hard coded values of a provisioned push certificate
// on the server used for the integration tests
func testAccGWSConnectionDataSourceConfig() string {
	return `
data "zentral_gws_connection" "check1_by_name" {
  name = "TF provider GitHub"
}

data "zentral_gws_connection" "check2_by_id" {
  id = "47aa8fc5-cb13-42c7-aed1-bb7950eec6f0"
}
`
}
