package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSantaEnrollmentDataSource(t *testing.T) {
	name := acctest.RandString(12)
	tagName := acctest.RandString(12)
	resourceName := "zentral_santa_enrollment.test"
	dataSourceName := "data.zentral_santa_enrollment.test"
	cfgResourceName := "zentral_santa_configuration.test"
	mbuResourceName := "zentral_meta_business_unit.test"
	tagResourceName := "zentral_tag.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSantaEnrollmentDataSourceConfig(name, tagName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name, platforms, tag
					resource.TestCheckResourceAttrPair(
						dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(
						dataSourceName, "configuration_id", cfgResourceName, "id"),
					resource.TestCheckResourceAttr(
						dataSourceName, "version", "1"),
					resource.TestCheckResourceAttrPair(
						dataSourceName, "meta_business_unit_id", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						dataSourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						dataSourceName, "tag_ids.*", tagResourceName, "id"),
					resource.TestCheckResourceAttr(
						dataSourceName, "serial_numbers.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						dataSourceName, "serial_numbers.*", "un"),
					resource.TestCheckTypeSetElemAttr(
						dataSourceName, "serial_numbers.*", "deux"),
					resource.TestCheckResourceAttr(
						dataSourceName, "udids.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						dataSourceName, "udids.*", "trois"),
					resource.TestCheckTypeSetElemAttr(
						dataSourceName, "udids.*", "quatre"),
					resource.TestCheckResourceAttr(
						dataSourceName, "quota", "5"),
				),
			},
		},
	})
}

func testAccSantaEnrollmentDataSourceConfig(name string, tagName string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_santa_configuration" "test" {
  name = %[1]q
}

resource "zentral_taxonomy" "test" {
  name = %[1]q
}

resource "zentral_tag" "test" {
  taxonomy_id = zentral_taxonomy.test.id
  name        = %[2]q
}

resource "zentral_santa_enrollment" "test" {
  configuration_id      = zentral_santa_configuration.test.id
  meta_business_unit_id = zentral_meta_business_unit.test.id
  tag_ids               = [zentral_tag.test.id]
  serial_numbers        = ["un", "deux"]
  udids                 = ["trois", "quatre"]
  quota                 = 5
}

data "zentral_santa_enrollment" "test" {
  id = zentral_santa_enrollment.test.id
}
`, name, tagName)
}
