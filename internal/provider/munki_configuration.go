package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type munkiConfiguration struct {
	ID                              types.Int64  `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	Description                     types.String `tfsdk:"description"`
	InventoryAppsFullInfoShard      types.Int64  `tfsdk:"inventory_apps_full_info_shard"`
	PrincipalUserDetectionSources   types.List   `tfsdk:"principal_user_detection_sources"`
	PrincipalUserDetectionDomains   types.Set    `tfsdk:"principal_user_detection_domains"`
	CollectedConditionKeys          types.Set    `tfsdk:"collected_condition_keys"`
	ManagedInstallsSyncIntervalDays types.Int64  `tfsdk:"managed_installs_sync_interval_days"`
	ScriptChecksRunIntervalSeconds  types.Int64  `tfsdk:"script_checks_run_interval_seconds"`
	AutoReinstallIncidents          types.Bool   `tfsdk:"auto_reinstall_incidents"`
	AutoFailedInstallIncidents      types.Bool   `tfsdk:"auto_failed_install_incidents"`
	Version                         types.Int64  `tfsdk:"version"`
}

func munkiConfigurationForState(mc *goztl.MunkiConfiguration) munkiConfiguration {
	pudss := make([]attr.Value, 0)
	for _, puds := range mc.PrincipalUserDetectionSources {
		pudss = append(pudss, types.StringValue(puds))
	}

	pudds := make([]attr.Value, 0)
	for _, pudd := range mc.PrincipalUserDetectionDomains {
		pudds = append(pudds, types.StringValue(pudd))
	}

	ccks := make([]attr.Value, 0)
	for _, cck := range mc.CollectedConditionKeys {
		ccks = append(ccks, types.StringValue(cck))
	}

	return munkiConfiguration{
		ID:                              types.Int64Value(int64(mc.ID)),
		Name:                            types.StringValue(mc.Name),
		Description:                     types.StringValue(mc.Description),
		InventoryAppsFullInfoShard:      types.Int64Value(int64(mc.InventoryAppsFullInfoShard)),
		PrincipalUserDetectionSources:   types.ListValueMust(types.StringType, pudss),
		PrincipalUserDetectionDomains:   types.SetValueMust(types.StringType, pudds),
		CollectedConditionKeys:          types.SetValueMust(types.StringType, ccks),
		ManagedInstallsSyncIntervalDays: types.Int64Value(int64(mc.ManagedInstallsSyncIntervalDays)),
		ScriptChecksRunIntervalSeconds:  types.Int64Value(int64(mc.ScriptChecksRunIntervalSeconds)),
		AutoReinstallIncidents:          types.BoolValue(mc.AutoReinstallIncidents),
		AutoFailedInstallIncidents:      types.BoolValue(mc.AutoFailedInstallIncidents),
		Version:                         types.Int64Value(int64(mc.Version)),
	}
}

func munkiConfigurationRequestWithState(data munkiConfiguration) *goztl.MunkiConfigurationRequest {
	pudss := make([]string, 0)
	for _, puds := range data.PrincipalUserDetectionSources.Elements() {
		pudss = append(pudss, puds.(types.String).ValueString())
	}

	pudds := make([]string, 0)
	for _, pudd := range data.PrincipalUserDetectionDomains.Elements() {
		pudds = append(pudds, pudd.(types.String).ValueString())
	}

	ccks := make([]string, 0)
	for _, cck := range data.CollectedConditionKeys.Elements() {
		ccks = append(ccks, cck.(types.String).ValueString())
	}

	return &goztl.MunkiConfigurationRequest{
		Name:                            data.Name.ValueString(),
		Description:                     data.Description.ValueString(),
		InventoryAppsFullInfoShard:      int(data.InventoryAppsFullInfoShard.ValueInt64()),
		PrincipalUserDetectionSources:   pudss,
		PrincipalUserDetectionDomains:   pudds,
		CollectedConditionKeys:          ccks,
		ManagedInstallsSyncIntervalDays: int(data.ManagedInstallsSyncIntervalDays.ValueInt64()),
		ScriptChecksRunIntervalSeconds:  int(data.ScriptChecksRunIntervalSeconds.ValueInt64()),
		AutoReinstallIncidents:          data.AutoReinstallIncidents.ValueBool(),
		AutoFailedInstallIncidents:      data.AutoFailedInstallIncidents.ValueBool(),
	}
}
