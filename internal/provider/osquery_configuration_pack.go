package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type osqueryConfigurationPack struct {
	ID              types.Int64 `tfsdk:"id"`
	ConfigurationID types.Int64 `tfsdk:"configuration_id"`
	PackID          types.Int64 `tfsdk:"pack_id"`
	TagIDs          types.Set   `tfsdk:"tag_ids"`
	ExcludedTagIDs  types.Set   `tfsdk:"excluded_tag_ids"`
}

func osqueryConfigurationPackForState(ocp *goztl.OsqueryConfigurationPack) osqueryConfigurationPack {
	return osqueryConfigurationPack{
		ID:              types.Int64Value(int64(ocp.ID)),
		ConfigurationID: types.Int64Value(int64(ocp.ConfigurationID)),
		PackID:          types.Int64Value(int64(ocp.PackID)),
		TagIDs:          int64SetForState(ocp.TagIDs),
		ExcludedTagIDs:  int64SetForState(ocp.ExcludedTagIDs),
	}
}

func osqueryConfigurationPackRequestWithState(data osqueryConfigurationPack) *goztl.OsqueryConfigurationPackRequest {
	return &goztl.OsqueryConfigurationPackRequest{
		ConfigurationID: int(data.ConfigurationID.ValueInt64()),
		PackID:          int(data.PackID.ValueInt64()),
		TagIDs:          intListWithState(data.TagIDs),
		ExcludedTagIDs:  intListWithState(data.ExcludedTagIDs),
	}
}
