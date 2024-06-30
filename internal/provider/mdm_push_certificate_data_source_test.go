package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMPushCertificateDataSource(t *testing.T) {
	ds1ResourceName := "data.zentral_mdm_push_certificate.check1_by_name"
	ds2ResourceName := "data.zentral_mdm_push_certificate.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMPushCertificateDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttr(
						ds1ResourceName, "id", "4"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "provisioning_uid", "TF provider GitHub"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", "TF provider GitHub"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "topic"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "certificate"),
					// Read by ID
					resource.TestCheckResourceAttr(
						ds2ResourceName, "id", "4"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "provisioning_uid", "TF provider GitHub"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", "TF provider GitHub"),
					resource.TestCheckNoResourceAttr(
						ds2ResourceName, "topic"),
					resource.TestCheckNoResourceAttr(
						ds2ResourceName, "certificate"),
				),
			},
		},
	})
}

// TODO: hard coded values of a provisioned push certificate
// on the server used for the integration tests
func testAccMDMPushCertificateDataSourceConfig() string {
	return `
data "zentral_mdm_push_certificate" "check1_by_name" {
  name = "TF provider GitHub"
}

data "zentral_mdm_push_certificate" "check2_by_id" {
  id = 4
}
`
}
