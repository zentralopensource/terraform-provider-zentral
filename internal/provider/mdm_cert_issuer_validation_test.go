package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// The IDent request limits are refused in the plan, before any request.
func TestMDMIDentRequestLimitsValidation(t *testing.T) {
	server := newOlderZentralServer(t)

	issuer := func(requestTimeout, maxRetries string) string {
		return testOlderZentralProviderConfig(server) + `
resource "zentral_mdm_acme_issuer" "test" {
  name           = "test"
  directory_url  = "https://www.example.com/acme"
  key_type       = "ECSECPrimeRandom"
  key_size       = 256
  usage_flags    = 1
  hardware_bound = true
  attest         = true
  backend        = "IDENT"
  ident = {
    url             = "https://www.example.com/ident"
    bearer_token    = "YoloFomo"
    request_timeout = ` + requestTimeout + `
    max_retries     = ` + maxRetries + `
  }
}
`
	}
	refused := func(config string, expectError string) resource.TestStep {
		return resource.TestStep{
			Config:      config,
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(expectError),
		}
	}
	accepted := func(config string) resource.TestStep {
		return resource.TestStep{
			Config:             config,
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		}
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			refused(issuer("31", "3"), `value\s+must\s+be\s+between\s+1\s+and\s+30,\s+got:\s+31`),
			refused(issuer("0", "3"), `value\s+must\s+be\s+between\s+1\s+and\s+30,\s+got:\s+0`),
			refused(issuer("30", "4"), `value\s+must\s+be\s+between\s+0\s+and\s+3,\s+got:\s+4`),
			refused(issuer("30", "-1"), `value\s+must\s+be\s+between\s+0\s+and\s+3,\s+got:\s+-1`),
			accepted(issuer("30", "0")),
			accepted(issuer("1", "3")),
		},
	})
}
