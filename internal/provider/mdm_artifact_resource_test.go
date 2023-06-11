package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMArtifactResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	raName := acctest.RandString(12)
	resourceName := "zentral_mdm_artifact.test"
	raResourceName := "zentral_mdm_artifact.required"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMArtifactResourceConfigBare(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "type", "Profile"),
					resource.TestCheckResourceAttr(
						resourceName, "channel", "Device"),
					resource.TestCheckResourceAttr(
						resourceName, "platforms.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "platforms.*", "iOS"),
					resource.TestCheckResourceAttr(
						resourceName, "install_during_setup_assistant", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_update", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "reinstall_interval", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "reinstall_on_os_update", "No"),
					resource.TestCheckResourceAttr(
						resourceName, "requires.#", "0"),
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
				Config: testAccMDMArtifactResourceConfigFull(secondName, raName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "type", "Profile"),
					resource.TestCheckResourceAttr(
						resourceName, "channel", "Device"),
					resource.TestCheckResourceAttr(
						resourceName, "platforms.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "platforms.*", "macOS"),
					resource.TestCheckResourceAttr(
						resourceName, "install_during_setup_assistant", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_update", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "reinstall_interval", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "reinstall_on_os_update", "Minor"),
					resource.TestCheckResourceAttr(
						resourceName, "requires.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "requires.*", raResourceName, "id"),
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

func testAccMDMArtifactResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_artifact" "test" {
  name      = %[1]q
  type      = "Profile"
  channel   = "Device"
  platforms = ["iOS"]
}
`, name)
}

func testAccMDMArtifactResourceConfigFull(name string, raName string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_artifact" "required" {
  name      = %[2]q
  type      = "Profile"
  channel   = "Device"
  platforms = ["macOS"]
}

resource "zentral_mdm_artifact" "test" {
  name                           = %[1]q
  type                           = "Profile"
  channel                        = "Device"
  platforms                      = ["macOS"]
  install_during_setup_assistant = true
  auto_update                    = false
  reinstall_interval             = 1
  reinstall_on_os_update         = "Minor"
  requires                       = [zentral_mdm_artifact.required.id]
}
`, name, raName)
}
