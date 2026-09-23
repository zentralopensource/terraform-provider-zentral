package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const (
	tfSantaPathRegexAllow = "ALLOW"
	tfSantaPathRegexBlock = "BLOCK"
)

type santaScopedPathRegex struct {
	ID                    types.Int64  `tfsdk:"id"`
	ConfigurationID       types.Int64  `tfsdk:"configuration_id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	Policy                types.String `tfsdk:"policy"`
	Regex                 types.String `tfsdk:"regex"`
	PrimaryUsers          types.Set    `tfsdk:"primary_users"`
	ExcludedPrimaryUsers  types.Set    `tfsdk:"excluded_primary_users"`
	SerialNumbers         types.Set    `tfsdk:"serial_numbers"`
	ExcludedSerialNumbers types.Set    `tfsdk:"excluded_serial_numbers"`
	TagIDs                types.Set    `tfsdk:"tag_ids"`
	ExcludedTagIDs        types.Set    `tfsdk:"excluded_tag_ids"`
}

func santaScopedPathRegexForState(sspr *goztl.SantaScopedPathRegex) santaScopedPathRegex {
	return santaScopedPathRegex{
		ID:                    types.Int64Value(int64(sspr.ID)),
		ConfigurationID:       types.Int64Value(int64(sspr.ConfigurationID)),
		Name:                  types.StringValue(sspr.Name),
		Description:           types.StringValue(sspr.Description),
		Policy:                types.StringValue(sspr.Policy),
		Regex:                 types.StringValue(sspr.Regex),
		PrimaryUsers:          stringSetForState(sspr.PrimaryUsers),
		ExcludedPrimaryUsers:  stringSetForState(sspr.ExcludedPrimaryUsers),
		SerialNumbers:         stringSetForState(sspr.SerialNumbers),
		ExcludedSerialNumbers: stringSetForState(sspr.ExcludedSerialNumbers),
		TagIDs:                int64SetForState(sspr.TagIDs),
		ExcludedTagIDs:        int64SetForState(sspr.ExcludedTagIDs),
	}
}

func santaScopedPathRegexRequestWithState(data santaScopedPathRegex) *goztl.SantaScopedPathRegexRequest {
	return &goztl.SantaScopedPathRegexRequest{
		ConfigurationID:       int(data.ConfigurationID.ValueInt64()),
		Name:                  data.Name.ValueString(),
		Description:           data.Description.ValueString(),
		Policy:                data.Policy.ValueString(),
		Regex:                 data.Regex.ValueString(),
		PrimaryUsers:          stringListWithStateSet(data.PrimaryUsers),
		ExcludedPrimaryUsers:  stringListWithStateSet(data.ExcludedPrimaryUsers),
		SerialNumbers:         stringListWithStateSet(data.SerialNumbers),
		ExcludedSerialNumbers: stringListWithStateSet(data.ExcludedSerialNumbers),
		TagIDs:                intListWithState(data.TagIDs),
		ExcludedTagIDs:        intListWithState(data.ExcludedTagIDs),
	}
}
