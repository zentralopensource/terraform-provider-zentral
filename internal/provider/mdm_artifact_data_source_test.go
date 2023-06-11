package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMArtifactDataSource(t *testing.T) {
	c1Name := acctest.RandString(12)
	c2Name := acctest.RandString(12)
	c1ResourceName := "zentral_mdm_artifact.check1"
	c2ResourceName := "zentral_mdm_artifact.check2"
	ds1ResourceName := "data.zentral_mdm_artifact.check1_by_name"
	ds2ResourceName := "data.zentral_mdm_artifact.check2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMDMArtifactDataSourceConfig(c1Name, c2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", c1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", c1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "type", "Profile"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "channel", "Device"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "platforms.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						ds1ResourceName, "platforms.*", "macOS"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "install_during_setup_assistant", "false"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "auto_update", "true"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "reinstall_interval", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "reinstall_on_os_update", "No"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "requires.#", "0"),
					// Read by ID
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", c2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", c2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "type", "Profile"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "channel", "Device"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "platforms.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						ds2ResourceName, "platforms.*", "macOS"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "install_during_setup_assistant", "true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "auto_update", "false"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "reinstall_interval", "1"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "reinstall_on_os_update", "Minor"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "requires.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						ds2ResourceName, "requires.*", ds1ResourceName, "id"),
				),
			},
		},
	})
}

func testAccMDMArtifactDataSourceConfig(c1Name string, c2Name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_artifact" "check1" {
  name      = %[1]q
  type      = "Profile"
  channel   = "Device"
  platforms = ["macOS"]
}

resource "zentral_mdm_artifact" "check2" {
  name                           = %[2]q
  type                           = "Profile"
  channel                        = "Device"
  platforms                      = ["macOS"]
  install_during_setup_assistant = true
  auto_update                    = false
  reinstall_interval             = 1
  reinstall_on_os_update         = "Minor"
  requires                       = [zentral_mdm_artifact.check1.id]
}

data "zentral_mdm_artifact" "check1_by_name" {
  name = zentral_mdm_artifact.check1.name
}

data "zentral_mdm_artifact" "check2_by_id" {
  id = zentral_mdm_artifact.check2.id
}
`, c1Name, c2Name)
}
