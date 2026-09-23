package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const minZentralVersionSantaScopedConfigurationItems = "v2026.6"

type santaScopedClientMode struct {
	ID                    types.Int64  `tfsdk:"id"`
	ConfigurationID       types.Int64  `tfsdk:"configuration_id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	ClientMode            types.String `tfsdk:"client_mode"`
	EventDetailSource     types.String `tfsdk:"event_detail_source"`
	EventDetailURL        types.String `tfsdk:"event_detail_url"`
	EventDetailText       types.String `tfsdk:"event_detail_text"`
	PrimaryUsers          types.Set    `tfsdk:"primary_users"`
	ExcludedPrimaryUsers  types.Set    `tfsdk:"excluded_primary_users"`
	SerialNumbers         types.Set    `tfsdk:"serial_numbers"`
	ExcludedSerialNumbers types.Set    `tfsdk:"excluded_serial_numbers"`
	TagIDs                types.Set    `tfsdk:"tag_ids"`
	ExcludedTagIDs        types.Set    `tfsdk:"excluded_tag_ids"`
}

func santaScopedClientModeForState(sscm *goztl.SantaScopedClientMode) santaScopedClientMode {
	clientMode := tfSantaMonitor
	if sscm.ClientMode == ztlSantaLockdown {
		clientMode = tfSantaLockdown
	}

	return santaScopedClientMode{
		ID:                    types.Int64Value(int64(sscm.ID)),
		ConfigurationID:       types.Int64Value(int64(sscm.ConfigurationID)),
		Name:                  types.StringValue(sscm.Name),
		Description:           types.StringValue(sscm.Description),
		ClientMode:            types.StringValue(clientMode),
		EventDetailSource:     types.StringValue(sscm.EventDetailSource),
		EventDetailURL:        types.StringValue(sscm.EventDetailURL),
		EventDetailText:       types.StringValue(sscm.EventDetailText),
		PrimaryUsers:          stringSetForState(sscm.PrimaryUsers),
		ExcludedPrimaryUsers:  stringSetForState(sscm.ExcludedPrimaryUsers),
		SerialNumbers:         stringSetForState(sscm.SerialNumbers),
		ExcludedSerialNumbers: stringSetForState(sscm.ExcludedSerialNumbers),
		TagIDs:                int64SetForState(sscm.TagIDs),
		ExcludedTagIDs:        int64SetForState(sscm.ExcludedTagIDs),
	}
}

func santaScopedClientModeRequestWithState(data santaScopedClientMode) *goztl.SantaScopedClientModeRequest {
	clientMode := ztlSantaMonitor
	if data.ClientMode.ValueString() == tfSantaLockdown {
		clientMode = ztlSantaLockdown
	}

	return &goztl.SantaScopedClientModeRequest{
		ConfigurationID:       int(data.ConfigurationID.ValueInt64()),
		Name:                  data.Name.ValueString(),
		Description:           data.Description.ValueString(),
		ClientMode:            clientMode,
		EventDetailSource:     data.EventDetailSource.ValueString(),
		EventDetailURL:        data.EventDetailURL.ValueString(),
		EventDetailText:       data.EventDetailText.ValueString(),
		PrimaryUsers:          stringListWithStateSet(data.PrimaryUsers),
		ExcludedPrimaryUsers:  stringListWithStateSet(data.ExcludedPrimaryUsers),
		SerialNumbers:         stringListWithStateSet(data.SerialNumbers),
		ExcludedSerialNumbers: stringListWithStateSet(data.ExcludedSerialNumbers),
		TagIDs:                intListWithState(data.TagIDs),
		ExcludedTagIDs:        intListWithState(data.ExcludedTagIDs),
	}
}
