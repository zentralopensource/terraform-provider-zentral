package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmProvisioningProfile struct {
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

func mdmProvisioningProfileForState(mpp *goztl.MDMProvisioningProfile) mdmProvisioningProfile {
	return mdmProvisioningProfile{
		ID:               types.StringValue(mpp.ID),
		Source:           types.StringValue(mpp.Source),
		ArtifactID:       types.StringValue(mpp.ArtifactID),
		IOS:              types.BoolValue(mpp.IOS),
		IOSMaxVersion:    types.StringValue(mpp.IOSMaxVersion),
		IOSMinVersion:    types.StringValue(mpp.IOSMinVersion),
		IPadOS:           types.BoolValue(mpp.IPadOS),
		IPadOSMaxVersion: types.StringValue(mpp.IPadOSMaxVersion),
		IPadOSMinVersion: types.StringValue(mpp.IPadOSMinVersion),
		MacOS:            types.BoolValue(mpp.MacOS),
		MacOSMaxVersion:  types.StringValue(mpp.MacOSMaxVersion),
		MacOSMinVersion:  types.StringValue(mpp.MacOSMinVersion),
		TVOS:             types.BoolValue(mpp.TVOS),
		TVOSMaxVersion:   types.StringValue(mpp.TVOSMaxVersion),
		TVOSMinVersion:   types.StringValue(mpp.TVOSMinVersion),
		DefaultShard:     types.Int64Value(int64(mpp.DefaultShard)),
		ShardModulo:      types.Int64Value(int64(mpp.ShardModulo)),
		ExcludedTagIDs:   int64SetForState(mpp.MDMArtifactVersion.ExcludedTagIDs),
		TagShards:        tagShardsForState(mpp.MDMArtifactVersion),
		Version:          types.Int64Value(int64(mpp.Version)),
	}
}

func mdmProvisioningProfileRequestWithState(data mdmProvisioningProfile) *goztl.MDMProvisioningProfileRequest {
	return &goztl.MDMProvisioningProfileRequest{
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
