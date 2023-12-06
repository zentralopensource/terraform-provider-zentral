package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMSoftwareUpdateEnforcementDataSource(t *testing.T) {
	c1Name := acctest.RandString(12)
	c2Name := acctest.RandString(12)
	c1ResourceName := "zentral_mdm_software_update_enforcement.check1"
	c2ResourceName := "zentral_mdm_software_update_enforcement.check2"
	ds1ResourceName := "data.zentral_mdm_software_update_enforcement.check1_by_name"
	ds2ResourceName := "data.zentral_mdm_software_update_enforcement.check2_by_id"
	tResourceName := "zentral_tag.check2"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMSoftwareUpdateEnforcementDataSourceConfig(c1Name, c2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", c1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "details_url", ""),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "platforms.#", "2"),
					resource.TestCheckTypeSetElemAttr(
						ds1ResourceName, "platforms.*", "iOS"),
					resource.TestCheckTypeSetElemAttr(
						ds1ResourceName, "platforms.*", "iPadOS"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "os_version", ""),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "build_version", ""),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "local_datetime"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "max_os_version", "15"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "delay_days", "3"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "local_time", "04:05:06"),
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", c2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "details_url", "https://www.example.com"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "platforms.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "platforms.*", "macOS"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						ds2ResourceName, "tag_ids.*", tResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "os_version", "14.1"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "build_version", "23B74"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "local_datetime", "2023-11-05T09:30:00"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "max_os_version", ""),
					resource.TestCheckNoResourceAttr(
						ds2ResourceName, "delay_days"),
					resource.TestCheckNoResourceAttr(
						ds2ResourceName, "local_time"),
				),
			},
		},
	})
}

func testAccMDMSoftwareUpdateEnforcementDataSourceConfig(c1Name string, c2Name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_software_update_enforcement" "check1" {
  name           = %[1]q
  platforms      = ["iPadOS", "iOS"]
  max_os_version = "15"
  delay_days     = 3
  local_time     = "04:05:06"
}

resource "zentral_tag" "check2" {
  name = %[2]q
}

resource "zentral_mdm_software_update_enforcement" "check2" {
  name           = %[2]q
  details_url    = "https://www.example.com"
  platforms      = ["macOS"]
  tag_ids        = [zentral_tag.check2.id]
  os_version     = "14.1"
  build_version  = "23B74"
  local_datetime = "2023-11-05T09:30:00"
}

data "zentral_mdm_software_update_enforcement" "check1_by_name" {
  name = zentral_mdm_software_update_enforcement.check1.name
}

data "zentral_mdm_software_update_enforcement" "check2_by_id" {
  id = zentral_mdm_software_update_enforcement.check2.id
}
`, c1Name, c2Name)
}
