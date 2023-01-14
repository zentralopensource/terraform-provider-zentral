package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

	tagIDs := make([]attr.Value, 0)
	for _, tagID := range se.Secret.TagIDs {
		tagIDs = append(tagIDs, types.Int64Value(int64(tagID)))
	}

	serialNumbers := make([]attr.Value, 0)
	for _, serialNumber := range se.Secret.SerialNumbers {
		serialNumbers = append(serialNumbers, types.StringValue(serialNumber))
	}

	udids := make([]attr.Value, 0)
	for _, udid := range se.Secret.UDIDs {
		udids = append(udids, types.StringValue(udid))
	}

	var quota types.Int64
	if se.Secret.Quota != nil {
		quota = types.Int64Value(int64(*se.Secret.Quota))
	} else {
		quota = types.Int64Null()
	}

	return santaEnrollment{
		ID:               types.Int64Value(int64(se.ID)),
		ConfigurationID:  types.Int64Value(int64(se.ConfigurationID)),
		ConfigProfileURL: types.StringValue(se.ConfigProfileURL),
		PlistURL:         types.StringValue(se.PlistURL),
		Version:          types.Int64Value(int64(se.Version)),
		// enrollment secret
		Secret:             types.StringValue(se.Secret.Secret),
		MetaBusinessUnitID: types.Int64Value(int64(se.Secret.MetaBusinessUnitID)),
		TagIDs:             types.SetValueMust(types.Int64Type, tagIDs),
		SerialNumbers:      types.SetValueMust(types.StringType, serialNumbers),
		UDIDs:              types.SetValueMust(types.StringType, udids),
		Quota:              quota,
	}
}

func santaEnrollmentRequestWithState(data santaEnrollment) *goztl.SantaEnrollmentRequest {
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

	santaEnrollmentRequest := &goztl.SantaEnrollmentRequest{
		ConfigurationID: int(data.ConfigurationID.ValueInt64()),
		Secret: goztl.EnrollmentSecretRequest{
			MetaBusinessUnitID: int(data.MetaBusinessUnitID.ValueInt64()),
			TagIDs:             tagIDs,
			SerialNumbers:      serialNumbers,
			UDIDs:              udids,
		},
	}

	if !data.Quota.IsNull() {
		santaEnrollmentRequest.Secret.Quota = goztl.Int(int(data.Quota.ValueInt64()))
	}

	return santaEnrollmentRequest
}
