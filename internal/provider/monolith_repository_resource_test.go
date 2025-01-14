package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	MonolithCloudfrontTestPrivKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBAKRhksp6Bvp6Iph7vxcAT1FO3p78ek34i3Zjv5p65Yve8SC5ZCef
d3ZfYpTLsq8Bagmv2McYu1BLQcP6808qf5cCAwEAAQJBAJGPOX4EOoO4fUQLaDYE
9zenoGimZ+L9cPl/8J3pr7R/ZcJkXMIj9t7cI1rY/Tk5N2ARBZ/H3NE4Unm7xZJU
lKECIQDXoiGSvGMSB3rLKZYqyAj75O/lsh9TtZRZgF/bUBBScQIhAMMnREkKtr9d
5W7eziXRABOnVdQjPPle1KiHlSaAFmaHAiB70nUW7qixFKx1dzbs8BsAknETZBpL
FkzOrEHfDPWicQIhAKN8I7Jk7U9HY8sLj/sSKVRNnJNIqe3mSZSdcI9+QkXFAiBg
Y5iiw7n52shShyNTBggl3Xp8BILhfrIgGJ6o8jOQwA==
-----END RSA PRIVATE KEY-----`
)

func TestAccMonolithRepositoryResource(t *testing.T) {
	firstName := acctest.RandString(12)
	secondName := acctest.RandString(12)
	resourceName := "zentral_monolith_repository.test"
	mbuResourceName := "zentral_meta_business_unit.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMonolithRepositoryResourceConfigBare(firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", firstName),
					resource.TestCheckNoResourceAttr(
						resourceName, "meta_business_unit_id"),
					resource.TestCheckResourceAttr(
						resourceName, "backend", "VIRTUAL"),
					resource.TestCheckNoResourceAttr(
						resourceName, "s3"),
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
				Config: testAccMonolithRepositoryResourceConfigAzure(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttrPair(
						resourceName, "meta_business_unit_id", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "backend", "AZURE"),
					resource.TestCheckResourceAttr(
						resourceName, "azure.storage_account", "yolo"),
					resource.TestCheckResourceAttr(
						resourceName, "azure.container", "fomo"),
					resource.TestCheckResourceAttr(
						resourceName, "azure.prefix", "prefix"),
					resource.TestCheckResourceAttr(
						resourceName, "azure.client_id", "de094c71-61cb-4646-a8f7-bc4900a7040a"),
					resource.TestCheckResourceAttr(
						resourceName, "azure.tenant_id", "5237b576-d8b2-447a-a84b-c67b9c659928"),
					resource.TestCheckResourceAttr(
						resourceName, "azure.client_secret", "1da65ea6bf7d397c0776fbea04aca2cb"),
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
				Config: testAccMonolithRepositoryResourceConfigS3(secondName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", secondName),
					resource.TestCheckResourceAttrPair(
						resourceName, "meta_business_unit_id", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName, "backend", "S3"),
					resource.TestCheckResourceAttr(
						resourceName, "s3.bucket", "bucket"),
					resource.TestCheckResourceAttr(
						resourceName, "s3.region_name", "us-east-1"),
					resource.TestCheckResourceAttr(
						resourceName, "s3.prefix", "prefix"),
					resource.TestCheckResourceAttr(
						resourceName, "s3.access_key_id", "11111111111111111111"),
					resource.TestCheckResourceAttr(
						resourceName, "s3.secret_access_key", "22222222222222222222"),
					resource.TestCheckResourceAttr(
						resourceName, "s3.assume_role_arn", "arn:aws:iam::123456789012:role/S3Access"),
					resource.TestCheckResourceAttr(
						resourceName, "s3.signature_version", "s3v4"),
					resource.TestCheckResourceAttr(
						resourceName, "s3.endpoint_url", "https://endpoint.example.com"),
					resource.TestCheckResourceAttr(
						resourceName, "s3.cloudfront_domain", "yolo.cloudfront.net"),
					resource.TestCheckResourceAttr(
						resourceName, "s3.cloudfront_key_id", "YOLO"),
					resource.TestCheckResourceAttr(
						resourceName, "s3.cloudfront_privkey_pem", MonolithCloudfrontTestPrivKeyPEM),
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

func testAccMonolithRepositoryResourceConfigBare(name string) string {
	return fmt.Sprintf(`
resource "zentral_monolith_repository" "test" {
  name    = %[1]q
  backend = "VIRTUAL"
}
`, name)
}

func testAccMonolithRepositoryResourceConfigAzure(name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_monolith_repository" "test" {
  name                  = %[1]q
  meta_business_unit_id = zentral_meta_business_unit.test.id
  backend               = "AZURE"
  azure = {
    storage_account = "yolo",
    container       = "fomo",
    prefix          = "prefix",
    client_id       = "de094c71-61cb-4646-a8f7-bc4900a7040a",
    tenant_id       = "5237b576-d8b2-447a-a84b-c67b9c659928",
    client_secret   = "1da65ea6bf7d397c0776fbea04aca2cb",
  }
}
`, name)
}

func testAccMonolithRepositoryResourceConfigS3(name string) string {
	return fmt.Sprintf(`
resource "zentral_meta_business_unit" "test" {
  name = %[1]q
}

resource "zentral_monolith_repository" "test" {
  name                  = %[1]q
  meta_business_unit_id = zentral_meta_business_unit.test.id
  backend               = "S3"
  s3 = {
    bucket                 = "bucket",
    region_name            = "us-east-1",
    prefix                 = "prefix",
    access_key_id          = "11111111111111111111",
    secret_access_key      = "22222222222222222222",
    assume_role_arn        = "arn:aws:iam::123456789012:role/S3Access",
    signature_version      = "s3v4",
    endpoint_url           = "https://endpoint.example.com",
    cloudfront_domain      = "yolo.cloudfront.net",
    cloudfront_key_id      = "YOLO",
    cloudfront_privkey_pem = %[2]q
  }
}
`, name, MonolithCloudfrontTestPrivKeyPEM)
}
