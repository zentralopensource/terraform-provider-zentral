package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const (
	ztlSantaMonitor  int    = 1
	ztlSantaLockdown        = 2
	tfSantaMonitor   string = "MONITOR"
	tfSantaLockdown         = "LOCKDOWN"
)

const minZentralVersionSantaEventDetail = "v2026.6"

// block notification button sources
const (
	tfSantaEventDetailLocal        = "LOCAL"
	tfSantaEventDetailInherit      = "INHERIT"
	tfSantaEventDetailVotingPortal = "VOTING_PORTAL"
	tfSantaEventDetailCustom       = "CUSTOM"
	tfSantaEventDetailNone         = "NONE"
)

// Only a custom button has a URL of its own, and only the voting portal and the custom buttons have
// a label. The server clears the attributes a configuration source does not use, so a value it
// clears would come back empty and fail the apply, and it stores them on a scoped client mode, where
// they would do nothing. defaultSource is the source of a configuration that leaves it out.
func validateSantaEventDetail(diags *diag.Diagnostics, source, url, text types.String, defaultSource string) {
	if source.IsUnknown() {
		return
	}
	effectiveSource := defaultSource
	if !source.IsNull() {
		effectiveSource = source.ValueString()
	}
	isSet := func(s types.String) bool {
		return !s.IsNull() && !s.IsUnknown() && s.ValueString() != ""
	}

	if effectiveSource == tfSantaEventDetailCustom {
		if !url.IsUnknown() && !isSet(url) {
			diags.AddAttributeError(
				path.Root("event_detail_url"),
				"Missing Attribute Value",
				fmt.Sprintf("event_detail_url is required when event_detail_source is %s.", tfSantaEventDetailCustom),
			)
		}
	} else if isSet(url) {
		diags.AddAttributeError(
			path.Root("event_detail_url"),
			"Invalid Attribute Combination",
			fmt.Sprintf("event_detail_url can only be set when event_detail_source is %s.", tfSantaEventDetailCustom),
		)
	}

	if isSet(text) && effectiveSource != tfSantaEventDetailCustom && effectiveSource != tfSantaEventDetailVotingPortal {
		diags.AddAttributeError(
			path.Root("event_detail_text"),
			"Invalid Attribute Combination",
			fmt.Sprintf(
				"event_detail_text can only be set when event_detail_source is %s or %s.",
				tfSantaEventDetailCustom, tfSantaEventDetailVotingPortal,
			),
		)
	}
}

// The source is the support marker of the three button attributes: the validation only lets a
// configuration set the URL or the label with a source other than LOCAL, which is what an older
// server does.
func checkSantaConfigurationEventDetailSupport(diags *diag.Diagnostics, planned santaConfiguration, echoed *goztl.SantaConfiguration) {
	checkStringAttributeSupport(
		diags, "event_detail_source", minZentralVersionSantaEventDetail,
		planned.EventDetailSource, echoed.EventDetailSource, tfSantaEventDetailLocal,
	)
}

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
	EventDetailSource         types.String `tfsdk:"event_detail_source"`
	EventDetailURL            types.String `tfsdk:"event_detail_url"`
	EventDetailText           types.String `tfsdk:"event_detail_text"`
	BlockUSBMount             types.Bool   `tfsdk:"block_usb_mount"`
	RemountUSBMode            types.Set    `tfsdk:"remount_usb_mode"`
	AllowUnknownShard         types.Int64  `tfsdk:"allow_unknown_shard"`
	EnableAllEventUploadShard types.Int64  `tfsdk:"enable_all_event_upload_shard"`
	SyncIncidentSeverity      types.Int64  `tfsdk:"sync_incident_severity"`
}

func santaConfigurationForState(sc *goztl.SantaConfiguration) santaConfiguration {

	clientMode := tfSantaMonitor // default to MONITOR
	if sc.ClientMode == ztlSantaLockdown {
		clientMode = tfSantaLockdown
	}

	// a server predating the button sends none, which is what the local configuration source does
	eventDetailSource := sc.EventDetailSource
	if eventDetailSource == "" {
		eventDetailSource = tfSantaEventDetailLocal
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
		EventDetailSource:         types.StringValue(eventDetailSource),
		EventDetailURL:            types.StringValue(sc.EventDetailURL),
		EventDetailText:           types.StringValue(sc.EventDetailText),
		BlockUSBMount:             types.BoolValue(sc.BlockUSBMount),
		RemountUSBMode:            stringSetForState(sc.RemountUSBMode),
		AllowUnknownShard:         types.Int64Value(int64(sc.AllowUnknownShard)),
		EnableAllEventUploadShard: types.Int64Value(int64(sc.EnableAllEventUploadShard)),
		SyncIncidentSeverity:      types.Int64Value(int64(sc.SyncIncidentSeverity)),
	}
}

func santaConfigurationRequestWithState(data santaConfiguration) *goztl.SantaConfigurationRequest {
	clientMode := ztlSantaMonitor // default to MONITOR
	if data.ClientMode.ValueString() == tfSantaLockdown {
		clientMode = ztlSantaLockdown
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
		EventDetailSource:         data.EventDetailSource.ValueString(),
		EventDetailURL:            data.EventDetailURL.ValueString(),
		EventDetailText:           data.EventDetailText.ValueString(),
		BlockUSBMount:             data.BlockUSBMount.ValueBool(),
		RemountUSBMode:            stringListWithStateSet(data.RemountUSBMode),
		AllowUnknownShard:         int(data.AllowUnknownShard.ValueInt64()),
		EnableAllEventUploadShard: int(data.EnableAllEventUploadShard.ValueInt64()),
		SyncIncidentSeverity:      int(data.SyncIncidentSeverity.ValueInt64()),
	}
}
