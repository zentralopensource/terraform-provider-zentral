package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// The attributes of a Santa configuration before v2026.6, which added the block notification button.
var olderZentralSantaConfigurationAttributes = []string{
	"name", "client_mode", "client_certificate_auth", "batch_size", "full_sync_interval",
	"enable_bundles", "enable_transitive_rules", "allowed_path_regex", "blocked_path_regex",
	"block_usb_mount", "remount_usb_mode", "allow_unknown_shard", "enable_all_event_upload_shard",
	"sync_incident_severity",
}

// newOlderZentralServer answers the Santa configuration endpoints like a Zentral that predates the
// block notification button and the scoped configuration items: DRF drops the keys it does not know
// from the request, and the scoped endpoints do not exist.
func newOlderZentralServer(t *testing.T) *httptest.Server {
	var mu sync.Mutex
	configurations := map[int]map[string]any{}
	nextID := 1

	const basePath = "/api/santa/configurations/"

	write := func(w http.ResponseWriter, status int, body any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}
	read := func(r *http.Request, into map[string]any) error {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return err
		}
		for _, attr := range olderZentralSantaConfigurationAttributes {
			if value, ok := body[attr]; ok {
				into[attr] = value
			}
		}
		return nil
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		if !strings.HasPrefix(r.URL.Path, basePath) {
			write(w, http.StatusNotFound, map[string]string{"detail": "Not found."})
			return
		}
		idPart := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, basePath), "/")

		if idPart == "" {
			if r.Method != http.MethodPost {
				write(w, http.StatusMethodNotAllowed, nil)
				return
			}
			configuration := map[string]any{
				"id":         nextID,
				"created_at": "2026-09-23T10:00:00.000000",
				"updated_at": "2026-09-23T10:00:00.000000",
			}
			if err := read(r, configuration); err != nil {
				write(w, http.StatusBadRequest, map[string]string{"detail": err.Error()})
				return
			}
			configurations[nextID] = configuration
			nextID++
			write(w, http.StatusCreated, configuration)
			return
		}

		id, err := strconv.Atoi(idPart)
		configuration, found := configurations[id]
		if err != nil || !found {
			write(w, http.StatusNotFound, map[string]string{"detail": "Not found."})
			return
		}
		switch r.Method {
		case http.MethodGet:
			write(w, http.StatusOK, configuration)
		case http.MethodPut:
			if err := read(r, configuration); err != nil {
				write(w, http.StatusBadRequest, map[string]string{"detail": err.Error()})
				return
			}
			write(w, http.StatusOK, configuration)
		case http.MethodDelete:
			delete(configurations, id)
			write(w, http.StatusNoContent, nil)
		default:
			write(w, http.StatusMethodNotAllowed, nil)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

func testOlderZentralProviderConfig(server *httptest.Server) string {
	return fmt.Sprintf(`
provider "zentral" {
  base_url = "%s/api/"
  token    = "older-zentral"
}
`, server.URL)
}

// A configuration that leaves the block notification button alone keeps working against a Zentral
// that predates it, and one that sets it is told which release it needs.
func TestSantaConfigurationResourceOlderZentral(t *testing.T) {
	server := newOlderZentralServer(t)
	resourceName := "zentral_santa_configuration.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read, then an empty plan
			{
				Config: testOlderZentralProviderConfig(server) + `
resource "zentral_santa_configuration" "test" {
  name = "older"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "event_detail_source", "LOCAL"),
					resource.TestCheckResourceAttr(resourceName, "event_detail_url", ""),
					resource.TestCheckResourceAttr(resourceName, "event_detail_text", ""),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read, with the default source spelled out
			{
				Config: testOlderZentralProviderConfig(server) + `
resource "zentral_santa_configuration" "test" {
  name                = "older, renamed"
  client_mode         = "LOCKDOWN"
  event_detail_source = "LOCAL"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "older, renamed"),
					resource.TestCheckResourceAttr(resourceName, "client_mode", "LOCKDOWN"),
					resource.TestCheckResourceAttr(resourceName, "event_detail_source", "LOCAL"),
				),
			},
			// Update with a button the server cannot store
			{
				Config: testOlderZentralProviderConfig(server) + `
resource "zentral_santa_configuration" "test" {
  name                = "older, renamed"
  client_mode         = "LOCKDOWN"
  event_detail_source = "CUSTOM"
  event_detail_url    = "https://www.example.com/blocked/"
}
`,
				ExpectError: regexp.MustCompile(`ignored\s+the\s+event_detail_source\s+attribute[\s\S]*requires\s+Zentral\s+v2026\.6`),
			},
		},
	})
}

func TestSantaScopedConfigurationItemsOlderZentral(t *testing.T) {
	server := newOlderZentralServer(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testOlderZentralProviderConfig(server) + `
resource "zentral_santa_configuration" "test" {
  name = "older"
}

resource "zentral_santa_scoped_client_mode" "test" {
  configuration_id = zentral_santa_configuration.test.id
  name             = "lockdown"
  client_mode      = "LOCKDOWN"
}
`,
				ExpectError: regexp.MustCompile(`no\s+Santa\s+scoped\s+client\s+mode\s+endpoint[\s\S]*requires\s+Zentral\s+v2026\.6`),
			},
			{
				Config: testOlderZentralProviderConfig(server) + `
resource "zentral_santa_configuration" "test" {
  name = "older"
}

resource "zentral_santa_scoped_path_regex" "test" {
  configuration_id = zentral_santa_configuration.test.id
  name             = "applications"
  policy           = "ALLOW"
  regex            = "/Applications/.+"
}
`,
				ExpectError: regexp.MustCompile(`no\s+Santa\s+scoped\s+path\s+regex\s+endpoint[\s\S]*requires\s+Zentral\s+v2026\.6`),
			},
		},
	})
}
