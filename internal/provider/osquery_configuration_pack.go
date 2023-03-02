package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type osqueryConfigurationPack struct {
	ID              types.Int64 `tfsdk:"id"`
	ConfigurationID types.Int64 `tfsdk:"configuration_id"`
	PackID          types.Int64 `tfsdk:"pack_id"`
	TagIDs          types.Set   `tfsdk:"tag_ids"`
}

func osqueryConfigurationPackForState(ocp *goztl.OsqueryConfigurationPack) osqueryConfigurationPack {
	tagIDs := make([]attr.Value, 0)
	for _, tagID := range ocp.TagIDs {
		tagIDs = append(tagIDs, types.Int64Value(int64(tagID)))
	}

	return osqueryConfigurationPack{
		ID:              types.Int64Value(int64(ocp.ID)),
		ConfigurationID: types.Int64Value(int64(ocp.ConfigurationID)),
		PackID:          types.Int64Value(int64(ocp.PackID)),
		TagIDs:          types.SetValueMust(types.Int64Type, tagIDs),
	}
}

func osqueryConfigurationPackRequestWithState(data osqueryConfigurationPack) *goztl.OsqueryConfigurationPackRequest {
	tagIDs := make([]int, 0)
	for _, tagID := range data.TagIDs.Elements() { // nil if null or unknown → no iterations
		tagIDs = append(tagIDs, int(tagID.(types.Int64).ValueInt64()))
	}

	return &goztl.OsqueryConfigurationPackRequest{
		ConfigurationID: int(data.ConfigurationID.ValueInt64()),
		PackID:          int(data.PackID.ValueInt64()),
		TagIDs:          tagIDs,
	}
}
