package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// testAccOldZentralServer stands in for a Zentral that predates an attribute. DRF drops the keys
// the serializer does not declare, so the stub echoes the request back without them, and answers
// the refresh and the destroy that follow the failed apply. The live tenant always has the
// attribute, which leaves this the only way to reach the unsupported branch.
func testAccOldZentralServer(t *testing.T, basePath string, droppedKeys ...string) *httptest.Server {
	t.Helper()

	stored := make(map[string]interface{})
	writeStored := func(w http.ResponseWriter, status int) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(stored); err != nil {
			t.Errorf("could not encode the stub response: %v", err)
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/"+basePath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&stored); err != nil {
			t.Errorf("could not decode the stub request: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		for _, key := range droppedKeys {
			delete(stored, key)
		}
		stored["id"] = 1
		writeStored(w, http.StatusCreated)
	})

	mux.HandleFunc("/"+basePath+"1/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeStored(w, http.StatusOK)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
		}
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server
}

// The diagnostic detail is wrapped in the terraform output, so only whole tokens are matched: the
// summary, the attribute the server dropped, and the release that added it.
func regexpUnsupportedZentralVersion(attribute string, minVersion string) *regexp.Regexp {
	return regexp.MustCompile(
		"(?s)Unsupported Zentral version.*" + regexp.QuoteMeta(attribute) + ".*" + regexp.QuoteMeta(minVersion),
	)
}

// testAccOldZentralDataAssetServer stands in for a Zentral that predates the data asset source
// attribute. It drops the key it does not declare, then rejects the request for the file_uri it
// still requires, and describes the fields it accepts on OPTIONS. acceptsSource puts the attribute
// back in the metadata, for the server that does support it and refuses the value for its own
// reasons.
func testAccOldZentralDataAssetServer(t *testing.T, acceptsSource bool, withActions bool) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("/mdm/data_assets/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodOptions:
			metadata := map[string]interface{}{"name": "Data Asset List"}
			if withActions {
				fields := map[string]interface{}{
					"type":     map[string]interface{}{"type": "choice", "required": true},
					"file_uri": map[string]interface{}{"type": "string", "required": true},
				}
				if acceptsSource {
					fields["source"] = map[string]interface{}{"type": "field", "required": false}
				}
				metadata["actions"] = map[string]interface{}{"POST": fields}
			}
			if err := json.NewEncoder(w).Encode(metadata); err != nil {
				t.Errorf("could not encode the stub metadata: %v", err)
			}
		case http.MethodPost:
			w.WriteHeader(http.StatusBadRequest)
			body := map[string][]string{"file_uri": {"This field is required."}}
			if err := json.NewEncoder(w).Encode(body); err != nil {
				t.Errorf("could not encode the stub response: %v", err)
			}
		default:
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
		}
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server
}

// The artifact is given by ID, so the stub only has to serve the data asset endpoint.
func testAccMDMDataAssetSourceConfig(baseURL string) string {
	return fmt.Sprintf(`
provider "zentral" {
  base_url = %[1]q
  token    = "stub"
}

resource "zentral_mdm_data_asset" "test" {
  artifact_id = "b89d21e8-76de-4ae5-948d-5627474ab8be"
  type        = "PLIST"
  source      = "PD94bWwgdmVyc2lvbj0iMS4wIj8+"
  version     = 1
  macos       = true
}
`, baseURL+"/")
}

func TestAccMDMDataAssetResourceUnsupportedSource(t *testing.T) {
	server := testAccOldZentralDataAssetServer(t, false, true)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccMDMDataAssetSourceConfig(server.URL),
				ExpectError: regexpUnsupportedZentralVersion("source", minZentralVersionDataAssetSource),
			},
		},
	})
}

// The endpoint accepts the attribute, so the 400 is about the value. Only the error of the API is
// reported, and the practitioner is not sent after a version that is not the problem.
func TestAccMDMDataAssetResourceSourceRejectedByASupportingServer(t *testing.T) {
	server := testAccOldZentralDataAssetServer(t, true, true)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccMDMDataAssetSourceConfig(server.URL),
				ExpectError: regexp.MustCompile("(?s)Unable to create MDM data asset"),
			},
		},
	})
}

// Without the add permission, DRF leaves the actions out. The support of the attribute is unknown,
// and a guess would be wrong as often as it is right.
func TestAccMDMDataAssetResourceSourceSupportUnknown(t *testing.T) {
	server := testAccOldZentralDataAssetServer(t, false, false)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccMDMDataAssetSourceConfig(server.URL),
				ExpectError: regexp.MustCompile("(?s)Unable to create MDM data asset"),
			},
		},
	})
}

func TestAccMDMFileVaultConfigResourceUnsupportedPRKRevealRotationDelay(t *testing.T) {
	server := testAccOldZentralServer(t, "mdm/filevault_configs/", "prk_reveal_rotation_delay")
	name := acctest.RandString(12)
	escName := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "zentral" {
  base_url = %[1]q
  token    = "stub"
}

resource "zentral_mdm_filevault_config" "test" {
  name                         = %[2]q
  escrow_location_display_name = %[3]q
  prk_reveal_rotation_delay    = 120
}
`, server.URL+"/", name, escName),
				ExpectError: regexpUnsupportedZentralVersion("prk_reveal_rotation_delay", minZentralVersionPRKRevealRotationDelay),
			},
		},
	})
}

func TestAccMDMRecoveryPasswordConfigResourceUnsupportedRevealRotationDelay(t *testing.T) {
	server := testAccOldZentralServer(t, "mdm/recovery_password_configs/", "reveal_rotation_delay")
	name := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "zentral" {
  base_url = %[1]q
  token    = "stub"
}

resource "zentral_mdm_recovery_password_config" "test" {
  name                  = %[2]q
  reveal_rotation_delay = 120
}
`, server.URL+"/", name),
				ExpectError: regexpUnsupportedZentralVersion("reveal_rotation_delay", minZentralVersionRevealRotationDelay),
			},
		},
	})
}
