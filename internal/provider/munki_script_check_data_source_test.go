package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMunkiScriptCheckDataSource(t *testing.T) {
	sc1Name := acctest.RandString(12)
	sc2Name := acctest.RandString(12)
	sc1ResourceName := "zentral_munki_script_check.check1"
	sc2ResourceName := "zentral_munki_script_check.check2"
	ds1ResourceName := "data.zentral_munki_script_check.check1_by_name"
	ds2ResourceName := "data.zentral_munki_script_check.check2_by_id"
	tagResourceName := "zentral_tag.test"
	excludedTagResourceName := "zentral_tag.test-excluded"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMunkiScriptCheckDataSourceConfig(sc1Name, sc2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name, platforms, tag
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", sc1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", sc1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "description", ""),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "type", "ZSH_STR"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "source", "echo test"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "expected_result", "test"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "arch_amd64", "true"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "arch_arm64", "true"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "min_os_version", ""),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "max_os_version", ""),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "excluded_tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "version", "1"),
					// Read by ID, no platforms, no tags
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", sc2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", sc2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "description", "Description"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "type", "ZSH_BOOL"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "source", "echo true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "expected_result", "true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "arch_amd64", "false"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "arch_arm64", "true"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "min_os_version", "14"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "max_os_version", "15"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						ds2ResourceName, "tag_ids.*", tagResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "excluded_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						ds2ResourceName, "excluded_tag_ids.*", excludedTagResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "version", "1"),
				),
			},
		},
	})
}

func testAccMunkiScriptCheckDataSourceConfig(sc1Name string, sc2Name string) string {
	return fmt.Sprintf(`
resource "zentral_munki_script_check" "check1" {
  name            = %[1]q
  source          = "echo test"
  expected_result = "test"
}

resource "zentral_tag" "test" {
  name = %[2]q
}

resource "zentral_tag" "test-excluded" {
  name = "%[2]s-excluded"
}

resource "zentral_munki_script_check" "check2" {
  name             = %[2]q
  description      = "Description"
  type             = "ZSH_BOOL"
  source           = "echo true"
  expected_result  = "true"
  arch_amd64       = false
  arch_arm64       = true
  min_os_version   = "14"
  max_os_version   = "15"
  tag_ids          = [zentral_tag.test.id]
  excluded_tag_ids = [zentral_tag.test-excluded.id]
}

data "zentral_munki_script_check" "check1_by_name" {
  name = zentral_munki_script_check.check1.name
}

data "zentral_munki_script_check" "check2_by_id" {
  id = zentral_munki_script_check.check2.id
}
`, sc1Name, sc2Name)
}
