package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMunkiScriptCheckResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_munki_script_check.test"
	tagResourceName := "zentral_tag.test"
	excludedTagResourceName := "zentral_tag.test-excluded"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMunkiScriptCheckResourceConfigBare(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "type", "ZSH_STR"),
					resource.TestCheckResourceAttr(
						resourceName, "source", "echo test"),
					resource.TestCheckResourceAttr(
						resourceName, "expected_result", "test"),
					resource.TestCheckResourceAttr(
						resourceName, "arch_amd64", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "arch_arm64", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "min_os_version", ""),
					resource.TestCheckResourceAttr(
						resourceName, "max_os_version", ""),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "1"),
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
				Config: testAccMunkiScriptCheckResourceConfigFull(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "Description"),
					resource.TestCheckResourceAttr(
						resourceName, "type", "ZSH_INT"),
					resource.TestCheckResourceAttr(
						resourceName, "source", "echo 10"),
					resource.TestCheckResourceAttr(
						resourceName, "expected_result", "10"),
					resource.TestCheckResourceAttr(
						resourceName, "arch_amd64", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "arch_arm64", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "min_os_version", "14"),
					resource.TestCheckResourceAttr(
						resourceName, "max_os_version", "15"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "tag_ids.*", tagResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "excluded_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "excluded_tag_ids.*", excludedTagResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "2"),
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

func testAccMunkiScriptCheckResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_munki_script_check" "test" {
  name            = %[1]q
  source          = "echo test"
  expected_result = "test"
}
`, name)
}

func testAccMunkiScriptCheckResourceConfigFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_tag" "test" {
  name = %[1]q
}

resource "zentral_tag" "test-excluded" {
  name = "%[1]s-excluded"
}

resource "zentral_munki_script_check" "test" {
  name             = %[1]q
  description      = "Description"
  type             = "ZSH_INT"
  source           = "echo 10"
  expected_result  = "10"
  arch_amd64       = false
  arch_arm64       = true
  min_os_version   = "14"
  max_os_version   = "15"
  tag_ids          = [zentral_tag.test.id]
  excluded_tag_ids = [zentral_tag.test-excluded.id]
}
`, name)
}
