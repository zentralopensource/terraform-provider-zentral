package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type santaEnrollment struct {
	ID               types.Int64  `tfsdk:"id"`
	ConfigurationID  types.Int64  `tfsdk:"configuration_id"`
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

func santaEnrollmentForState(se *goztl.SantaEnrollment) santaEnrollment {

	return santaEnrollment{
		ID:               types.Int64Value(int64(se.ID)),
		ConfigurationID:  types.Int64Value(int64(se.ConfigurationID)),
		ConfigProfileURL: types.StringValue(se.ConfigProfileURL),
		PlistURL:         types.StringValue(se.PlistURL),
		Version:          types.Int64Value(int64(se.Version)),
		// enrollment secret
		Secret:             types.StringValue(se.Secret.Secret),
		MetaBusinessUnitID: types.Int64Value(int64(se.Secret.MetaBusinessUnitID)),
		TagIDs:             int64SetForState(se.Secret.TagIDs),
		SerialNumbers:      stringSetForState(se.Secret.SerialNumbers),
		UDIDs:              stringSetForState(se.Secret.UDIDs),
		Quota:              optionalInt64ForState(se.Secret.Quota),
	}
}

func santaEnrollmentRequestWithState(data santaEnrollment) *goztl.SantaEnrollmentRequest {
	return &goztl.SantaEnrollmentRequest{
		ConfigurationID: int(data.ConfigurationID.ValueInt64()),
		Secret: goztl.EnrollmentSecretRequest{
			MetaBusinessUnitID: int(data.MetaBusinessUnitID.ValueInt64()),
			TagIDs:             intListWithState(data.TagIDs),
			SerialNumbers:      stringListWithStateSet(data.SerialNumbers),
			UDIDs:              stringListWithStateSet(data.UDIDs),
			Quota:              optionalIntWithState(data.Quota),
		},
	}
}
