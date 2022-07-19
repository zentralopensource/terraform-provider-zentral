package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTagDataSource(t *testing.T) {
	taxName := acctest.RandString(12)
	t1Name := acctest.RandString(12)
	t2Name := acctest.RandString(12)
	t1Color := "0079bf" // default color
	t2Color := acctest.RandStringFromCharSet(6, "abcdef0123456789")
	tResourceName := "zentral_taxonomy.test"
	t1ResourceName := "zentral_tag.test1"
	t2ResourceName := "zentral_tag.test2"
	ds1ResourceName := "data.zentral_tag.test1_by_name"
	ds2ResourceName := "data.zentral_tag.test2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagDataSourceConfig(taxName, t1Name, t2Name, t2Color),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name, no taxonomy, default color
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", t1ResourceName, "id"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "taxonomy_id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", t1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "color", t1Color),
					// Read by ID, taxonomy, color
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", t2ResourceName, "id"),
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "taxonomy_id", tResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", t2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "color", t2Color),
				),
			},
		},
	})
}

func testAccTagDataSourceConfig(taxName string, t1Name string, t2Name string, t2Color string) string {
	return fmt.Sprintf(`
resource "zentral_taxonomy" "test" {
  name = %q
}

resource "zentral_tag" "test1" {
  name = %q
}

resource "zentral_tag" "test2" {
  taxonomy_id = zentral_taxonomy.test.id
  name        = %q
  color       = %q
}

data "zentral_tag" "test1_by_name" {
  name = zentral_tag.test1.name
}

data "zentral_tag" "test2_by_id" {
  id = zentral_tag.test2.id
}
`, taxName, t1Name, t2Name, t2Color)
}
