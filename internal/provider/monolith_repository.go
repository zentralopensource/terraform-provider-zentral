package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const (
	tfMonolithAzureBackend   string = "AZURE"
	tfMonolithS3Backend             = "S3"
	tfMonolithVirtualBackend        = "VIRTUAL"
)

type monolithRepository struct {
	ID                 types.Int64  `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	MetaBusinessUnitID types.Int64  `tfsdk:"meta_business_unit_id"`
	Backend            types.String `tfsdk:"backend"`
	Azure              types.Object `tfsdk:"azure"`
	S3                 types.Object `tfsdk:"s3"`
}

var azureAttrTypes = map[string]attr.Type{
	"storage_account": types.StringType,
	"container":       types.StringType,
	"prefix":          types.StringType,
	"client_id":       types.StringType,
	"tenant_id":       types.StringType,
	"client_secret":   types.StringType,
}

var s3AttrTypes = map[string]attr.Type{
	"bucket":                 types.StringType,
	"region_name":            types.StringType,
	"prefix":                 types.StringType,
	"access_key_id":          types.StringType,
	"secret_access_key":      types.StringType,
	"assume_role_arn":        types.StringType,
	"signature_version":      types.StringType,
	"endpoint_url":           types.StringType,
	"cloudfront_domain":      types.StringType,
	"cloudfront_key_id":      types.StringType,
	"cloudfront_privkey_pem": types.StringType,
}

func monolithRepositoryForState(mr *goztl.MonolithRepository) monolithRepository {
	var mbu types.Int64
	if mr.MetaBusinessUnitID != nil {
		mbu = types.Int64Value(int64(*mr.MetaBusinessUnitID))
	} else {
		mbu = types.Int64Null()
	}

	var az types.Object
	if mr.Azure != nil {
		az = types.ObjectValueMust(
			azureAttrTypes,
			map[string]attr.Value{
				"storage_account": types.StringValue(mr.Azure.StorageAccount),
				"container":       types.StringValue(mr.Azure.Container),
				"prefix":          types.StringValue(mr.Azure.Prefix),
				"client_id":       types.StringValue(mr.Azure.ClientID),
				"tenant_id":       types.StringValue(mr.Azure.TenantID),
				"client_secret":   types.StringValue(mr.Azure.ClientSecret),
			},
		)
	} else {
		az = types.ObjectNull(azureAttrTypes)
	}

	var s3 types.Object
	if mr.S3 != nil {
		s3 = types.ObjectValueMust(
			s3AttrTypes,
			map[string]attr.Value{
				"bucket":                 types.StringValue(mr.S3.Bucket),
				"region_name":            types.StringValue(mr.S3.RegionName),
				"prefix":                 types.StringValue(mr.S3.Prefix),
				"access_key_id":          types.StringValue(mr.S3.AccessKeyID),
				"secret_access_key":      types.StringValue(mr.S3.SecretAccessKey),
				"assume_role_arn":        types.StringValue(mr.S3.AssumeRoleARN),
				"signature_version":      types.StringValue(mr.S3.SignatureVersion),
				"endpoint_url":           types.StringValue(mr.S3.EndpointURL),
				"cloudfront_domain":      types.StringValue(mr.S3.CloudfrontDomain),
				"cloudfront_key_id":      types.StringValue(mr.S3.CloudfrontKeyID),
				"cloudfront_privkey_pem": types.StringValue(mr.S3.CloudfrontPrivkeyPEM),
			},
		)
	} else {
		s3 = types.ObjectNull(s3AttrTypes)
	}

	return monolithRepository{
		ID:                 types.Int64Value(int64(mr.ID)),
		Name:               types.StringValue(mr.Name),
		MetaBusinessUnitID: mbu,
		Backend:            types.StringValue(mr.Backend),
		Azure:              az,
		S3:                 s3,
	}
}

func monolithRepositoryRequestWithState(data monolithRepository) *goztl.MonolithRepositoryRequest {
	var mbu *int
	if !data.MetaBusinessUnitID.IsNull() {
		mbu = goztl.Int(int(data.MetaBusinessUnitID.ValueInt64()))
	}

	req := &goztl.MonolithRepositoryRequest{
		Name:               data.Name.ValueString(),
		MetaBusinessUnitID: mbu,
		Backend:            data.Backend.ValueString(),
	}

	if !data.Azure.IsNull() {
		azMap := data.Azure.Attributes()
		if azMap != nil {
			azBackend := &goztl.MonolithAzureBackend{
				StorageAccount: azMap["storage_account"].(types.String).ValueString(),
				Container:      azMap["container"].(types.String).ValueString(),
				Prefix:         azMap["prefix"].(types.String).ValueString(),
				ClientID:       azMap["client_id"].(types.String).ValueString(),
				TenantID:       azMap["tenant_id"].(types.String).ValueString(),
				ClientSecret:   azMap["client_secret"].(types.String).ValueString(),
			}
			req.Azure = azBackend
		}
	}

	if !data.S3.IsNull() {
		s3Map := data.S3.Attributes()
		if s3Map != nil {
			s3Backend := &goztl.MonolithS3Backend{
				Bucket:               s3Map["bucket"].(types.String).ValueString(),
				RegionName:           s3Map["region_name"].(types.String).ValueString(),
				Prefix:               s3Map["prefix"].(types.String).ValueString(),
				AccessKeyID:          s3Map["access_key_id"].(types.String).ValueString(),
				SecretAccessKey:      s3Map["secret_access_key"].(types.String).ValueString(),
				AssumeRoleARN:        s3Map["assume_role_arn"].(types.String).ValueString(),
				SignatureVersion:     s3Map["signature_version"].(types.String).ValueString(),
				EndpointURL:          s3Map["endpoint_url"].(types.String).ValueString(),
				CloudfrontDomain:     s3Map["cloudfront_domain"].(types.String).ValueString(),
				CloudfrontKeyID:      s3Map["cloudfront_key_id"].(types.String).ValueString(),
				CloudfrontPrivkeyPEM: s3Map["cloudfront_privkey_pem"].(types.String).ValueString(),
			}
			req.S3 = s3Backend
		}
	}

	return req
}
