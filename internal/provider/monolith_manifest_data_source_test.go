package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonolithManifestDataSource(t *testing.T) {
	m1Name := acctest.RandString(12)
	m2Name := acctest.RandString(12)
	m1ResourceName := "zentral_monolith_manifest.test1"
	m2ResourceName := "zentral_monolith_manifest.test2"
	mbuResourceName := "zentral_meta_business_unit.test"
	ds1ResourceName := "data.zentral_monolith_manifest.test_by_id"
	ds2ResourceName := "data.zentral_monolith_manifest.test_by_name"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonolithManifestDataSourceConfig(m1Name, m2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by id
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", m1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", m1Name),
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "meta_business_unit_id", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "version", "1"),
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", m2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", m2Name),
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "meta_business_unit_id", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "version", "1"),
				),
			},
		},
	})
}

func testAccMonolithManifestDataSourceConfig(m1Name string, m2Name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_monolith_manifest" "test1" {
  name                  = %[1]q
  meta_business_unit_id = zentral_meta_business_unit.test.id
}

resource "zentral_monolith_manifest" "test2" {
  name                  = %[2]q
  meta_business_unit_id = zentral_meta_business_unit.test.id
}

data "zentral_monolith_manifest" "test_by_id" {
  id = zentral_monolith_manifest.test1.id
}

data "zentral_monolith_manifest" "test_by_name" {
  name = zentral_monolith_manifest.test2.name
}
`, m1Name, m2Name)
}
