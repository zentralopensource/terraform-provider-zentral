package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmDEPEnrollment struct {
	ID                         types.Int64  `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	DisplayName                types.String `tfsdk:"display_name"`
	UseRealmUser               types.Bool   `tfsdk:"use_realm_user"`
	UsernamePattern            types.String `tfsdk:"username_pattern"`
	RealmUserIsAdmin           types.Bool   `tfsdk:"realm_user_is_admin"`
	AdminFullName              types.String `tfsdk:"admin_full_name"`
	AdminShortName             types.String `tfsdk:"admin_short_name"`
	HiddenAdmin                types.Bool   `tfsdk:"hidden_admin"`
	AdminPasswordComplexity    types.Int64  `tfsdk:"admin_password_complexity"`
	AdminPasswordRotationDelay types.Int64  `tfsdk:"admin_password_rotation_delay"`
	AllowPairing               types.Bool   `tfsdk:"allow_pairing"`
	AutoAdvanceSetup           types.Bool   `tfsdk:"auto_advance_setup"`
	AwaitDeviceConfigured      types.Bool   `tfsdk:"await_device_configured"`
	Department                 types.String `tfsdk:"department"`
	IsMandatory                types.Bool   `tfsdk:"is_mandatory"`
	IsMdmRemovable             types.Bool   `tfsdk:"is_mdm_removable"`
	IsMultiUser                types.Bool   `tfsdk:"is_multi_user"`
	IsSupervised               types.Bool   `tfsdk:"is_supervised"`
	Language                   types.String `tfsdk:"language"`
	OrgMagic                   types.String `tfsdk:"org_magic"`
	Region                     types.String `tfsdk:"region"`
	SkipSetupItems             types.Set    `tfsdk:"skip_setup_items"`
	SupportEmailAddress        types.String `tfsdk:"support_email_address"`
	SupportPhoneNumber         types.String `tfsdk:"support_phone_number"`
	IncludeTlsCertificates     types.Bool   `tfsdk:"include_tls_certificates"`
	IosMaxVersion              types.String `tfsdk:"ios_max_version"`
	IosMinVersion              types.String `tfsdk:"ios_min_version"`
	MacosMaxVersion            types.String `tfsdk:"macos_max_version"`
	MacosMinVersion            types.String `tfsdk:"macos_min_version"`
	PushCertificateID          types.Int64  `tfsdk:"push_certificate_id"`
	ACMEIssuerUUID             types.String `tfsdk:"acme_issuer_id"`
	SCEPIssuerUUID             types.String `tfsdk:"scep_issuer_id"`
	BlueprintID                types.Int64  `tfsdk:"blueprint_id"`
	RealmUUID                  types.String `tfsdk:"realm_uuid"`
	VirtualServerID            types.Int64  `tfsdk:"virtual_server_id"`
	// enrollment secret
	Secret             types.String `tfsdk:"secret"`
	MetaBusinessUnitID types.Int64  `tfsdk:"meta_business_unit_id"`
	TagIDs             types.Set    `tfsdk:"tag_ids"`
	SerialNumbers      types.Set    `tfsdk:"serial_numbers"`
	UDIDs              types.Set    `tfsdk:"udids"`
	Quota              types.Int64  `tfsdk:"quota"`
}

func mdmDEPEnrollmentForState(enrollment *goztl.MDMDEPEnrollment) mdmDEPEnrollment {
	return mdmDEPEnrollment{
		ID:                         types.Int64Value(int64(enrollment.ID)),
		Name:                       types.StringValue(enrollment.Name),
		DisplayName:                types.StringValue(enrollment.DisplayName),
		UseRealmUser:               types.BoolValue(enrollment.UseRealmUser),
		UsernamePattern:            types.StringValue(enrollment.UsernamePattern),
		RealmUserIsAdmin:           types.BoolValue(enrollment.RealmUserIsAdmin),
		AdminFullName:              optionalStringForState(enrollment.AdminFullName),
		AdminShortName:             optionalStringForState(enrollment.AdminShortName),
		HiddenAdmin:                types.BoolValue(enrollment.HiddenAdmin),
		AdminPasswordComplexity:    types.Int64Value(int64(enrollment.AdminPasswordComplexity)),
		AdminPasswordRotationDelay: types.Int64Value(int64(enrollment.AdminPasswordRotationDelay)),
		AllowPairing:               types.BoolValue(enrollment.AllowPairing),
		AutoAdvanceSetup:           types.BoolValue(enrollment.AutoAdvanceSetup),
		AwaitDeviceConfigured:      types.BoolValue(enrollment.AwaitDeviceConfigured),
		Department:                 types.StringValue(enrollment.Department),
		IsMandatory:                types.BoolValue(enrollment.IsMandatory),
		IsMdmRemovable:             types.BoolValue(enrollment.IsMdmRemovable),
		IsMultiUser:                types.BoolValue(enrollment.IsMultiUser),
		IsSupervised:               types.BoolValue(enrollment.IsSupervised),
		Language:                   types.StringValue(enrollment.Language),
		OrgMagic:                   types.StringValue(enrollment.OrgMagic),
		Region:                     types.StringValue(enrollment.Region),
		SkipSetupItems:             stringSetForState(enrollment.SkipSetupItems),
		SupportEmailAddress:        types.StringValue(enrollment.SupportEmailAddress),
		SupportPhoneNumber:         types.StringValue(enrollment.SupportPhoneNumber),
		IncludeTlsCertificates:     types.BoolValue(enrollment.IncludeTlsCertificates),
		IosMaxVersion:              types.StringValue(enrollment.IosMaxVersion),
		IosMinVersion:              types.StringValue(enrollment.IosMinVersion),
		MacosMaxVersion:            types.StringValue(enrollment.MacosMaxVersion),
		MacosMinVersion:            types.StringValue(enrollment.MacosMinVersion),
		PushCertificateID:          types.Int64Value(int64(enrollment.PushCertificateID)),
		ACMEIssuerUUID:             optionalStringForState(enrollment.ACMEIssuerUUID),
		SCEPIssuerUUID:             types.StringValue(enrollment.SCEPIssuerUUID),
		BlueprintID:                optionalInt64ForState(enrollment.BlueprintID),
		RealmUUID:                  optionalStringForState(enrollment.RealmUUID),
		VirtualServerID:            types.Int64Value(int64(enrollment.VirtualServer)),
		// enrollment secret
		Secret:             types.StringValue(enrollment.Secret.Secret),
		MetaBusinessUnitID: types.Int64Value(int64(enrollment.Secret.MetaBusinessUnitID)),
		TagIDs:             int64SetForState(enrollment.Secret.TagIDs),
		SerialNumbers:      stringSetForState(enrollment.Secret.SerialNumbers),
		UDIDs:              stringSetForState(enrollment.Secret.UDIDs),
		Quota:              optionalInt64ForState(enrollment.Secret.Quota),
	}
}

func mdmDEPEnrollmentCreateRequestWithState(data mdmDEPEnrollment) *goztl.MDMDEPEnrollmentCreationRequest {
	mdmDEPEnrollmentRequest := &goztl.MDMDEPEnrollmentCreationRequest{
		Name:        data.Name.ValueString(),
		DisplayName: data.DisplayName.ValueString(),
		Secret: goztl.EnrollmentSecretRequest{
			MetaBusinessUnitID: int(data.MetaBusinessUnitID.ValueInt64()),
			TagIDs:             intListWithState(data.TagIDs),
			SerialNumbers:      stringListWithStateSet(data.SerialNumbers),
			UDIDs:              stringListWithStateSet(data.UDIDs),
			Quota:              optionalIntWithState(data.Quota),
		},
		UseRealmUser:               data.UseRealmUser.ValueBool(),
		UsernamePattern:            data.UsernamePattern.ValueString(),
		RealmUserIsAdmin:           data.RealmUserIsAdmin.ValueBool(),
		AdminFullName:              optionalStringWithState(data.AdminFullName),
		AdminShortName:             optionalStringWithState(data.AdminShortName),
		HiddenAdmin:                data.HiddenAdmin.ValueBool(),
		AdminPasswordComplexity:    int(data.AdminPasswordComplexity.ValueInt64()),
		AdminPasswordRotationDelay: int(data.AdminPasswordRotationDelay.ValueInt64()),
		AllowPairing:               data.AllowPairing.ValueBool(),
		AutoAdvanceSetup:           data.AutoAdvanceSetup.ValueBool(),
		AwaitDeviceConfigured:      data.AwaitDeviceConfigured.ValueBool(),
		Department:                 data.Department.ValueString(),
		IsMandatory:                data.IsMandatory.ValueBool(),
		IsMdmRemovable:             data.IsMdmRemovable.ValueBool(),
		IsMultiUser:                data.IsMultiUser.ValueBool(),
		IsSupervised:               data.IsSupervised.ValueBool(),
		Language:                   data.Language.ValueString(),
		OrgMagic:                   data.OrgMagic.ValueString(),
		Region:                     data.Region.ValueString(),
		SkipSetupItems:             stringListWithStateSet(data.SkipSetupItems),
		SupportEmailAddress:        data.SupportEmailAddress.ValueString(),
		SupportPhoneNumber:         data.SupportPhoneNumber.ValueString(),
		IncludeTlsCertificates:     data.IncludeTlsCertificates.ValueBool(),
		IosMaxVersion:              data.IosMaxVersion.ValueString(),
		IosMinVersion:              data.IosMinVersion.ValueString(),
		MacosMaxVersion:            data.MacosMaxVersion.ValueString(),
		MacosMinVersion:            data.MacosMinVersion.ValueString(),
		PushCertificateID:          int(data.PushCertificateID.ValueInt64()),
		ACMEIssuerUUID:             optionalStringWithState(data.ACMEIssuerUUID),
		SCEPIssuerUUID:             data.SCEPIssuerUUID.ValueString(),
		BlueprintID:                optionalIntWithState(data.BlueprintID),
		RealmUUID:                  optionalStringWithState(data.RealmUUID),
		VirtualServer:              int(data.VirtualServerID.ValueInt64()),
	}

	return mdmDEPEnrollmentRequest
}

func mdmDEPEnrollmentUpdateRequestWithState(data mdmDEPEnrollment) *goztl.MDMDEPEnrollmentUpdateRequest {
	mdmDEPEnrollmentRequest := &goztl.MDMDEPEnrollmentUpdateRequest{
		Name:        data.Name.ValueString(),
		DisplayName: data.DisplayName.ValueString(),
		Secret: goztl.EnrollmentSecretRequest{
			MetaBusinessUnitID: int(data.MetaBusinessUnitID.ValueInt64()),
			TagIDs:             intListWithState(data.TagIDs),
			SerialNumbers:      stringListWithStateSet(data.SerialNumbers),
			UDIDs:              stringListWithStateSet(data.UDIDs),
			Quota:              optionalIntWithState(data.Quota),
		},
		UseRealmUser:               data.UseRealmUser.ValueBool(),
		UsernamePattern:            data.UsernamePattern.ValueString(),
		RealmUserIsAdmin:           data.RealmUserIsAdmin.ValueBool(),
		AdminFullName:              optionalStringWithState(data.AdminFullName),
		AdminShortName:             optionalStringWithState(data.AdminShortName),
		HiddenAdmin:                data.HiddenAdmin.ValueBool(),
		AdminPasswordComplexity:    int(data.AdminPasswordComplexity.ValueInt64()),
		AdminPasswordRotationDelay: int(data.AdminPasswordRotationDelay.ValueInt64()),
		AllowPairing:               data.AllowPairing.ValueBool(),
		AutoAdvanceSetup:           data.AutoAdvanceSetup.ValueBool(),
		AwaitDeviceConfigured:      data.AwaitDeviceConfigured.ValueBool(),
		Department:                 data.Department.ValueString(),
		IsMandatory:                data.IsMandatory.ValueBool(),
		IsMdmRemovable:             data.IsMdmRemovable.ValueBool(),
		IsMultiUser:                data.IsMultiUser.ValueBool(),
		IsSupervised:               data.IsSupervised.ValueBool(),
		Language:                   data.Language.ValueString(),
		OrgMagic:                   data.OrgMagic.ValueString(),
		Region:                     data.Region.ValueString(),
		SkipSetupItems:             stringListWithStateSet(data.SkipSetupItems),
		SupportEmailAddress:        data.SupportEmailAddress.ValueString(),
		SupportPhoneNumber:         data.SupportPhoneNumber.ValueString(),
		IncludeTlsCertificates:     data.IncludeTlsCertificates.ValueBool(),
		IosMaxVersion:              data.IosMaxVersion.ValueString(),
		IosMinVersion:              data.IosMinVersion.ValueString(),
		MacosMaxVersion:            data.MacosMaxVersion.ValueString(),
		MacosMinVersion:            data.MacosMinVersion.ValueString(),
		PushCertificateID:          int(data.PushCertificateID.ValueInt64()),
		ACMEIssuerUUID:             optionalStringWithState(data.ACMEIssuerUUID),
		SCEPIssuerUUID:             data.SCEPIssuerUUID.ValueString(),
		BlueprintID:                optionalIntWithState(data.BlueprintID),
		RealmUUID:                  optionalStringWithState(data.RealmUUID),
	}

	return mdmDEPEnrollmentRequest
}
