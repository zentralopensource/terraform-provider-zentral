package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	platforms := make([]attr.Value, 0)
	for _, p := range ma.Platforms {
		platforms = append(platforms, types.StringValue(p))
	}

	requires := make([]attr.Value, 0)
	for _, raID := range ma.Requires {
		requires = append(requires, types.StringValue(raID))
	}

	return mdmArtifact{
		ID:                          types.StringValue(ma.ID),
		Name:                        types.StringValue(ma.Name),
		Type:                        types.StringValue(ma.Type),
		Channel:                     types.StringValue(ma.Channel),
		Platforms:                   types.SetValueMust(types.StringType, platforms),
		InstallDuringSetupAssistant: types.BoolValue(ma.InstallDuringSetupAssistant),
		AutoUpdate:                  types.BoolValue(ma.AutoUpdate),
		ReinstallInterval:           types.Int64Value(int64(ma.ReinstallInterval)),
		ReinstallOnOSUpdate:         types.StringValue(ma.ReinstallOnOSUpdate),
		Requires:                    types.SetValueMust(types.StringType, requires),
	}
}

func mdmArtifactRequestWithState(data mdmArtifact) *goztl.MDMArtifactRequest {
	platforms := make([]string, 0)
	for _, p := range data.Platforms.Elements() { // nil if null or unknown → no iterations
		platforms = append(platforms, p.(types.String).ValueString())
	}

	requires := make([]string, 0)
	for _, raID := range data.Requires.Elements() { // nil if null or unknown → no iterations
		requires = append(requires, raID.(types.String).ValueString())
	}

	return &goztl.MDMArtifactRequest{
		Name:                        data.Name.ValueString(),
		Type:                        data.Type.ValueString(),
		Channel:                     data.Channel.ValueString(),
		Platforms:                   platforms,
		InstallDuringSetupAssistant: data.InstallDuringSetupAssistant.ValueBool(),
		AutoUpdate:                  data.AutoUpdate.ValueBool(),
		ReinstallInterval:           int(data.ReinstallInterval.ValueInt64()),
		ReinstallOnOSUpdate:         data.ReinstallOnOSUpdate.ValueString(),
		Requires:                    requires,
	}
}
