package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmArtifact struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	Type                        types.String `tfsdk:"type"`
	Channel                     types.String `tfsdk:"channel"`
	Platforms                   types.Set    `tfsdk:"platforms"`
	InstallDuringSetupAssistant types.Bool   `tfsdk:"install_during_setup_assistant"`
	AutoUpdate                  types.Bool   `tfsdk:"auto_update"`
	ReinstallInterval           types.Int64  `tfsdk:"reinstall_interval"`
	ReinstallOnOSUpdate         types.String `tfsdk:"reinstall_on_os_update"`
	Requires                    types.Set    `tfsdk:"requires"`
}

func mdmArtifactForState(ma *goztl.MDMArtifact) mdmArtifact {
	return mdmArtifact{
		ID:                          types.StringValue(ma.ID),
		Name:                        types.StringValue(ma.Name),
		Type:                        types.StringValue(ma.Type),
		Channel:                     types.StringValue(ma.Channel),
		Platforms:                   stringSetForState(ma.Platforms),
		InstallDuringSetupAssistant: types.BoolValue(ma.InstallDuringSetupAssistant),
		AutoUpdate:                  types.BoolValue(ma.AutoUpdate),
		ReinstallInterval:           types.Int64Value(int64(ma.ReinstallInterval)),
		ReinstallOnOSUpdate:         types.StringValue(ma.ReinstallOnOSUpdate),
		Requires:                    stringSetForState(ma.Requires),
	}
}

func mdmArtifactRequestWithState(data mdmArtifact) *goztl.MDMArtifactRequest {
	return &goztl.MDMArtifactRequest{
		Name:                        data.Name.ValueString(),
		Type:                        data.Type.ValueString(),
		Channel:                     data.Channel.ValueString(),
		Platforms:                   stringListWithStateSet(data.Platforms),
		InstallDuringSetupAssistant: data.InstallDuringSetupAssistant.ValueBool(),
		AutoUpdate:                  data.AutoUpdate.ValueBool(),
		ReinstallInterval:           int(data.ReinstallInterval.ValueInt64()),
		ReinstallOnOSUpdate:         data.ReinstallOnOSUpdate.ValueString(),
		Requires:                    stringListWithStateSet(data.Requires),
	}
}
