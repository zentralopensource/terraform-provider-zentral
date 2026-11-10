package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// The max_retries of an HTTP store is refused in the plan, before any request.
func TestStoreHTTPMaxRetriesValidation(t *testing.T) {
	server := newOlderZentralServer(t)

	store := func(maxRetries string) string {
		return testOlderZentralProviderConfig(server) + `
resource "zentral_store" "test" {
  name    = "test"
  backend = "HTTP"
  http = {
    endpoint_url = "https://www.example.com/post"
    max_retries  = ` + maxRetries + `
  }
}
`
	}
	refused := func(maxRetries string) resource.TestStep {
		return resource.TestStep{
			Config:   store(maxRetries),
			PlanOnly: true,
			ExpectError: regexp.MustCompile(
				`value\s+must\s+be\s+between\s+0\s+and\s+5,\s+got:\s+` + regexp.QuoteMeta(maxRetries),
			),
		}
	}
	accepted := func(maxRetries string) resource.TestStep {
		return resource.TestStep{
			Config:             store(maxRetries),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		}
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			refused("-1"),
			refused("6"),
			accepted("0"),
			accepted("5"),
		},
	})
}
