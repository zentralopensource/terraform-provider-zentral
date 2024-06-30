package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMSCEPConfigDataSource(t *testing.T) {
	ds1ResourceName := "data.zentral_mdm_scep_config.check1_by_name"
	ds2ResourceName := "data.zentral_mdm_scep_config.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMSCEPConfigDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttr(
						ds1ResourceName, "id", "2"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "provisioning_uid", "TF provider GitHub"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", "TF provider GitHub"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "url", "https://github.com/zentralopensource/terraform-provider-zentral"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "key_usage", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "key_is_extractable", "false"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "key_size", "2048"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "allow_all_apps_access", "false"),
					// Read by ID
					resource.TestCheckResourceAttr(
						ds2ResourceName, "id", "2"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "provisioning_uid", "TF provider GitHub"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", "TF provider GitHub"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "url", "https://github.com/zentralopensource/terraform-provider-zentral"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "key_usage", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "key_is_extractable", "false"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "key_size", "2048"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "allow_all_apps_access", "false"),
				),
			},
		},
	})
}

// TODO: hard coded values of a provisioned push certificate
// on the server used for the integration tests
func testAccMDMSCEPConfigDataSourceConfig() string {
	return `
data "zentral_mdm_scep_config" "check1_by_name" {
  name = "TF provider GitHub"
}

data "zentral_mdm_scep_config" "check2_by_id" {
  id = 2
}
`
}
