package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMDEPVirtualServerDataSource(t *testing.T) {
	ds1ResourceName := "data.zentral_mdm_dep_virtual_server.check1_by_name"
	ds2ResourceName := "data.zentral_mdm_dep_virtual_server.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMDEPVirtualServerDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", "TF provider GitHub"),
					// Read by ID
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", "TF provider GitHub"),
				),
			},
		},
	})
}

// TODO: hard coded values of a virtual server
// on the server used for the integration tests
func testAccMDMDEPVirtualServerDataSourceConfig() string {
	return `
data "zentral_mdm_dep_virtual_server" "check1_by_name" {
  name = "TF provider GitHub"
}

data "zentral_mdm_dep_virtual_server" "check2_by_id" {
   id = data.zentral_mdm_dep_virtual_server.check1_by_name.id
}
`
}
