package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type osqueryEnrollment struct {
	ID                  types.Int64  `tfsdk:"id"`
	ConfigurationID     types.Int64  `tfsdk:"configuration_id"`
	OsqueryRelease      types.String `tfsdk:"osquery_release"`
	PackageURL          types.String `tfsdk:"package_url"`
	ScriptURL           types.String `tfsdk:"script_url"`
	PowershellScriptURL types.String `tfsdk:"powershell_script_url"`
	Version             types.Int64  `tfsdk:"version"`
	// enrollment secret
	Secret             types.String `tfsdk:"secret"`
	MetaBusinessUnitID types.Int64  `tfsdk:"meta_business_unit_id"`
	TagIDs             types.Set    `tfsdk:"tag_ids"`
	SerialNumbers      types.Set    `tfsdk:"serial_numbers"`
	UDIDs              types.Set    `tfsdk:"udids"`
	Quota              types.Int64  `tfsdk:"quota"`
}

func osqueryEnrollmentForState(oe *goztl.OsqueryEnrollment) osqueryEnrollment {

	return osqueryEnrollment{
		ID:                  types.Int64Value(int64(oe.ID)),
		ConfigurationID:     types.Int64Value(int64(oe.ConfigurationID)),
		OsqueryRelease:      types.StringValue(oe.OsqueryRelease),
		PackageURL:          types.StringValue(oe.PackageURL),
		ScriptURL:           types.StringValue(oe.ScriptURL),
		PowershellScriptURL: types.StringValue(oe.PowershellScriptURL),
		Version:             types.Int64Value(int64(oe.Version)),
		// enrollment secret
		Secret:             types.StringValue(oe.Secret.Secret),
		MetaBusinessUnitID: types.Int64Value(int64(oe.Secret.MetaBusinessUnitID)),
		TagIDs:             int64SetForState(oe.Secret.TagIDs),
		SerialNumbers:      stringSetForState(oe.Secret.SerialNumbers),
		UDIDs:              stringSetForState(oe.Secret.UDIDs),
		Quota:              optionalInt64ForState(oe.Secret.Quota),
	}
}

func osqueryEnrollmentRequestWithState(data osqueryEnrollment) *goztl.OsqueryEnrollmentRequest {
	return &goztl.OsqueryEnrollmentRequest{
		ConfigurationID: int(data.ConfigurationID.ValueInt64()),
		OsqueryRelease:  data.OsqueryRelease.ValueString(),
		Secret: goztl.EnrollmentSecretRequest{
			MetaBusinessUnitID: int(data.MetaBusinessUnitID.ValueInt64()),
			TagIDs:             intListWithState(data.TagIDs),
			SerialNumbers:      stringListWithStateSet(data.SerialNumbers),
			UDIDs:              stringListWithStateSet(data.UDIDs),
			Quota:              optionalIntWithState(data.Quota),
		},
	}
}
