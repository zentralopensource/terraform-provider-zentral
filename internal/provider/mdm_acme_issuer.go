package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const (
	tfACMEIssuerRSAKeyType              string = "RSA"
	tfACMEIssuerECSECPrimeRandomKeyType        = "ECSECPrimeRandom"
	tfACMEIssuerDefaultKeyUsage         int64  = 0
)

type mdmACMEIssuer struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	DirectoryURL     types.String `tfsdk:"directory_url"`
	KeyType          types.String `tfsdk:"key_type"`
	KeySize          types.Int64  `tfsdk:"key_size"`
	UsageFlags       types.Int64  `tfsdk:"usage_flags"`
	ExtendedKeyUsage types.Set    `tfsdk:"extended_key_usage"`
	HardwareBound    types.Bool   `tfsdk:"hardware_bound"`
	Attest           types.Bool   `tfsdk:"attest"`
	Backend          types.String `tfsdk:"backend"`
	IDent            types.Object `tfsdk:"ident"`
	MicrosoftCA      types.Object `tfsdk:"microsoft_ca"`
	OktaCA           types.Object `tfsdk:"okta_ca"`
	StaticChallenge  types.Object `tfsdk:"static_challenge"`
}

func mdmACMEIssuerForState(mai *goztl.MDMACMEIssuer) mdmACMEIssuer {
	return mdmACMEIssuer{
		ID:               types.StringValue(mai.ID),
		Name:             types.StringValue(mai.Name),
		Description:      types.StringValue(mai.Description),
		DirectoryURL:     types.StringValue(mai.DirectoryURL),
		KeyType:          types.StringValue(mai.KeyType),
		KeySize:          types.Int64Value(int64(mai.KeySize)),
		UsageFlags:       types.Int64Value(int64(mai.UsageFlags)),
		ExtendedKeyUsage: stringSetForState(mai.ExtendedKeyUsage),
		HardwareBound:    types.BoolValue(mai.HardwareBound),
		Attest:           types.BoolValue(mai.Attest),
		Backend:          optionalStringForState(mai.Backend),
		IDent:            identBackendForState(mai.IDent),
		MicrosoftCA:      microsoftCABackendForState(mai.MicrosoftCA),
		OktaCA:           microsoftCABackendForState(mai.OktaCA),
		StaticChallenge:  staticChallengeBackendForState(mai.StaticChallenge),
	}
}

func mdmACMEIssuerRequestWithState(data mdmACMEIssuer) *goztl.MDMACMEIssuerRequest {
	return &goztl.MDMACMEIssuerRequest{
		Name:             data.Name.ValueString(),
		Description:      data.Description.ValueString(),
		DirectoryURL:     data.DirectoryURL.ValueString(),
		KeyType:          data.KeyType.ValueString(),
		KeySize:          int(data.KeySize.ValueInt64()),
		UsageFlags:       int(data.UsageFlags.ValueInt64()),
		ExtendedKeyUsage: stringListWithStateSet(data.ExtendedKeyUsage),
		HardwareBound:    data.HardwareBound.ValueBool(),
		Attest:           data.Attest.ValueBool(),
		Backend:          data.Backend.ValueString(),
		IDent:            identBackendWithState(data.IDent),
		MicrosoftCA:      microsoftCABackendWithState(data.MicrosoftCA),
		OktaCA:           microsoftCABackendWithState(data.OktaCA),
		StaticChallenge:  staticChallengeBackendWithState(data.StaticChallenge),
	}
}
