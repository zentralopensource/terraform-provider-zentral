package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProbeResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_probe.test"
	actionResourceName := "zentral_probe_action.test"
	mbuResourceName := "zentral_meta_business_unit.test"
	tagResourceName := "zentral_tag.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccProbeResourceBase(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckNoResourceAttr(
						resourceName, "incident_severity"),
					resource.TestCheckResourceAttr(
						resourceName, "active", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "action_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_filters.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "metadata_filters.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "payload_filters.#", "0"),
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
				Config: testAccProbeResourceFull(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "First description"),
					resource.TestCheckResourceAttr(
						resourceName, "incident_severity", "300"),
					resource.TestCheckResourceAttr(
						resourceName, "active", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "action_ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						resourceName, "action_ids.0", actionResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_filters.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_filters.0.meta_business_unit_ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						resourceName, "inventory_filters.0.meta_business_unit_ids.0", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_filters.0.platforms.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_filters.0.platforms.0", "MACOS"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_filters.0.tag_ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						resourceName, "inventory_filters.0.tag_ids.0", tagResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_filters.0.types.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "metadata_filters.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "metadata_filters.0.event_types.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "metadata_filters.0.event_types.0", "zentral_login"),
					resource.TestCheckResourceAttr(
						resourceName, "metadata_filters.0.event_routing_keys.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "metadata_filters.0.event_tags.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "metadata_filters.0.event_tags.0", "zentral"),
					resource.TestCheckResourceAttr(
						resourceName, "payload_filters.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "payload_filters.0.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "payload_filters.0.0.attribute", "yolo"),
					resource.TestCheckResourceAttr(
						resourceName, "payload_filters.0.0.operator", "IN"),
					resource.TestCheckResourceAttr(
						resourceName, "payload_filters.0.0.values.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "payload_filters.0.0.values.0", "un"),
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
				Config: testAccProbeResourceFullUpdated(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "Second description"),
					resource.TestCheckNoResourceAttr(
						resourceName, "incident_severity"),
					resource.TestCheckResourceAttr(
						resourceName, "active", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "action_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_filters.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_filters.0.meta_business_unit_ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						resourceName, "inventory_filters.0.meta_business_unit_ids.0", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_filters.0.platforms.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_filters.0.platforms.0", "MACOS"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_filters.0.tag_ids.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_filters.0.types.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "inventory_filters.0.types.0", "LAPTOP"),
					resource.TestCheckResourceAttr(
						resourceName, "metadata_filters.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "metadata_filters.0.event_types.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "metadata_filters.0.event_routing_keys.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "metadata_filters.0.event_routing_keys.0", "yolo"),
					resource.TestCheckResourceAttr(
						resourceName, "metadata_filters.0.event_tags.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "payload_filters.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "payload_filters.0.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "payload_filters.0.0.attribute", "fomo"),
					resource.TestCheckResourceAttr(
						resourceName, "payload_filters.0.0.operator", "NOT_IN"),
					resource.TestCheckResourceAttr(
						resourceName, "payload_filters.0.0.values.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "payload_filters.0.0.values.0", "trois"),
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

func testAccProbeResourceBase(name string) string {
	return fmt.Sprintf(`
resource "zentral_probe" "test" {
  name    = %[1]q
}
`, name)
}

func testAccProbeResourceFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_tag" "test" {
  name = %[1]q
}

resource "zentral_probe_action" "test" {
  name    = %[1]q
  backend = "HTTP_POST"
  http_post = {
    url = "https://www.example.com/post"
  }
}

resource "zentral_probe" "test" {
  name              = %[1]q
  description       = "First description"
  incident_severity = 300
  active            = true
  action_ids        = [zentral_probe_action.test.id]
  inventory_filters = [
    {
      meta_business_unit_ids = [zentral_meta_business_unit.test.id]
      platforms              = ["MACOS"]
      tag_ids                = [zentral_tag.test.id]
    }
  ]
  metadata_filters = [
    {
      event_types = ["zentral_login"]
      event_tags  = ["zentral"]
    }
  ]
  payload_filters = [
    [
      {
        attribute = "yolo"
        operator  = "IN"
        values    = ["un"]
      }
    ]
  ]
}
`, name)
}

func testAccProbeResourceFullUpdated(name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_probe_action" "test" {
  name    = %[1]q
  backend = "HTTP_POST"
  http_post = {
    url = "https://www.example.com/post"
  }
}

resource "zentral_probe" "test" {
  name              = %[1]q
  description       = "Second description"
  active            = true
  inventory_filters = [
    {
      meta_business_unit_ids = [zentral_meta_business_unit.test.id]
      platforms              = ["MACOS"]
      types                  = ["LAPTOP"]
    }
  ]
  metadata_filters = [
    {
      event_routing_keys = ["yolo"]
    }
  ]
  payload_filters = [
    [
      {
        attribute = "fomo"
        operator  = "NOT_IN"
        values    = ["trois"]
      }
    ]
  ]
}
`, name)
}
