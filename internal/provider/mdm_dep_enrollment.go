package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const (
	tfMDMDEPEnrollmentDefaultAdminPasswordComplexity    int64 = 3
	tfMDMDEPEnrollmentMinAdminPasswordComplexity              = 1
	tfMDMDEPEnrollmentMaxAdminPasswordComplexity              = 3
	tfMDMDEPEnrollmentDefaultAdminPasswordRotationDelay       = 60
	tfMDMDEPEnrollmentMinAdminPasswordRotationDelay           = 0
	tfMDMDEPEnrollmentMaxAdminPasswordRotationDelay           = 1440
)

// Authentication
var mdeAuthenticationAttrTypes = map[string]attr.Type{
	"realm_uuid":                       types.StringType,
	"use_for_setup_assistant_user":     types.BoolType,
	"setup_assistant_user_is_admin":    types.BoolType,
	"setup_assistant_username_pattern": types.StringType,
}

func mdeAuthenticationDefault() types.Object {
	return types.ObjectValueMust(
		mdeAuthenticationAttrTypes,
		map[string]attr.Value{
			"realm_uuid":                       types.StringNull(),
			"use_for_setup_assistant_user":     types.BoolValue(false),
			"setup_assistant_user_is_admin":    types.BoolValue(false), // different in the Zentral model, but that would block
			"setup_assistant_username_pattern": types.StringValue(""),
		},
	)
}

func mdeAuthenticationForState(mde *goztl.MDMDEPEnrollment) types.Object {
	return types.ObjectValueMust(
		mdeAuthenticationAttrTypes,
		map[string]attr.Value{
			"realm_uuid":                       optionalStringForState(mde.RealmUUID),
			"use_for_setup_assistant_user":     types.BoolValue(mde.UseRealmUser),
			"setup_assistant_user_is_admin":    types.BoolValue(mde.RealmUserIsAdmin),
			"setup_assistant_username_pattern": types.StringValue(mde.UsernamePattern),
		},
	)
}

func addAuthenticationToReqWithState(data mdmDEPEnrollment, mder *goztl.MDMDEPEnrollmentRequest) {
	oMap := data.Authentication.Attributes()
	mder.RealmUUID = optionalStringWithState(oMap["realm_uuid"].(types.String))
	mder.UseRealmUser = oMap["use_for_setup_assistant_user"].(types.Bool).ValueBool()
	mder.RealmUserIsAdmin = oMap["setup_assistant_user_is_admin"].(types.Bool).ValueBool()
	mder.UsernamePattern = oMap["setup_assistant_username_pattern"].(types.String).ValueString()
}

// Extra admin
var mdeExtraAdminAttrTypes = map[string]attr.Type{
	"hidden":                  types.BoolType,
	"full_name":               types.StringType,
	"short_name":              types.StringType,
	"password_complexity":     types.Int64Type,
	"password_rotation_delay": types.Int64Type,
}

func mdeExtraAdminDefault() types.Object {
	return types.ObjectValueMust(
		mdeExtraAdminAttrTypes,
		map[string]attr.Value{
			"hidden":                  types.BoolValue(false),
			"full_name":               types.StringNull(),
			"short_name":              types.StringNull(),
			"password_complexity":     types.Int64Value(tfMDMDEPEnrollmentDefaultAdminPasswordComplexity),
			"password_rotation_delay": types.Int64Value(tfMDMDEPEnrollmentDefaultAdminPasswordRotationDelay),
		},
	)
}

func mdeExtraAdminForState(mde *goztl.MDMDEPEnrollment) types.Object {
	return types.ObjectValueMust(
		mdeExtraAdminAttrTypes,
		map[string]attr.Value{
			"hidden":                  types.BoolValue(mde.HiddenAdmin),
			"full_name":               optionalStringForState(mde.AdminFullName),
			"short_name":              optionalStringForState(mde.AdminShortName),
			"password_complexity":     types.Int64Value(int64(mde.AdminPasswordComplexity)),
			"password_rotation_delay": types.Int64Value(int64(mde.AdminPasswordRotationDelay)),
		},
	)
}

func addExtraAdminToReqWithState(data mdmDEPEnrollment, mder *goztl.MDMDEPEnrollmentRequest) {
	oMap := data.ExtraAdmin.Attributes()
	mder.HiddenAdmin = oMap["hidden"].(types.Bool).ValueBool()
	mder.AdminFullName = optionalStringWithState(oMap["full_name"].(types.String))
	mder.AdminShortName = optionalStringWithState(oMap["short_name"].(types.String))
	mder.AdminPasswordComplexity = int(oMap["password_complexity"].(types.Int64).ValueInt64())
	mder.AdminPasswordRotationDelay = int(oMap["password_rotation_delay"].(types.Int64).ValueInt64())
}

