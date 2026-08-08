package provider

import (
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
	return munkiConfiguration{
		ID:                              types.Int64Value(int64(mc.ID)),
		Name:                            types.StringValue(mc.Name),
		Description:                     types.StringValue(mc.Description),
		InventoryAppsFullInfoShard:      types.Int64Value(int64(mc.InventoryAppsFullInfoShard)),
		PrincipalUserDetectionSources:   stringListForState(mc.PrincipalUserDetectionSources),
		PrincipalUserDetectionDomains:   stringSetForState(mc.PrincipalUserDetectionDomains),
		CollectedConditionKeys:          stringSetForState(mc.CollectedConditionKeys),
		ManagedInstallsSyncIntervalDays: types.Int64Value(int64(mc.ManagedInstallsSyncIntervalDays)),
		ScriptChecksRunIntervalSeconds:  types.Int64Value(int64(mc.ScriptChecksRunIntervalSeconds)),
		AutoReinstallIncidents:          types.BoolValue(mc.AutoReinstallIncidents),
		AutoFailedInstallIncidents:      types.BoolValue(mc.AutoFailedInstallIncidents),
		Version:                         types.Int64Value(int64(mc.Version)),
	}
}

func munkiConfigurationRequestWithState(data munkiConfiguration) *goztl.MunkiConfigurationRequest {
	return &goztl.MunkiConfigurationRequest{
		Name:                            data.Name.ValueString(),
		Description:                     data.Description.ValueString(),
		InventoryAppsFullInfoShard:      int(data.InventoryAppsFullInfoShard.ValueInt64()),
		PrincipalUserDetectionSources:   stringListWithStateList(data.PrincipalUserDetectionSources),
		PrincipalUserDetectionDomains:   stringListWithStateSet(data.PrincipalUserDetectionDomains),
		CollectedConditionKeys:          stringListWithStateSet(data.CollectedConditionKeys),
		ManagedInstallsSyncIntervalDays: int(data.ManagedInstallsSyncIntervalDays.ValueInt64()),
		ScriptChecksRunIntervalSeconds:  int(data.ScriptChecksRunIntervalSeconds.ValueInt64()),
		AutoReinstallIncidents:          data.AutoReinstallIncidents.ValueBool(),
		AutoFailedInstallIncidents:      data.AutoFailedInstallIncidents.ValueBool(),
	}
}
