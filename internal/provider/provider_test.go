package provider

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"zentral": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testQueryArg(t *testing.T, r *http.Request, arg string, want string) {
	t.Helper()
	if got := r.URL.Query().Get(arg); got != want {
		t.Errorf("Request query arg %q: value %q, want %q", arg, got, want)
	}
}

func setupMockedReverseProxy(t *testing.T) (mux *http.ServeMux, srvURL string, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	realBaseURL := "https://zaio.zentral.dev/api/"

	real, err := url.Parse(realBaseURL)
	if err != nil {
		t.Fatalf("invalid real base url %q: %v", realBaseURL, err)
	}
	rp := httputil.NewSingleHostReverseProxy(real)

	origDirector := rp.Director
	rp.Director = func(r *http.Request) {
		origDirector(r)
		r.Host = real.Host
	}

	mux.Handle("/", rp)

	return mux, server.URL, server.Close
}
