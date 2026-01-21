package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMDMEnrollmentCustomViewResource_basic(t *testing.T) {
	name := acctest.RandString(12)
	description := acctest.RandString(12)
	resourceName := "zentral_mdm_enrollment_custom_view.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMEnrollmentCustomViewResourceConfigBare(name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName, "description", description),
					resource.TestCheckResourceAttr(
						resourceName, "html", "<html/>"),
					resource.TestCheckResourceAttr(
						resourceName, "requires_authentication", "false"),
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
				Config: testAccMDMEnrollmentCustomViewResourceConfigFull(name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName, "description", description),
					resource.TestCheckResourceAttr(
						resourceName, "html", "<html><body>updated</body></html>"),
					resource.TestCheckResourceAttr(
						resourceName, "requires_authentication", "true"),
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

func TestAccMDMEnrollmentCustomViewResource_file(t *testing.T) {
	name := acctest.RandString(12)
	description := acctest.RandString(12)
	resourceName := "zentral_mdm_enrollment_custom_view.test"

	// Create temp HTML files for html_file testing
	tmpDir := t.TempDir()
	htmlPath1 := filepath.Join(tmpDir, "view1.html")
	htmlPath2 := filepath.Join(tmpDir, "view2.html")
	if err := os.WriteFile(htmlPath1, []byte("<html/>"), 0o644); err != nil {
		t.Fatalf("write htmlPath1: %v", err)
	}
	if err := os.WriteFile(htmlPath2, []byte("<html><body>updated</body></html>"), 0o644); err != nil {
		t.Fatalf("write htmlPath2: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMDMEnrollmentCustomViewResourceConfigBareFile(name, description, htmlPath1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName, "description", description),
					resource.TestCheckResourceAttr(
						resourceName, "html", "<html/>"),
					resource.TestCheckResourceAttr(
						resourceName, "requires_authentication", "false"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"html_file",
				},
			},
			// Update and Read
			{
				Config: testAccMDMEnrollmentCustomViewResourceConfigFullFile(name, description, htmlPath2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName, "description", description),
					resource.TestCheckResourceAttr(
						resourceName, "html", "<html><body>updated</body></html>"),
					resource.TestCheckResourceAttr(
						resourceName, "requires_authentication", "true"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"html_file",
				},
			},
		},
	})
}

func testAccMDMEnrollmentCustomViewResourceConfigBare(name string, description string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_enrollment_custom_view" "test" {
	name                  			= %[1]q
	description						= %[2]q
	html							= "<html/>"
}
`, name, description)
}

func testAccMDMEnrollmentCustomViewResourceConfigFull(name string, description string) string {
	return fmt.Sprintf(`
resource "zentral_mdm_enrollment_custom_view" "test" {
	name                  			= %[1]q
	description						= %[2]q
	html							= "<html><body>updated</body></html>"
	requires_authentication			= true
}
`, name, description)
}

func testAccMDMEnrollmentCustomViewResourceConfigBareFile(name string, description string, htmlFile string) string {
	return fmt.Sprintf(`
 resource "zentral_mdm_enrollment_custom_view" "test" {
 	name                  			= %[1]q
 	description						= %[2]q
	html_file						= %[3]q
 }
`, name, description, htmlFile)
}

func testAccMDMEnrollmentCustomViewResourceConfigFullFile(name string, description string, htmlFile string) string {
	return fmt.Sprintf(`
 resource "zentral_mdm_enrollment_custom_view" "test" {
 	name                  			= %[1]q
 	description						= %[2]q
	html_file						= %[3]q
 	requires_authentication			= true
 }
`, name, description, htmlFile)
}
