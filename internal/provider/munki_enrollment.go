package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type munkiEnrollment struct {
	ID              types.Int64  `tfsdk:"id"`
	ConfigurationID types.Int64  `tfsdk:"configuration_id"`
	PackageURL      types.String `tfsdk:"package_url"`
	Version         types.Int64  `tfsdk:"version"`
	// enrollment secret
	Secret             types.String `tfsdk:"secret"`
	MetaBusinessUnitID types.Int64  `tfsdk:"meta_business_unit_id"`
	TagIDs             types.Set    `tfsdk:"tag_ids"`
	SerialNumbers      types.Set    `tfsdk:"serial_numbers"`
	UDIDs              types.Set    `tfsdk:"udids"`
	Quota              types.Int64  `tfsdk:"quota"`
}

func munkiEnrollmentForState(me *goztl.MunkiEnrollment) munkiEnrollment {

	return munkiEnrollment{
		ID:              types.Int64Value(int64(me.ID)),
		ConfigurationID: types.Int64Value(int64(me.ConfigurationID)),
		PackageURL:      types.StringValue(me.PackageURL),
		Version:         types.Int64Value(int64(me.Version)),
		// enrollment secret
		Secret:             types.StringValue(me.Secret.Secret),
		MetaBusinessUnitID: types.Int64Value(int64(me.Secret.MetaBusinessUnitID)),
		TagIDs:             int64SetForState(me.Secret.TagIDs),
		SerialNumbers:      stringSetForState(me.Secret.SerialNumbers),
		UDIDs:              stringSetForState(me.Secret.UDIDs),
		Quota:              optionalInt64ForState(me.Secret.Quota),
	}
}

func munkiEnrollmentRequestWithState(data munkiEnrollment) *goztl.MunkiEnrollmentRequest {
	return &goztl.MunkiEnrollmentRequest{
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
