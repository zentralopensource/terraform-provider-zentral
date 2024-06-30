package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmOTAEnrollment struct {
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	DisplayName       types.String `tfsdk:"display_name"`
	BlueprintID       types.Int64  `tfsdk:"blueprint_id"`
	PushCertificateID types.Int64  `tfsdk:"push_certificate_id"`
	RealmID           types.Int64  `tfsdk:"realm_id"`
	SCEPConfigID      types.Int64  `tfsdk:"scep_config_id"`
	SCEPVerification  types.Bool   `tfsdk:"scep_verification"`
	// enrollment secret
	Secret             types.String `tfsdk:"secret"`
	MetaBusinessUnitID types.Int64  `tfsdk:"meta_business_unit_id"`
	TagIDs             types.Set    `tfsdk:"tag_ids"`
	SerialNumbers      types.Set    `tfsdk:"serial_numbers"`
	UDIDs              types.Set    `tfsdk:"udids"`
	Quota              types.Int64  `tfsdk:"quota"`
}

func mdmOTAEnrollmentForState(moe *goztl.MDMOTAEnrollment) mdmOTAEnrollment {
	var blueprintID types.Int64
	if moe.BlueprintID != nil {
		blueprintID = types.Int64Value(int64(*moe.BlueprintID))
	} else {
		blueprintID = types.Int64Null()
	}

	var realmID types.Int64
	if moe.RealmID != nil {
		realmID = types.Int64Value(int64(*moe.RealmID))
	} else {
		realmID = types.Int64Null()
	}

	tagIDs := make([]attr.Value, 0)
	for _, tagID := range moe.Secret.TagIDs {
		tagIDs = append(tagIDs, types.Int64Value(int64(tagID)))
	}

	serialNumbers := make([]attr.Value, 0)
	for _, serialNumber := range moe.Secret.SerialNumbers {
		serialNumbers = append(serialNumbers, types.StringValue(serialNumber))
	}

	udids := make([]attr.Value, 0)
	for _, udid := range moe.Secret.UDIDs {
		udids = append(udids, types.StringValue(udid))
	}

	var quota types.Int64
	if moe.Secret.Quota != nil {
		quota = types.Int64Value(int64(*moe.Secret.Quota))
	} else {
		quota = types.Int64Null()
	}

	return mdmOTAEnrollment{
		ID:                types.Int64Value(int64(moe.ID)),
		Name:              types.StringValue(moe.Name),
		DisplayName:       types.StringValue(moe.DisplayName),
		BlueprintID:       blueprintID,
		PushCertificateID: types.Int64Value(int64(moe.PushCertificateID)),
		RealmID:           realmID,
		SCEPConfigID:      types.Int64Value(int64(moe.SCEPConfigID)),
		SCEPVerification:  types.BoolValue(moe.SCEPVerification),
		// enrollment secret
		Secret:             types.StringValue(moe.Secret.Secret),
		MetaBusinessUnitID: types.Int64Value(int64(moe.Secret.MetaBusinessUnitID)),
		TagIDs:             types.SetValueMust(types.Int64Type, tagIDs),
		SerialNumbers:      types.SetValueMust(types.StringType, serialNumbers),
		UDIDs:              types.SetValueMust(types.StringType, udids),
		Quota:              quota,
	}
}

func mdmOTAEnrollmentRequestWithState(data mdmOTAEnrollment) *goztl.MDMOTAEnrollmentRequest {
	var bpID *int
	if !data.BlueprintID.IsNull() {
		bpID = goztl.Int(int(data.BlueprintID.ValueInt64()))
	}

	var rID *int
	if !data.RealmID.IsNull() {
		rID = goztl.Int(int(data.RealmID.ValueInt64()))
	}

	var dn *string
	if !data.DisplayName.IsNull() {
		dn = goztl.String(data.DisplayName.ValueString())
	}

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

	mdmOTAEnrollmentRequest := &goztl.MDMOTAEnrollmentRequest{
		Name:              data.Name.ValueString(),
		DisplayName:       dn,
		BlueprintID:       bpID,
		PushCertificateID: int(data.PushCertificateID.ValueInt64()),
		RealmID:           rID,
		SCEPConfigID:      int(data.SCEPConfigID.ValueInt64()),
		SCEPVerification:  data.SCEPVerification.ValueBool(),
		Secret: goztl.EnrollmentSecretRequest{
			MetaBusinessUnitID: int(data.MetaBusinessUnitID.ValueInt64()),
			TagIDs:             tagIDs,
			SerialNumbers:      serialNumbers,
			UDIDs:              udids,
		},
	}

	if !data.Quota.IsNull() {
		mdmOTAEnrollmentRequest.Secret.Quota = goztl.Int(int(data.Quota.ValueInt64()))
	}

	return mdmOTAEnrollmentRequest
}