// Profile
var mdeProfileAttrTypes = map[string]attr.Type{
	"virtual_server_id": types.Int64Type,
	"name":              types.StringType,

	"allow_pairing":           types.BoolType,
	"auto_advance_setup":      types.BoolType,
	"await_device_configured": types.BoolType,

	"is_mandatory":     types.BoolType,
	"is_mdm_removable": types.BoolType,
	"is_multi_user":    types.BoolType,
	"is_supervised":    types.BoolType,

	"include_anchor_certs": types.BoolType,

	"skip_setup_items": types.SetType{ElemType: types.StringType},

	"language": types.StringType,
	"region":   types.StringType,

	"department": types.StringType,
	"org_magic":  types.StringType,

	"support_email_address": types.StringType,
	"support_phone_number":  types.StringType,
}

func mdeProfileForState(mde *goztl.MDMDEPEnrollment) types.Object {
	return types.ObjectValueMust(
		mdeProfileAttrTypes,
		map[string]attr.Value{
			"virtual_server_id":       types.Int64Value(int64(mde.VirtualServerID)),
			"name":                    types.StringValue(mde.Name),
			"allow_pairing":           types.BoolValue(mde.AllowPairing),
			"auto_advance_setup":      types.BoolValue(mde.AutoAdvanceSetup),
			"await_device_configured": types.BoolValue(mde.AwaitDeviceConfigured),
			"is_mandatory":            types.BoolValue(mde.IsMandatory),
			"is_mdm_removable":        types.BoolValue(mde.IsMDMRemovable),
			"is_multi_user":           types.BoolValue(mde.IsMultiUser),
			"is_supervised":           types.BoolValue(mde.IsSupervised),
			"include_anchor_certs":    types.BoolValue(mde.IncludeTLSCertificates),
			"skip_setup_items":        stringSetForState(mde.SkipSetupItems),
			"language":                types.StringValue(mde.Language),
			"region":                  types.StringValue(mde.Region),
			"department":              types.StringValue(mde.Department),
			"org_magic":               types.StringValue(mde.OrgMagic),
			"support_email_address":   types.StringValue(mde.SupportEmailAddress),
			"support_phone_number":    types.StringValue(mde.SupportPhoneNumber),
		},
	)
}

func addProfileToReqWithState(data mdmDEPEnrollment, mder *goztl.MDMDEPEnrollmentRequest) {
	oMap := data.Profile.Attributes()
	mder.VirtualServerID = int(oMap["virtual_server_id"].(types.Int64).ValueInt64())
	mder.Name = oMap["name"].(types.String).ValueString()
	mder.AllowPairing = oMap["allow_pairing"].(types.Bool).ValueBool()
	mder.AutoAdvanceSetup = oMap["auto_advance_setup"].(types.Bool).ValueBool()
	mder.AwaitDeviceConfigured = oMap["await_device_configured"].(types.Bool).ValueBool()
	mder.IsMandatory = oMap["is_mandatory"].(types.Bool).ValueBool()
	mder.IsMDMRemovable = oMap["is_mdm_removable"].(types.Bool).ValueBool()
	mder.IsMultiUser = oMap["is_multi_user"].(types.Bool).ValueBool()
	mder.IsSupervised = oMap["is_supervised"].(types.Bool).ValueBool()
	mder.IncludeTLSCertificates = oMap["include_anchor_certs"].(types.Bool).ValueBool()
	mder.SkipSetupItems = stringListWithStateSet(oMap["skip_setup_items"].(types.Set))
	mder.Language = oMap["language"].(types.String).ValueString()
	mder.Region = oMap["region"].(types.String).ValueString()
	mder.Department = oMap["department"].(types.String).ValueString()
	mder.OrgMagic = oMap["org_magic"].(types.String).ValueString()
	mder.SupportEmailAddress = oMap["support_email_address"].(types.String).ValueString()
	mder.SupportPhoneNumber = oMap["support_phone_number"].(types.String).ValueString()
}

// OS version enforcement
var mdeOSVersionEnforcementAttrTypes = map[string]attr.Type{
	"ios_min_version":              types.StringType,
	"auto_ios_min_version_until":   types.StringType,
	"macos_min_version":            types.StringType,
	"auto_macos_min_version_until": types.StringType,
}

