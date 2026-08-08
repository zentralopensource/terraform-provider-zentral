package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type turboEnrollment struct {
	ID               types.Int64  `tfsdk:"id"`
	ConfigurationID  types.String `tfsdk:"configuration_id"`
	ConfigProfileURL types.String `tfsdk:"configuration_profile_url"`
	PlistURL         types.String `tfsdk:"plist_url"`
	Version          types.Int64  `tfsdk:"version"`
	// enrollment secret
	Secret             types.String `tfsdk:"secret"`
	MetaBusinessUnitID types.Int64  `tfsdk:"meta_business_unit_id"`
	TagIDs             types.Set    `tfsdk:"tag_ids"`
	SerialNumbers      types.Set    `tfsdk:"serial_numbers"`
	UDIDs              types.Set    `tfsdk:"udids"`
	Quota              types.Int64  `tfsdk:"quota"`
}

func turboEnrollmentForState(te *goztl.TurboEnrollment) turboEnrollment {
	return turboEnrollment{
		ID:               types.Int64Value(int64(te.ID)),
		ConfigurationID:  types.StringValue(te.ConfigurationID),
		ConfigProfileURL: types.StringValue(te.ConfigProfileURL),
		PlistURL:         types.StringValue(te.PlistURL),
		Version:          types.Int64Value(int64(te.Version)),
		// enrollment secret
		Secret:             types.StringValue(te.Secret.Secret),
		MetaBusinessUnitID: types.Int64Value(int64(te.Secret.MetaBusinessUnitID)),
		TagIDs:             int64SetForState(te.Secret.TagIDs),
		SerialNumbers:      stringSetForState(te.Secret.SerialNumbers),
		UDIDs:              stringSetForState(te.Secret.UDIDs),
		Quota:              optionalInt64ForState(te.Secret.Quota),
	}
}

func turboEnrollmentRequestWithState(data turboEnrollment) *goztl.TurboEnrollmentRequest {
	return &goztl.TurboEnrollmentRequest{
		ConfigurationID: data.ConfigurationID.ValueString(),
		Secret: goztl.EnrollmentSecretRequest{
			MetaBusinessUnitID: int(data.MetaBusinessUnitID.ValueInt64()),
			// the helpers iterate the elements: nil if null or unknown → empty slice
			TagIDs:        intListWithState(data.TagIDs),
			SerialNumbers: stringListWithStateSet(data.SerialNumbers),
			UDIDs:         stringListWithStateSet(data.UDIDs),
			Quota:         optionalIntWithState(data.Quota),
		},
	}
}
