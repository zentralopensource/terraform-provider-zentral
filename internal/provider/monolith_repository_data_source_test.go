package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonolithRepositoryDataSource(t *testing.T) {
	r1Name := acctest.RandString(12)
	r2Name := acctest.RandString(12)
	r1ResourceName := "zentral_monolith_repository.test1"
	r2ResourceName := "zentral_monolith_repository.test2"
	mbuResourceName := "zentral_meta_business_unit.test2"
	ds1ResourceName := "data.zentral_monolith_repository.test_by_id"
	ds2ResourceName := "data.zentral_monolith_repository.test_by_name"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonolithRepositoryDataSourceConfig(r1Name, r2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by id
					resource.TestCheckResourceAttrPair(
						ds1ResourceName, "id", r1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "name", r1Name),
					resource.TestCheckResourceAttr(
						ds1ResourceName, "backend", "VIRTUAL"),
					resource.TestCheckNoResourceAttr(
						ds1ResourceName, "meta_business_unit_id"),
					// Read by name
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "id", r2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "name", r2Name),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "backend", "S3"),
					resource.TestCheckResourceAttrPair(
						ds2ResourceName, "meta_business_unit_id", mbuResourceName, "id"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "s3.bucket", "bucket"),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "s3.region_name", ""),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "s3.prefix", ""),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "s3.access_key_id", ""),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "s3.secret_access_key", ""),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "s3.assume_role_arn", ""),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "s3.signature_version", ""),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "s3.endpoint_url", ""),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "s3.cloudfront_domain", ""),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "s3.cloudfront_key_id", ""),
					resource.TestCheckResourceAttr(
						ds2ResourceName, "s3.cloudfront_privkey_pem", ""),
				),
			},
		},
	})
}

func testAccMonolithRepositoryDataSourceConfig(r1Name string, r2Name string) string {
	return fmt.Sprintf(`
resource "zentral_monolith_repository" "test1" {
  name    = %[1]q
  backend = "VIRTUAL"
}

resource "zentral_meta_business_unit" "test2" {
  name = %[2]q
}

resource "zentral_monolith_repository" "test2" {
  name                  = %[2]q
  meta_business_unit_id = zentral_meta_business_unit.test2.id
  backend               = "S3"
  s3 = {
    bucket = "bucket",
  }
}

data "zentral_monolith_repository" "test_by_id" {
  id = zentral_monolith_repository.test1.id
}

data "zentral_monolith_repository" "test_by_name" {
  name = zentral_monolith_repository.test2.name
}
`, r1Name, r2Name)
}
