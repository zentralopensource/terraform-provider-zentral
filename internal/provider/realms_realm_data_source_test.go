package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRealmsRealmDataSource(t *testing.T) {
	ds1ResourceName := "data.zentral_realm.check1_by_name"
	ds2ResourceName := "data.zentral_realm.check2_by_uuid"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRealmsRealmDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttr(
						ds1ResourceName, "id", "73ce6db8-ebcf-4ba2-b624-0df721512fae"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "uuid", "73ce6db8-ebcf-4ba2-b624-0df721512fae"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", "TF provider GitHub"),
					// Read by UUID
					resource.TestCheckResourceAttr(
						ds2ResourceName, "id", "73ce6db8-ebcf-4ba2-b624-0df721512fae"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "uuid", "73ce6db8-ebcf-4ba2-b624-0df721512fae"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", "TF provider GitHub"),
				),
			},
		},
	})
}

// TODO: hard coded values of a realm
// on the server used for the integration tests
func testAccRealmsRealmDataSourceConfig() string {
	return `
data "zentral_realm" "check1_by_name" {
  name = "TF provider GitHub"
}

data "zentral_realm" "check2_by_uuid" {
  uuid = "73ce6db8-ebcf-4ba2-b624-0df721512fae"
}
`
}
