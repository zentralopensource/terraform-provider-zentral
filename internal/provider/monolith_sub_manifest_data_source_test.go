package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonolithSubManifestDataSource(t *testing.T) {
	sm1Name := acctest.RandString(12)
	sm2Name := acctest.RandString(12)
	sm1ResourceName := "zentral_monolith_sub_manifest.test1"
	sm2ResourceName := "zentral_monolith_sub_manifest.test2"
	mbuResourceName := "zentral_meta_business_unit.test"
	ds1ResourceName := "data.zentral_monolith_sub_manifest.test_by_id"
	ds2ResourceName := "data.zentral_monolith_sub_manifest.test_by_name"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonolithSubManifestDataSourceConfig(sm1Name, sm2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by id
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", sm1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", sm1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "description", ""),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "meta_business_unit_id"),
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", sm2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", sm2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "description", "Description"),
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "meta_business_unit_id", mbuResourceName, "id"),
				),
			},
		},
	})
}

func testAccMonolithSubManifestDataSourceConfig(sm1Name string, sm2Name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[2]q
}

resource "zentral_monolith_sub_manifest" "test1" {
  name = %[1]q
}

resource "zentral_monolith_sub_manifest" "test2" {
  name                  = %[2]q
  description           = "Description"
  meta_business_unit_id = zentral_meta_business_unit.test.id
}

data "zentral_monolith_sub_manifest" "test_by_id" {
  id = zentral_monolith_sub_manifest.test1.id
}

data "zentral_monolith_sub_manifest" "test_by_name" {
  name = zentral_monolith_sub_manifest.test2.name
}
`, sm1Name, sm2Name)
}
