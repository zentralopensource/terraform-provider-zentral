package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmOTAEnrollment struct {
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	DisplayName       types.String `tfsdk:"display_name"`
	BlueprintID       types.Int64  `tfsdk:"blueprint_id"`
	PushCertificateID types.Int64  `tfsdk:"push_certificate_id"`
	RealmUUID         types.String `tfsdk:"realm_uuid"`
	ACMEIssuerUUID    types.String `tfsdk:"acme_issuer_id"`
	SCEPIssuerUUID    types.String `tfsdk:"scep_issuer_id"`
	// enrollment secret
	Secret             types.String `tfsdk:"secret"`
	MetaBusinessUnitID types.Int64  `tfsdk:"meta_business_unit_id"`
	TagIDs             types.Set    `tfsdk:"tag_ids"`
	SerialNumbers      types.Set    `tfsdk:"serial_numbers"`
	UDIDs              types.Set    `tfsdk:"udids"`
	Quota              types.Int64  `tfsdk:"quota"`
}

func mdmOTAEnrollmentForState(moe *goztl.MDMOTAEnrollment) mdmOTAEnrollment {
	return mdmOTAEnrollment{
		ID:                types.Int64Value(int64(moe.ID)),
		Name:              types.StringValue(moe.Name),
		DisplayName:       types.StringValue(moe.DisplayName),
		BlueprintID:       optionalInt64ForState(moe.BlueprintID),
		PushCertificateID: types.Int64Value(int64(moe.PushCertificateID)),
		RealmUUID:         optionalStringForState(moe.RealmUUID),
		ACMEIssuerUUID:    optionalStringForState(moe.ACMEIssuerUUID),
		SCEPIssuerUUID:    types.StringValue(moe.SCEPIssuerUUID),
		// enrollment secret
		Secret:             types.StringValue(moe.Secret.Secret),
		MetaBusinessUnitID: types.Int64Value(int64(moe.Secret.MetaBusinessUnitID)),
		TagIDs:             int64SetForState(moe.Secret.TagIDs),
		SerialNumbers:      stringSetForState(moe.Secret.SerialNumbers),
		UDIDs:              stringSetForState(moe.Secret.UDIDs),
		Quota:              optionalInt64ForState(moe.Secret.Quota),
	}
}

func mdmOTAEnrollmentRequestWithState(data mdmOTAEnrollment) *goztl.MDMOTAEnrollmentRequest {
	mdmOTAEnrollmentRequest := &goztl.MDMOTAEnrollmentRequest{
		Name:              data.Name.ValueString(),
		DisplayName:       optionalStringWithState(data.DisplayName),
		BlueprintID:       optionalIntWithState(data.BlueprintID),
		PushCertificateID: int(data.PushCertificateID.ValueInt64()),
		RealmUUID:         optionalStringWithState(data.RealmUUID),
		ACMEIssuerUUID:    optionalStringWithState(data.ACMEIssuerUUID),
		SCEPIssuerUUID:    data.SCEPIssuerUUID.ValueString(),
		Secret: goztl.EnrollmentSecretRequest{
			MetaBusinessUnitID: int(data.MetaBusinessUnitID.ValueInt64()),
			TagIDs:             intListWithState(data.TagIDs),
			SerialNumbers:      stringListWithStateSet(data.SerialNumbers),
			UDIDs:              stringListWithStateSet(data.UDIDs),
			Quota:              optionalIntWithState(data.Quota),
		},
	}

	return mdmOTAEnrollmentRequest
}
