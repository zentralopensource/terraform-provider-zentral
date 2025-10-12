package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMLocationDataSource(t *testing.T) {
	ds1ResourceName := "data.zentral_mdm_location.check1_by_id"
	ds2ResourceName := "data.zentral_mdm_location.check2_by_mdm_info_id"
	ds3ResourceName := "data.zentral_mdm_location.check3_by_name"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMLocationDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by ID
					resource.TestCheckResourceAttr(
						ds1ResourceName, "id", "3"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", "Terraform Provider CI/CD"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "organization_name", "Zentral Pro Services GmbH"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "mdm_info_id", "6e3d30f0-4812-443d-9e0f-2e4b80709fef"),
					// Read by MDM info ID
					resource.TestCheckResourceAttr(
						ds2ResourceName, "id", "3"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", "Terraform Provider CI/CD"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "organization_name", "Zentral Pro Services GmbH"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "mdm_info_id", "6e3d30f0-4812-443d-9e0f-2e4b80709fef"),
					// Read by Name
					resource.TestCheckResourceAttr(
						ds3ResourceName, "id", "3"),
					resource.TestCheckResourceAttr(
						ds3ResourceName, "name", "Terraform Provider CI/CD"),
					resource.TestCheckResourceAttr(
						ds3ResourceName, "organization_name", "Zentral Pro Services GmbH"),
					resource.TestCheckResourceAttr(
						ds3ResourceName, "mdm_info_id", "6e3d30f0-4812-443d-9e0f-2e4b80709fef"),
				),
			},
		},
	})
}

// TODO: hard coded values of a provisioned location
// on the server used for the integration tests
func testAccMDMLocationDataSourceConfig() string {
	return `
data "zentral_mdm_location" "check1_by_id" {
  id = 3
}

data "zentral_mdm_location" "check2_by_mdm_info_id" {
  mdm_info_id = "6e3d30f0-4812-443d-9e0f-2e4b80709fef"
}

data "zentral_mdm_location" "check3_by_name" {
  name = "Terraform Provider CI/CD"
}
`
}