func mdeOSVersionEnforcementDefault() types.Object {
	return types.ObjectValueMust(
		mdeOSVersionEnforcementAttrTypes,
		map[string]attr.Value{
			"ios_min_version":              types.StringValue(""),
			"auto_ios_min_version_until":   types.StringValue(""),
			"macos_min_version":            types.StringValue(""),
			"auto_macos_min_version_until": types.StringValue(""),
		},
	)
}

func mdeOSVersionEnforcementForState(mde *goztl.MDMDEPEnrollment) types.Object {
	return types.ObjectValueMust(
		mdeOSVersionEnforcementAttrTypes,
		map[string]attr.Value{
			"ios_min_version":              types.StringValue(mde.IOSMinVersion),
			"auto_ios_min_version_until":   types.StringValue(mde.IOSMaxVersion),
			"macos_min_version":            types.StringValue(mde.MacOSMinVersion),
			"auto_macos_min_version_until": types.StringValue(mde.MacOSMaxVersion),
		},
	)
}

func addOSVersionEnforcementToReqWithState(data mdmDEPEnrollment, mder *goztl.MDMDEPEnrollmentRequest) {
	oMap := data.OSVersionEnforcement.Attributes()
	mder.IOSMinVersion = oMap["ios_min_version"].(types.String).ValueString()
	mder.IOSMaxVersion = oMap["auto_ios_min_version_until"].(types.String).ValueString()
	mder.MacOSMinVersion = oMap["macos_min_version"].(types.String).ValueString()
	mder.MacOSMaxVersion = oMap["auto_macos_min_version_until"].(types.String).ValueString()
}

// MDM DEP enrollment
type mdmDEPEnrollment struct {
	ID          types.Int64  `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`

	PushCertificateID types.Int64 `tfsdk:"push_certificate_id"`
	BlueprintID       types.Int64 `tfsdk:"blueprint_id"`

	ACMEIssuerUUID types.String `tfsdk:"acme_issuer_id"`
	SCEPIssuerUUID types.String `tfsdk:"scep_issuer_id"`

	EnrollmentSecret     types.Object `tfsdk:"enrollment"`
	Authentication       types.Object `tfsdk:"authentication"`
	ExtraAdmin           types.Object `tfsdk:"extra_admin"`
	Profile              types.Object `tfsdk:"profile"`
	OSVersionEnforcement types.Object `tfsdk:"os_version_enforcement"`
}

func mdmDEPEnrollmentForState(mde *goztl.MDMDEPEnrollment) mdmDEPEnrollment {
	return mdmDEPEnrollment{
		ID:          types.Int64Value(int64(mde.ID)),
		DisplayName: types.StringValue(mde.DisplayName),

		PushCertificateID: types.Int64Value(int64(mde.PushCertificateID)),
		BlueprintID:       optionalInt64ForState(mde.BlueprintID),

		ACMEIssuerUUID: optionalStringForState(mde.ACMEIssuerUUID),
		SCEPIssuerUUID: types.StringValue(mde.SCEPIssuerUUID),

		EnrollmentSecret:     enrollmentSecretForState(mde.Secret),
		Authentication:       mdeAuthenticationForState(mde),
		ExtraAdmin:           mdeExtraAdminForState(mde),
		Profile:              mdeProfileForState(mde),
		OSVersionEnforcement: mdeOSVersionEnforcementForState(mde),
	}
}

func mdmDEPEnrollmentRequestWithState(data mdmDEPEnrollment) *goztl.MDMDEPEnrollmentRequest {
	mdmDEPEnrollmentRequest := &goztl.MDMDEPEnrollmentRequest{
		DisplayName: data.DisplayName.ValueString(),

		PushCertificateID: int(data.PushCertificateID.ValueInt64()),
		BlueprintID:       optionalIntWithState(data.BlueprintID),

		ACMEIssuerUUID: optionalStringWithState(data.ACMEIssuerUUID),
		SCEPIssuerUUID: data.SCEPIssuerUUID.ValueString(),

		Secret: enrollmentSecretRequestWithState(data.EnrollmentSecret),
	}

	addAuthenticationToReqWithState(data, mdmDEPEnrollmentRequest)
	addExtraAdminToReqWithState(data, mdmDEPEnrollmentRequest)
	addProfileToReqWithState(data, mdmDEPEnrollmentRequest)
	addOSVersionEnforcementToReqWithState(data, mdmDEPEnrollmentRequest)

	return mdmDEPEnrollmentRequest
}
