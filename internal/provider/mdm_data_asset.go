package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmDataAsset struct {
	ID               types.String `tfsdk:"id"`
	Type             types.String `tfsdk:"type"`
	FileURI          types.String `tfsdk:"file_uri"`
	FileSHA256       types.String `tfsdk:"file_sha256"`
	FileSize         types.Int64  `tfsdk:"file_size"`
	Filename         types.String `tfsdk:"filename"`
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
	Version          types.Int64  `tfsdk:"version"`
}

func mdmDataAssetForState(mda *goztl.MDMDataAsset, fileURI types.String) mdmDataAsset {
	exTagIDs := exTagIDsForState(mda.MDMArtifactVersion)
	tagShards := tagShardsForState(mda.MDMArtifactVersion)

	return mdmDataAsset{
		ID:               types.StringValue(mda.ID),
		Type:             types.StringValue(mda.Type),
		FileURI:          fileURI,
		FileSHA256:       types.StringValue(mda.FileSHA256),
		FileSize:         types.Int64Value(mda.FileSize),
		Filename:         types.StringValue(mda.Filename),
		ArtifactID:       types.StringValue(mda.ArtifactID),
		IOS:              types.BoolValue(mda.IOS),
		IOSMaxVersion:    types.StringValue(mda.IOSMaxVersion),
		IOSMinVersion:    types.StringValue(mda.IOSMinVersion),
		IPadOS:           types.BoolValue(mda.IPadOS),
		IPadOSMaxVersion: types.StringValue(mda.IPadOSMaxVersion),
		IPadOSMinVersion: types.StringValue(mda.IPadOSMinVersion),
		MacOS:            types.BoolValue(mda.MacOS),
		MacOSMaxVersion:  types.StringValue(mda.MacOSMaxVersion),
		MacOSMinVersion:  types.StringValue(mda.MacOSMinVersion),
		TVOS:             types.BoolValue(mda.TVOS),
		TVOSMaxVersion:   types.StringValue(mda.TVOSMaxVersion),
		TVOSMinVersion:   types.StringValue(mda.TVOSMinVersion),
		DefaultShard:     types.Int64Value(int64(mda.DefaultShard)),
		ShardModulo:      types.Int64Value(int64(mda.ShardModulo)),
		ExcludedTagIDs:   types.SetValueMust(types.Int64Type, exTagIDs),
		TagShards:        types.SetValueMust(types.ObjectType{AttrTypes: tagShardAttrTypes}, tagShards),
		Version:          types.Int64Value(int64(mda.Version)),
	}
}

func mdmDataAssetRequestWithState(data mdmDataAsset) *goztl.MDMDataAssetRequest {
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

	return &goztl.MDMDataAssetRequest{
		Type:       data.Type.ValueString(),
		FileURI:    data.FileURI.ValueString(),
		FileSHA256: data.FileSHA256.ValueString(),
		MDMArtifactVersionRequest: goztl.MDMArtifactVersionRequest{
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
			Version:          int(data.Version.ValueInt64()),
		},
	}
}
