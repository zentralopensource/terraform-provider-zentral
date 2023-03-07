package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type monolithSubManifestPkgInfo struct {
	ID             types.Int64  `tfsdk:"id"`
	SubManifestID  types.Int64  `tfsdk:"sub_manifest_id"`
	Key            types.String `tfsdk:"key"`
	PkgInfoName    types.String `tfsdk:"pkg_info_name"`
	FeaturedItem   types.Bool   `tfsdk:"featured_item"`
	ConditionID    types.Int64  `tfsdk:"condition_id"`
	ShardModulo    types.Int64  `tfsdk:"shard_modulo"`
	DefaultShard   types.Int64  `tfsdk:"default_shard"`
	ExcludedTagIDs types.Set    `tfsdk:"excluded_tag_ids"`
	TagShards      types.Set    `tfsdk:"tag_shards"`
}

var tagShardAttrTypes = map[string]attr.Type{
	"tag_id": types.Int64Type,
	"shard":  types.Int64Type,
}

func monolithSubManifestPkgInfoForState(msmpi *goztl.MonolithSubManifestPkgInfo) monolithSubManifestPkgInfo {
	var cID types.Int64
	if msmpi.ConditionID != nil {
		cID = types.Int64Value(int64(*msmpi.ConditionID))
	} else {
		cID = types.Int64Null()
	}

	exTagIDs := make([]attr.Value, 0)
	for _, exTagID := range msmpi.ExcludedTagIDs {
		exTagIDs = append(exTagIDs, types.Int64Value(int64(exTagID)))
	}

	tagShards := make([]attr.Value, 0)
	for _, tagShard := range msmpi.TagShards {
		tagShards = append(
			tagShards,
			types.ObjectValueMust(
				tagShardAttrTypes,
				map[string]attr.Value{
					"tag_id": types.Int64Value(int64(tagShard.TagID)),
					"shard":  types.Int64Value(int64(tagShard.Shard)),
				},
			),
		)
	}

	return monolithSubManifestPkgInfo{
		ID:             types.Int64Value(int64(msmpi.ID)),
		SubManifestID:  types.Int64Value(int64(msmpi.SubManifestID)),
		Key:            types.StringValue(msmpi.Key),
		PkgInfoName:    types.StringValue(msmpi.PkgInfoName),
		FeaturedItem:   types.BoolValue(msmpi.FeaturedItem),
		ConditionID:    cID,
		ShardModulo:    types.Int64Value(int64(msmpi.ShardModulo)),
		DefaultShard:   types.Int64Value(int64(msmpi.DefaultShard)),
		ExcludedTagIDs: types.SetValueMust(types.Int64Type, exTagIDs),
		TagShards:      types.SetValueMust(types.ObjectType{AttrTypes: tagShardAttrTypes}, tagShards),
	}
}

func monolithSubManifestPkgInfoRequestWithState(data monolithSubManifestPkgInfo) *goztl.MonolithSubManifestPkgInfoRequest {
	var cID *int
	if !data.ConditionID.IsNull() {
		cID = goztl.Int(int(data.ConditionID.ValueInt64()))
	}

	exTagIDs := make([]int, 0)
	for _, exTagID := range data.ExcludedTagIDs.Elements() { // nil if null or unknown → no iterations
		exTagIDs = append(exTagIDs, int(exTagID.(types.Int64).ValueInt64()))
	}

	tagShards := make([]goztl.TagShard, 0)
	for _, tagShard := range data.TagShards.Elements() { // nil if null or unknown → no iterations
		tagShardMap := tagShard.(types.Object).Attributes()
		if tagShardMap != nil {
			tagShards = append(
				tagShards,
				goztl.TagShard{
					TagID: int(tagShardMap["tag_id"].(types.Int64).ValueInt64()),
					Shard: int(tagShardMap["shard"].(types.Int64).ValueInt64()),
				},
			)
		}
	}

	return &goztl.MonolithSubManifestPkgInfoRequest{
		SubManifestID:  int(data.SubManifestID.ValueInt64()),
		Key:            data.Key.ValueString(),
		PkgInfoName:    data.PkgInfoName.ValueString(),
		FeaturedItem:   data.FeaturedItem.ValueBool(),
		ConditionID:    cID,
		ShardModulo:    int(data.ShardModulo.ValueInt64()),
		DefaultShard:   int(data.DefaultShard.ValueInt64()),
		ExcludedTagIDs: exTagIDs,
		TagShards:      tagShards,
	}
}
