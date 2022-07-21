package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetaBusinessUnitDataSource(t *testing.T) {
	r1Name := acctest.RandString(12)
	r2Name := acctest.RandString(12)
	mbu1ResourceName := "zentral_meta_business_unit.test_aee"
	mbu2ResourceName := "zentral_meta_business_unit.test_aed"
	ds1ResourceName := "data.zentral_meta_business_unit.by_id"
	ds2ResourceName := "data.zentral_meta_business_unit.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMetaBusinessUnitDataSourceConfig(r1Name, r2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", mbu1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", r1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "api_enrollment_enabled", "true"),
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", mbu2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", r2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "api_enrollment_enabled", "false"),
				),
			},
		},
	})
}

func testAccMetaBusinessUnitDataSourceConfig(r1Name string, r2Name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test_aee" {
  name = %q
}

data "zentral_meta_business_unit" "by_id" {
  id = zentral_meta_business_unit.test_aee.id
}

resource "zentral_meta_business_unit" "test_aed" {
  name = %q
  api_enrollment_enabled = false
}

data "zentral_meta_business_unit" "by_name" {
  name = zentral_meta_business_unit.test_aed.name
}
`, r1Name, r2Name)
}
