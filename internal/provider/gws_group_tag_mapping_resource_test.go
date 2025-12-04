package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGWSGroupTagMappingResource(t *testing.T) {
	connectionName := "TF provider GitHub"
	tagName := acctest.RandString(12)
	otherTagName := acctest.RandString(12)

	resourceName := "zentral_gws_group_tag_mapping.test"
	connectionResourceName := "data.zentral_gws_connection.test"
	tagResourceName := "zentral_tag.test"
	otherTagResourceName := "zentral_tag.test_other"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccGWSGroupTagMappingResourceFirstConfig(connectionName, tagName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name
					resource.TestCheckResourceAttr(
						resourceName, "group_email", "info@zentral.com"),
					resource.TestCheckResourceAttrPair(
						resourceName, "connection_id", connectionResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "tag_ids.*", tagResourceName, "id"),
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
				Config: testAccGWSGroupTagMappingResourceSecondConfig(connectionName, otherTagName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "group_email", "info@zentral.com"),
					resource.TestCheckResourceAttrPair(
						resourceName, "connection_id", connectionResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						resourceName, "tag_ids.*", otherTagResourceName, "id"),
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

func testAccGWSGroupTagMappingResourceFirstConfig(connectionName string, tagName string) string {
	return fmt.Sprintf(`
data "zentral_gws_connection" "test" {
  name = %[1]q
}

resource "zentral_tag" "test" {
  name = %[2]q
}

resource  "zentral_gws_group_tag_mapping" "test" {
  group_email   = "info@zentral.com"
  connection_id = data.zentral_gws_connection.test.id
  tag_ids       = [zentral_tag.test.id]
}
`, connectionName, tagName)
}

func testAccGWSGroupTagMappingResourceSecondConfig(connectionName string, otherTagName string) string {
	return fmt.Sprintf(`
data "zentral_gws_connection" "test" {
  name = %[1]q
}

resource "zentral_tag" "test_other" {
  name = %[2]q
}

resource  "zentral_gws_group_tag_mapping" "test" {
  group_email   = "info@zentral.com"
  connection_id = data.zentral_gws_connection.test.id
  tag_ids       = [zentral_tag.test_other.id]
}
`, connectionName, otherTagName)
}
