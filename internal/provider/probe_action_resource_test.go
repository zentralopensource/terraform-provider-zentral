package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProbeActionResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_probe_action.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccProbeActionResourceConfigHTTPPostBase(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckResourceAttr(
						resourceName, "description", ""),
					resource.TestCheckResourceAttr(
						resourceName, "backend", "HTTP_POST"),
					resource.TestCheckResourceAttr(
						resourceName, "http_post.url", "https://www.example.com/post"),
					resource.TestCheckNoResourceAttr(
						resourceName, "http_post.username"),
					resource.TestCheckNoResourceAttr(
						resourceName, "http_post.password"),
					resource.TestCheckResourceAttr(
						resourceName, "http_post.headers.#", "0"),
					resource.TestCheckNoResourceAttr(
						resourceName, "slack_incoming_webhook"),
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
				Config: testAccProbeActionResourceConfigHTTPPostFull(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "First description"),
					resource.TestCheckResourceAttr(
						resourceName, "backend", "HTTP_POST"),
					resource.TestCheckResourceAttr(
						resourceName, "http_post.url", "https://www.example.com/post"),
					resource.TestCheckResourceAttr(
						resourceName, "http_post.username", "yolo"),
					resource.TestCheckResourceAttr(
						resourceName, "http_post.password", "fomo"),
					resource.TestCheckResourceAttr(
						resourceName, "http_post.headers.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "http_post.headers.0.name", "X-Custom-Header"),
					resource.TestCheckResourceAttr(
						resourceName, "http_post.headers.0.value", "Value"),
					resource.TestCheckNoResourceAttr(
						resourceName, "slack_incoming_webhook"),
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
				Config: testAccProbeActionResourceConfigSlack(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttr(
						resourceName, "description", "Second description"),
					resource.TestCheckResourceAttr(
						resourceName, "backend", "SLACK_INCOMING_WEBHOOK"),
					resource.TestCheckNoResourceAttr(
						resourceName, "http_post"),
					resource.TestCheckResourceAttr(
						resourceName, "slack_incoming_webhook.url", "https://www.example.com/post"),
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

func testAccProbeActionResourceConfigHTTPPostBase(name string) string {
	return fmt.Sprintf(`
resource "zentral_probe_action" "test" {
  name    = %[1]q
  backend = "HTTP_POST"
  http_post = {
    url = "https://www.example.com/post"
  }
}
`, name)
}

func testAccProbeActionResourceConfigHTTPPostFull(name string) string {
	return fmt.Sprintf(`
resource "zentral_probe_action" "test" {
  name        = %[1]q
  description = "First description"
  backend     = "HTTP_POST"
  http_post = {
    url      = "https://www.example.com/post"
    username = "yolo"
    password = "fomo"
    headers = [
      { name = "X-Custom-Header"
      value = "Value" }
    ]
  }
}
`, name)
}

func testAccProbeActionResourceConfigSlack(name string) string {
	return fmt.Sprintf(`
resource "zentral_probe_action" "test" {
  name        = %[1]q
  description = "Second description"
  backend     = "SLACK_INCOMING_WEBHOOK"
  slack_incoming_webhook = {
    url = "https://www.example.com/post"
  }
}
`, name)
}
