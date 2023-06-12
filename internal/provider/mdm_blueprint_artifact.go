package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmBlueprintArtifact struct {
	ID               types.Int64  `tfsdk:"id"`
	BlueprintID      types.Int64  `tfsdk:"blueprint_id"`
	ArtifactID       types.String `tfsdk:"artifact_id"`
	IOS              types.Bool   `tfsdk:"ios"`
	IOSMaxVersion    types.String `tfsdk:"ios_max_version"`
	IOSMinVersion    types.String `tfsdk:"ios_min_version"`
	IPadOS           types.Bool   `tfsdk:"ipados"`
	IPadOSMaxVersion types.String `tfsdk:"ipados_max_version"`
	IPadOSMinVersion types.String `tfsdk:"ipados_min_version"`
	MacOS            types.Bool   `tfsdk:"macos"`
	MacOSMaxVersion  types.String `tfsdk:"macos_max_version"`
	MacOSMinVersion  types.String `tfsdk:"macos_min_version"`
	TVOS             types.Bool   `tfsdk:"tvos"`
	TVOSMaxVersion   types.String `tfsdk:"tvos_max_version"`
	TVOSMinVersion   types.String `tfsdk:"tvos_min_version"`
	DefaultShard     types.Int64  `tfsdk:"default_shard"`
	ShardModulo      types.Int64  `tfsdk:"shard_modulo"`
	ExcludedTagIDs   types.Set    `tfsdk:"excluded_tag_ids"`
	TagShards        types.Set    `tfsdk:"tag_shards"`
}

func mdmBlueprintArtifactForState(mba *goztl.MDMBlueprintArtifact) mdmBlueprintArtifact {
	exTagIDs := make([]attr.Value, 0)
	for _, exTagID := range mba.ExcludedTagIDs {
		exTagIDs = append(exTagIDs, types.Int64Value(int64(exTagID)))
	}

	tagShards := make([]attr.Value, 0)
	for _, tagShard := range mba.TagShards {
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

	return mdmBlueprintArtifact{
		ID:               types.Int64Value(int64(mba.ID)),
		BlueprintID:      types.Int64Value(int64(mba.BlueprintID)),
		ArtifactID:       types.StringValue(mba.ArtifactID),
		IOS:              types.BoolValue(mba.IOS),
		IOSMaxVersion:    types.StringValue(mba.IOSMaxVersion),
		IOSMinVersion:    types.StringValue(mba.IOSMinVersion),
		IPadOS:           types.BoolValue(mba.IPadOS),
		IPadOSMaxVersion: types.StringValue(mba.IPadOSMaxVersion),
		IPadOSMinVersion: types.StringValue(mba.IPadOSMinVersion),
		MacOS:            types.BoolValue(mba.MacOS),
		MacOSMaxVersion:  types.StringValue(mba.MacOSMaxVersion),
		MacOSMinVersion:  types.StringValue(mba.MacOSMinVersion),
		TVOS:             types.BoolValue(mba.TVOS),
		TVOSMaxVersion:   types.StringValue(mba.TVOSMaxVersion),
		TVOSMinVersion:   types.StringValue(mba.TVOSMinVersion),
		DefaultShard:     types.Int64Value(int64(mba.DefaultShard)),
		ShardModulo:      types.Int64Value(int64(mba.ShardModulo)),
		ExcludedTagIDs:   types.SetValueMust(types.Int64Type, exTagIDs),
		TagShards:        types.SetValueMust(types.ObjectType{AttrTypes: tagShardAttrTypes}, tagShards),
	}
}

func mdmBlueprintArtifactRequestWithState(data mdmBlueprintArtifact) *goztl.MDMBlueprintArtifactRequest {
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

	return &goztl.MDMBlueprintArtifactRequest{
		BlueprintID:      int(data.BlueprintID.ValueInt64()),
		ArtifactID:       data.ArtifactID.ValueString(),
		IOS:              data.IOS.ValueBool(),
		IOSMaxVersion:    data.IOSMaxVersion.ValueString(),
		IOSMinVersion:    data.IOSMinVersion.ValueString(),
		IPadOS:           data.IPadOS.ValueBool(),
		IPadOSMaxVersion: data.IPadOSMaxVersion.ValueString(),
		IPadOSMinVersion: data.IPadOSMinVersion.ValueString(),
		MacOS:            data.MacOS.ValueBool(),
		MacOSMaxVersion:  data.MacOSMaxVersion.ValueString(),
		MacOSMinVersion:  data.MacOSMinVersion.ValueString(),
		TVOS:             data.TVOS.ValueBool(),
		TVOSMaxVersion:   data.TVOSMaxVersion.ValueString(),
		TVOSMinVersion:   data.TVOSMinVersion.ValueString(),
		DefaultShard:     int(data.DefaultShard.ValueInt64()),
		ShardModulo:      int(data.ShardModulo.ValueInt64()),
		ExcludedTagIDs:   exTagIDs,
		TagShards:        tagShards,
	}
}
