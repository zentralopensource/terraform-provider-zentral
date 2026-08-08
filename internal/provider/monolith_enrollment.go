package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type monolithEnrollment struct {
	ID               types.Int64  `tfsdk:"id"`
	ManifestID       types.Int64  `tfsdk:"manifest_id"`
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

func monolithEnrollmentForState(me *goztl.MonolithEnrollment) monolithEnrollment {

	return monolithEnrollment{
		ID:               types.Int64Value(int64(me.ID)),
		ManifestID:       types.Int64Value(int64(me.ManifestID)),
		ConfigProfileURL: types.StringValue(me.ConfigProfileURL),
		PlistURL:         types.StringValue(me.PlistURL),
		Version:          types.Int64Value(int64(me.Version)),
		// enrollment secret
		Secret:             types.StringValue(me.Secret.Secret),
		MetaBusinessUnitID: types.Int64Value(int64(me.Secret.MetaBusinessUnitID)),
		TagIDs:             int64SetForState(me.Secret.TagIDs),
		SerialNumbers:      stringSetForState(me.Secret.SerialNumbers),
		UDIDs:              stringSetForState(me.Secret.UDIDs),
		Quota:              optionalInt64ForState(me.Secret.Quota),
	}
}

func monolithEnrollmentRequestWithState(data monolithEnrollment) *goztl.MonolithEnrollmentRequest {
	return &goztl.MonolithEnrollmentRequest{
		ManifestID: int(data.ManifestID.ValueInt64()),
		Secret: goztl.EnrollmentSecretRequest{
			MetaBusinessUnitID: int(data.MetaBusinessUnitID.ValueInt64()),
			TagIDs:             intListWithState(data.TagIDs),
			SerialNumbers:      stringListWithStateSet(data.SerialNumbers),
			UDIDs:              stringListWithStateSet(data.UDIDs),
			Quota:              optionalIntWithState(data.Quota),
		},
	}
}
