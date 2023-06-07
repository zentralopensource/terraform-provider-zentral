package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMBlueprintDataSource(t *testing.T) {
	c1Name := acctest.RandString(12)
	c2Name := acctest.RandString(12)
	c1ResourceName := "zentral_mdm_blueprint.check1"
	c2ResourceName := "zentral_mdm_blueprint.check2"
	ds1ResourceName := "data.zentral_mdm_blueprint.check1_by_name"
	ds2ResourceName := "data.zentral_mdm_blueprint.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMBlueprintDataSourceConfig(c1Name, c2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", c1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", c1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "inventory_interval", "86400"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "collect_apps", "NO"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "collect_certificates", "NO"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "collect_profiles", "NO"),
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", c2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", c2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "inventory_interval", "77777"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "collect_apps", "MANAGED_ONLY"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "collect_certificates", "ALL"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "collect_profiles", "MANAGED_ONLY"),
				),
			},
		},
	})
}

func testAccMDMBlueprintDataSourceConfig(c1Name string, c2Name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_blueprint" "check1" {
  name = %[1]q
}

resource "zentral_mdm_blueprint" "check2" {
  name                 = %[2]q
  inventory_interval   = 77777
  collect_apps         = "MANAGED_ONLY"
  collect_certificates = "ALL"
  collect_profiles     = "MANAGED_ONLY"
}

data "zentral_mdm_blueprint" "check1_by_name" {
  name = zentral_mdm_blueprint.check1.name
}

data "zentral_mdm_blueprint" "check2_by_id" {
  id = zentral_mdm_blueprint.check2.id
}
`, c1Name, c2Name)
}
