package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmEnterpriseApp struct {
	ID               types.String `tfsdk:"id"`
	PackageURI       types.String `tfsdk:"package_uri"`
	PackageSHA256    types.String `tfsdk:"package_sha256"`
	IOSApp           types.Bool   `tfsdk:"ios_app"`
	Configuration    types.String `tfsdk:"configuration"`
	InstallAsManaged types.Bool   `tfsdk:"install_as_managed"`
	RemoveOnUnenroll types.Bool   `tfsdk:"remove_on_unenroll"`
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

func mdmEnterpriseAppForState(mea *goztl.MDMEnterpriseApp) mdmEnterpriseApp {
	exTagIDs := exTagIDsForState(mea.MDMArtifactVersion)
	tagShards := tagShardsForState(mea.MDMArtifactVersion)

	var configuration types.String
	if mea.Configuration != nil {
		configuration = types.StringValue(*mea.Configuration)
	} else {
		configuration = types.StringNull()
	}

	return mdmEnterpriseApp{
		ID:               types.StringValue(mea.ID),
		PackageURI:       types.StringValue(mea.PackageURI),
		PackageSHA256:    types.StringValue(mea.PackageSHA256),
		IOSApp:           types.BoolValue(mea.IOSApp),
		Configuration:    configuration,
		InstallAsManaged: types.BoolValue(mea.InstallAsManaged),
		RemoveOnUnenroll: types.BoolValue(mea.RemoveOnUnenroll),
		ArtifactID:       types.StringValue(mea.ArtifactID),
		IOS:              types.BoolValue(mea.IOS),
		IOSMaxVersion:    types.StringValue(mea.IOSMaxVersion),
		IOSMinVersion:    types.StringValue(mea.IOSMinVersion),
		IPadOS:           types.BoolValue(mea.IPadOS),
		IPadOSMaxVersion: types.StringValue(mea.IPadOSMaxVersion),
		IPadOSMinVersion: types.StringValue(mea.IPadOSMinVersion),
		MacOS:            types.BoolValue(mea.MacOS),
		MacOSMaxVersion:  types.StringValue(mea.MacOSMaxVersion),
		MacOSMinVersion:  types.StringValue(mea.MacOSMinVersion),
		TVOS:             types.BoolValue(mea.TVOS),
		TVOSMaxVersion:   types.StringValue(mea.TVOSMaxVersion),
		TVOSMinVersion:   types.StringValue(mea.TVOSMinVersion),
		DefaultShard:     types.Int64Value(int64(mea.DefaultShard)),
		ShardModulo:      types.Int64Value(int64(mea.ShardModulo)),
		ExcludedTagIDs:   types.SetValueMust(types.Int64Type, exTagIDs),
		TagShards:        types.SetValueMust(types.ObjectType{AttrTypes: tagShardAttrTypes}, tagShards),
		Version:          types.Int64Value(int64(mea.Version)),
	}
}

func mdmEnterpriseAppRequestWithState(data mdmEnterpriseApp) *goztl.MDMEnterpriseAppRequest {
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

	var configuration *string
	if !data.Configuration.IsNull() {
		configuration = goztl.String(data.Configuration.ValueString())
	}

	return &goztl.MDMEnterpriseAppRequest{
		PackageURI:       data.PackageURI.ValueString(),
		PackageSHA256:    data.PackageSHA256.ValueString(),
		IOSApp:           data.IOSApp.ValueBool(),
		Configuration:    configuration,
		InstallAsManaged: data.InstallAsManaged.ValueBool(),
		RemoveOnUnenroll: data.RemoveOnUnenroll.ValueBool(),
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
