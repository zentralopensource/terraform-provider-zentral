package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmProfile struct {
	ID               types.String `tfsdk:"id"`
	Source           types.String `tfsdk:"source"`
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

func mdmProfileForState(mp *goztl.MDMProfile) mdmProfile {
	return mdmProfile{
		ID:               types.StringValue(mp.ID),
		Source:           types.StringValue(mp.Source),
		ArtifactID:       types.StringValue(mp.ArtifactID),
		IOS:              types.BoolValue(mp.IOS),
		IOSMaxVersion:    types.StringValue(mp.IOSMaxVersion),
		IOSMinVersion:    types.StringValue(mp.IOSMinVersion),
		IPadOS:           types.BoolValue(mp.IPadOS),
		IPadOSMaxVersion: types.StringValue(mp.IPadOSMaxVersion),
		IPadOSMinVersion: types.StringValue(mp.IPadOSMinVersion),
		MacOS:            types.BoolValue(mp.MacOS),
		MacOSMaxVersion:  types.StringValue(mp.MacOSMaxVersion),
		MacOSMinVersion:  types.StringValue(mp.MacOSMinVersion),
		TVOS:             types.BoolValue(mp.TVOS),
		TVOSMaxVersion:   types.StringValue(mp.TVOSMaxVersion),
		TVOSMinVersion:   types.StringValue(mp.TVOSMinVersion),
		DefaultShard:     types.Int64Value(int64(mp.DefaultShard)),
		ShardModulo:      types.Int64Value(int64(mp.ShardModulo)),
		ExcludedTagIDs:   int64SetForState(mp.MDMArtifactVersion.ExcludedTagIDs),
		TagShards:        tagShardsForState(mp.MDMArtifactVersion),
		Version:          types.Int64Value(int64(mp.Version)),
	}
}

func mdmProfileRequestWithState(data mdmProfile) *goztl.MDMProfileRequest {
	return &goztl.MDMProfileRequest{
		Source: data.Source.ValueString(),
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
			ExcludedTagIDs:   intListWithState(data.ExcludedTagIDs),
			TagShards:        tagShardsWithState(data.TagShards),
			Version:          int(data.Version.ValueInt64()),
		},
	}
}
