package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const (
	ztlClientModeMonitor  int    = 1
	ztlClientModeLockdown        = 2
	tfClientModeMonitor   string = "MONITOR"
	tfClientModeLockdown         = "LOCKDOWN"
)

type santaConfiguration struct {
	ID                        types.Int64  `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	ClientMode                types.String `tfsdk:"client_mode"`
	ClientCertificateAuth     types.Bool   `tfsdk:"client_certificate_auth"`
	BatchSize                 types.Int64  `tfsdk:"batch_size"`
	FullSyncInterval          types.Int64  `tfsdk:"full_sync_interval"`
	EnableBundles             types.Bool   `tfsdk:"enable_bundles"`
	EnableTransitiveRules     types.Bool   `tfsdk:"enable_transitive_rules"`
	AllowedPathRegex          types.String `tfsdk:"allowed_path_regex"`
	BlockedPathRegex          types.String `tfsdk:"blocked_path_regex"`
	BlockUSBMount             types.Bool   `tfsdk:"block_usb_mount"`
	RemountUSBMode            types.Set    `tfsdk:"remount_usb_mode"`
	AllowUnknownShard         types.Int64  `tfsdk:"allow_unknown_shard"`
	EnableAllEventUploadShard types.Int64  `tfsdk:"enable_all_event_upload_shard"`
	SyncIncidentSeverity      types.Int64  `tfsdk:"sync_incident_severity"`
}

func santaConfigurationForState(sc *goztl.SantaConfiguration) santaConfiguration {

	remountUSBModes := make([]attr.Value, 0)
	for _, rumv := range sc.RemountUSBMode {
		remountUSBModes = append(remountUSBModes, types.StringValue(rumv))
	}

	clientMode := tfClientModeMonitor // default to MONITOR
	if sc.ClientMode == ztlClientModeLockdown {
		clientMode = tfClientModeLockdown
	}

	return santaConfiguration{
		ID:                        types.Int64Value(int64(sc.ID)),
		Name:                      types.StringValue(sc.Name),
		ClientMode:                types.StringValue(clientMode),
		ClientCertificateAuth:     types.BoolValue(sc.ClientCertificateAuth),
		BatchSize:                 types.Int64Value(int64(sc.BatchSize)),
		FullSyncInterval:          types.Int64Value(int64(sc.FullSyncInterval)),
		EnableBundles:             types.BoolValue(sc.EnableBundles),
		EnableTransitiveRules:     types.BoolValue(sc.EnableTransitiveRules),
		AllowedPathRegex:          types.StringValue(sc.AllowedPathRegex),
		BlockedPathRegex:          types.StringValue(sc.BlockedPathRegex),
		BlockUSBMount:             types.BoolValue(sc.BlockUSBMount),
		RemountUSBMode:            types.SetValueMust(types.StringType, remountUSBModes),
		AllowUnknownShard:         types.Int64Value(int64(sc.AllowUnknownShard)),
		EnableAllEventUploadShard: types.Int64Value(int64(sc.EnableAllEventUploadShard)),
		SyncIncidentSeverity:      types.Int64Value(int64(sc.SyncIncidentSeverity)),
	}
}

func santaConfigurationRequestWithState(data santaConfiguration) *goztl.SantaConfigurationRequest {
	remountUSBMode := make([]string, 0)
	for _, rumv := range data.RemountUSBMode.Elements() { // nil if null or unknown → no iterations
		remountUSBMode = append(remountUSBMode, rumv.(types.String).ValueString())
	}

	clientMode := ztlClientModeMonitor // default to MONITOR
	if data.ClientMode.ValueString() == tfClientModeLockdown {
		clientMode = ztlClientModeLockdown
	}

	return &goztl.SantaConfigurationRequest{
		Name:                      data.Name.ValueString(),
		ClientMode:                clientMode,
		ClientCertificateAuth:     data.ClientCertificateAuth.ValueBool(),
		BatchSize:                 int(data.BatchSize.ValueInt64()),
		FullSyncInterval:          int(data.FullSyncInterval.ValueInt64()),
		EnableBundles:             data.EnableBundles.ValueBool(),
		EnableTransitiveRules:     data.EnableTransitiveRules.ValueBool(),
		AllowedPathRegex:          data.AllowedPathRegex.ValueString(),
		BlockedPathRegex:          data.BlockedPathRegex.ValueString(),
		BlockUSBMount:             data.BlockUSBMount.ValueBool(),
		RemountUSBMode:            remountUSBMode,
		AllowUnknownShard:         int(data.AllowUnknownShard.ValueInt64()),
		EnableAllEventUploadShard: int(data.EnableAllEventUploadShard.ValueInt64()),
		SyncIncidentSeverity:      int(data.SyncIncidentSeverity.ValueInt64()),
	}
}
