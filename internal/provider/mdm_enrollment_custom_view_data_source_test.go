package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMEnrollmentCustomViewDataSource(t *testing.T) {
	bidName := acctest.RandString(12)
	bnName := acctest.RandString(12)
	bidResourceName := "zentral_mdm_enrollment_custom_view.by_id"
	bnResourceName := "zentral_mdm_enrollment_custom_view.by_name"
	ds1ResourceName := "data.zentral_mdm_enrollment_custom_view.by_id"
	ds2ResourceName := "data.zentral_mdm_enrollment_custom_view.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMEnrollmentCustomViewDataSourceConfig(bidName, bnName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by id
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", bidResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", bidName),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "html", "<html/>"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "requires_authentication", "false"),
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", bnResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", bnName),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "html", "<html/>"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "requires_authentication", "false"),
				),
			},
		},
	})
}

func testAccMDMEnrollmentCustomViewDataSourceConfig(bidName string, bnName string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_enrollment_custom_view" "by_id" {
	name                  			= %[1]q
	description						= %[1]q
	html							= "<html/>"
}

resource "zentral_mdm_enrollment_custom_view" "by_name" {
	name                  			= %[2]q
	description						= %[2]q
	html							= "<html/>"
}

data "zentral_mdm_enrollment_custom_view" "by_id" {
  id = zentral_mdm_enrollment_custom_view.by_id.id
}

data "zentral_mdm_enrollment_custom_view" "by_name" {
  name = zentral_mdm_enrollment_custom_view.by_name.name
}
`, bidName, bnName)
}
