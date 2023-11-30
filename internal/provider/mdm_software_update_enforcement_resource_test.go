package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMSoftwareUpdateEnforcementResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_mdm_software_update_enforcement.test"
	tResourceName := "zentral_tag.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMSoftwareUpdateEnforcementResourceConfigLatest(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "details_url", ""),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "os_version", ""),
					resource.TestCheckResourceAttr(
						resourceName, "build_version", ""),
					resource.TestCheckNoResourceAttr(
						resourceName, "local_datetime"),
					resource.TestCheckResourceAttr(
						resourceName, "max_os_version", "15"),
					resource.TestCheckResourceAttr(
						resourceName, "delay_days", "14"),
					resource.TestCheckResourceAttr(
						resourceName, "local_time", "09:30:00"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccMDMSoftwareUpdateEnforcementResourceConfigOneTime(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "details_url", "https://www.example.com"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "tag_ids.*", tResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "os_version", "14.1"),
					resource.TestCheckResourceAttr(
						resourceName, "build_version", "23B74"),
					resource.TestCheckResourceAttr(
						resourceName, "local_datetime", "2023-11-05T09:30:00"),
					resource.TestCheckResourceAttr(
						resourceName, "max_os_version", ""),
					resource.TestCheckNoResourceAttr(
						resourceName, "delay_days"),
					resource.TestCheckNoResourceAttr(
						resourceName, "local_time"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMDMSoftwareUpdateEnforcementResourceConfigLatest(name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_software_update_enforcement" "test" {
  name           = %[1]q
  max_os_version = "15"
}
`, name)
}

func testAccMDMSoftwareUpdateEnforcementResourceConfigOneTime(name string) string {
	return fmt.Sprintf(`
resource "zentral_tag" "test" {
  name = %[1]q
}

resource "zentral_mdm_software_update_enforcement" "test" {
  name           = %[1]q
  details_url    = "https://www.example.com"
  tag_ids        = [zentral_tag.test.id]
  os_version     = "14.1"
  build_version  = "23B74"
  local_datetime = "2023-11-05T09:30:00"
}
`, name)
}
