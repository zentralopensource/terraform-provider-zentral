package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

	tagIDs := make([]attr.Value, 0)
	for _, tagID := range oe.Secret.TagIDs {
		tagIDs = append(tagIDs, types.Int64Value(int64(tagID)))
	}

	serialNumbers := make([]attr.Value, 0)
	for _, serialNumber := range oe.Secret.SerialNumbers {
		serialNumbers = append(serialNumbers, types.StringValue(serialNumber))
	}

	udids := make([]attr.Value, 0)
	for _, udid := range oe.Secret.UDIDs {
		udids = append(udids, types.StringValue(udid))
	}

	var quota types.Int64
	if oe.Secret.Quota != nil {
		quota = types.Int64Value(int64(*oe.Secret.Quota))
	} else {
		quota = types.Int64Null()
	}

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
		TagIDs:             types.SetValueMust(types.Int64Type, tagIDs),
		SerialNumbers:      types.SetValueMust(types.StringType, serialNumbers),
		UDIDs:              types.SetValueMust(types.StringType, udids),
		Quota:              quota,
	}
}

func osqueryEnrollmentRequestWithState(data osqueryEnrollment) *goztl.OsqueryEnrollmentRequest {
	tagIDs := make([]int, 0)
	for _, tagID := range data.TagIDs.Elements() { // nil if null or unknown → no iterations
		tagIDs = append(tagIDs, int(tagID.(types.Int64).ValueInt64()))
	}

	serialNumbers := make([]string, 0)
	for _, serialNumber := range data.SerialNumbers.Elements() { // nil if null or unknown → no iterations
		serialNumbers = append(serialNumbers, serialNumber.(types.String).ValueString())
	}

	udids := make([]string, 0)
	for _, udid := range data.UDIDs.Elements() { // nil if null or unknown → no iterations
		udids = append(udids, udid.(types.String).ValueString())
	}

	osqueryEnrollmentRequest := &goztl.OsqueryEnrollmentRequest{
		ConfigurationID: int(data.ConfigurationID.ValueInt64()),
		OsqueryRelease:  data.OsqueryRelease.ValueString(),
		Secret: goztl.EnrollmentSecretRequest{
			MetaBusinessUnitID: int(data.MetaBusinessUnitID.ValueInt64()),
			TagIDs:             tagIDs,
			SerialNumbers:      serialNumbers,
			UDIDs:              udids,
		},
	}

	if !data.Quota.IsNull() {
		osqueryEnrollmentRequest.Secret.Quota = goztl.Int(int(data.Quota.ValueInt64()))
	}

	return osqueryEnrollmentRequest
}
