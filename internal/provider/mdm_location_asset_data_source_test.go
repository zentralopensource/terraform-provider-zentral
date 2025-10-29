package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMLocationAssetDataSource(t *testing.T) {
	ds1ResourceName := "data.zentral_mdm_location_asset.check1"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMLocationAssetDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read
					resource.TestCheckResourceAttr(
						ds1ResourceName, "id", "19"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "location_id", "3"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "asset_id", "1"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "adam_id", "803453959"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "pricing_param", "STDQ"),
				),
			},
		},
	})
}

// TODO: hard coded values of a provisioned location asset
// on the server used for the integration tests
func testAccMDMLocationAssetDataSourceConfig() string {
	return `
data "zentral_mdm_location" "check1" {
  name = "Terraform Provider CI/CD"
}

data "zentral_mdm_location_asset" "check1" {
  location_id   = data.zentral_mdm_location.check1.id
  adam_id       = "803453959"
  pricing_param = "STDQ"
}
`
}
