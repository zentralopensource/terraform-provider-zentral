package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmSCEPConfig struct {
	ID                 types.Int64  `tfsdk:"id"`
	ProvisioningUID    types.String `tfsdk:"provisioning_uid"`
	Name               types.String `tfsdk:"name"`
	URL                types.String `tfsdk:"url"`
	KeyUsage           types.Int64  `tfsdk:"key_usage"`
	KeyIsExtractable   types.Bool   `tfsdk:"key_is_extractable"`
	KeySize            types.Int64  `tfsdk:"key_size"`
	AllowAllAppsAccess types.Bool   `tfsdk:"allow_all_apps_access"`
}

func mdmSCEPConfigForState(msc *goztl.MDMSCEPConfig) mdmSCEPConfig {
	var provisioningUID types.String
	if msc.ProvisioningUID != nil {
		provisioningUID = types.StringValue(*msc.ProvisioningUID)
	} else {
		provisioningUID = types.StringNull()
	}

	return mdmSCEPConfig{
		ID:                 types.Int64Value(int64(msc.ID)),
		ProvisioningUID:    provisioningUID,
		Name:               types.StringValue(msc.Name),
		URL:                types.StringValue(msc.URL),
		KeyUsage:           types.Int64Value(int64(msc.KeyUsage)),
		KeyIsExtractable:   types.BoolValue(msc.KeyIsExtractable),
		KeySize:            types.Int64Value(int64(msc.KeySize)),
		AllowAllAppsAccess: types.BoolValue(msc.AllowAllAppsAccess),
	}
}
