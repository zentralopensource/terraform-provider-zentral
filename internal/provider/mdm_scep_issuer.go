package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const (
	tfSCEPIssuerDefaultKeySize  int64 = 2048
	tfSCEPIssuerDefaultKeyUsage       = 0
)

type mdmSCEPIssuer struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	URL             types.String `tfsdk:"url"`
	KeySize         types.Int64  `tfsdk:"key_size"`
	KeyUsage        types.Int64  `tfsdk:"key_usage"`
	Backend         types.String `tfsdk:"backend"`
	IDent           types.Object `tfsdk:"ident"`
	MicrosoftCA     types.Object `tfsdk:"microsoft_ca"`
	OktaCA          types.Object `tfsdk:"okta_ca"`
	StaticChallenge types.Object `tfsdk:"static_challenge"`
}

func mdmSCEPIssuerForState(msi *goztl.MDMSCEPIssuer) mdmSCEPIssuer {
	return mdmSCEPIssuer{
		ID:              types.StringValue(msi.ID),
		Name:            types.StringValue(msi.Name),
		Description:     types.StringValue(msi.Description),
		URL:             types.StringValue(msi.URL),
		KeySize:         types.Int64Value(int64(msi.KeySize)),
		KeyUsage:        types.Int64Value(int64(msi.KeyUsage)),
		Backend:         optionalStringForState(msi.Backend),
		IDent:           identBackendForState(msi.IDent),
		MicrosoftCA:     microsoftCABackendForState(msi.MicrosoftCA),
		OktaCA:          microsoftCABackendForState(msi.OktaCA),
		StaticChallenge: staticChallengeBackendForState(msi.StaticChallenge),
	}
}

func mdmSCEPIssuerRequestWithState(data mdmSCEPIssuer) *goztl.MDMSCEPIssuerRequest {
	return &goztl.MDMSCEPIssuerRequest{
		Name:            data.Name.ValueString(),
		Description:     data.Description.ValueString(),
		URL:             data.URL.ValueString(),
		KeySize:         int(data.KeySize.ValueInt64()),
		KeyUsage:        int(data.KeyUsage.ValueInt64()),
		Backend:         data.Backend.ValueString(),
		IDent:           identBackendWithState(data.IDent),
		MicrosoftCA:     microsoftCABackendWithState(data.MicrosoftCA),
		OktaCA:          microsoftCABackendWithState(data.OktaCA),
		StaticChallenge: staticChallengeBackendWithState(data.StaticChallenge),
	}
}
