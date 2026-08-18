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
