package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProbeActionDataSource(t *testing.T) {
	r1Name := acctest.RandString(12)
	r2Name := acctest.RandString(12)
	r3Name := acctest.RandString(12)
	r1ResourceName := "zentral_probe_action.test1"
	r2ResourceName := "zentral_probe_action.test2"
	r3ResourceName := "zentral_probe_action.test3"
	ds1ResourceName := "data.zentral_probe_action.test_by_id"
	ds2ResourceName := "data.zentral_probe_action.test_by_name_http"
	ds3ResourceName := "data.zentral_probe_action.test_by_name_slack"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProbeActionDataSourceConfig(r1Name, r2Name, r3Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by id
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", r1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", r1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "description", ""),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "backend", "HTTP_POST"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "http_post.url", "https://www.example.com/post"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "http_post.username"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "http_post.password"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "http_post.headers.#", "0"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "slack_incoming_webhook"),
					// Read by name HTTP_POST
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", r2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", r2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "description", "First description"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "backend", "HTTP_POST"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "http_post.url", "https://www.example.com/post"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "http_post.username", "yolo"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "http_post.password", "fomo"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "http_post.headers.#", "1"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "http_post.headers.0.name", "X-Custom-Header"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "http_post.headers.0.value", "Value"),
					resource.TestCheckNoResourceAttr(
						ds2ResourceName, "slack_incoming_webhook"),
					// Read by name SLACK_INCOMING_WEBHOOK
					resource.TestCheckResourceAttrPair(
						ds3ResourceName, "id", r3ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds3ResourceName, "name", r3Name),
					resource.TestCheckResourceAttr(
						ds3ResourceName, "description", "Second description"),
					resource.TestCheckResourceAttr(
						ds3ResourceName, "backend", "SLACK_INCOMING_WEBHOOK"),
					resource.TestCheckNoResourceAttr(
						ds3ResourceName, "http_post"),
					resource.TestCheckResourceAttr(
						ds3ResourceName, "slack_incoming_webhook.url", "https://www.example.com/post"),
				),
			},
		},
	})
}

func testAccProbeActionDataSourceConfig(r1Name string, r2Name string, r3Name string) string {
	return fmt.Sprintf(`
resource "zentral_probe_action" "test1" {
  name    = %[1]q
  backend = "HTTP_POST"
  http_post = {
    url = "https://www.example.com/post"
  }
}

resource "zentral_probe_action" "test2" {
  name        = %[2]q
  description = "First description"
  backend     = "HTTP_POST"
  http_post = {
    url = "https://www.example.com/post"
    username = "yolo"
    password = "fomo"
    headers = [
      {name = "X-Custom-Header",
      value = "Value"}
    ]
  }
}

resource "zentral_probe_action" "test3" {
  name        = %[3]q
  description = "Second description"
  backend     = "SLACK_INCOMING_WEBHOOK"
  slack_incoming_webhook = {
    url = "https://www.example.com/post"
  }
}

data "zentral_probe_action" "test_by_id" {
  id = zentral_probe_action.test1.id
}

data "zentral_probe_action" "test_by_name_http" {
  name = zentral_probe_action.test2.name
}

data "zentral_probe_action" "test_by_name_slack" {
  name = zentral_probe_action.test3.name
}
`, r1Name, r2Name, r3Name)
}
