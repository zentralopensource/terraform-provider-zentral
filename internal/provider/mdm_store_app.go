package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmStoreApp struct {
	ID                                     types.String `tfsdk:"id"`
	LocationAssetID                        types.Int64  `tfsdk:"location_asset_id"`
	AssociatedDomains                      types.List   `tfsdk:"associated_domains"`
	AssociatedDomainsEnableDirectDownloads types.Bool   `tfsdk:"associated_domains_enable_direct_downloads"`
	Configuration                          types.String `tfsdk:"configuration"`
	ContentFilterUUID                      types.String `tfsdk:"content_filter_uuid"`
	DNSProxyUUID                           types.String `tfsdk:"dns_proxy_uuid"`
	VPNUUID                                types.String `tfsdk:"vpn_uuid"`
	PreventBackup                          types.Bool   `tfsdk:"prevent_backup"`
	Removable                              types.Bool   `tfsdk:"removable"`
	RemoveOnUnenroll                       types.Bool   `tfsdk:"remove_on_unenroll"`
	ArtifactID                             types.String `tfsdk:"artifact_id"`
	IOS                                    types.Bool   `tfsdk:"ios"`
	IOSMaxVersion                          types.String `tfsdk:"ios_max_version"`
	IOSMinVersion                          types.String `tfsdk:"ios_min_version"`
	IPadOS                                 types.Bool   `tfsdk:"ipados"`
	IPadOSMaxVersion                       types.String `tfsdk:"ipados_max_version"`
	IPadOSMinVersion                       types.String `tfsdk:"ipados_min_version"`
	MacOS                                  types.Bool   `tfsdk:"macos"`
	MacOSMaxVersion                        types.String `tfsdk:"macos_max_version"`
	MacOSMinVersion                        types.String `tfsdk:"macos_min_version"`
	TVOS                                   types.Bool   `tfsdk:"tvos"`
	TVOSMaxVersion                         types.String `tfsdk:"tvos_max_version"`
	TVOSMinVersion                         types.String `tfsdk:"tvos_min_version"`
	DefaultShard                           types.Int64  `tfsdk:"default_shard"`
	ShardModulo                            types.Int64  `tfsdk:"shard_modulo"`
	ExcludedTagIDs                         types.Set    `tfsdk:"excluded_tag_ids"`
	TagShards                              types.Set    `tfsdk:"tag_shards"`
	Version                                types.Int64  `tfsdk:"version"`
}

func mdmStoreAppForState(msa *goztl.MDMStoreApp) mdmStoreApp {
	return mdmStoreApp{
		ID:                                     types.StringValue(msa.ID),
		LocationAssetID:                        types.Int64Value(int64(msa.LocationAssetID)),
		AssociatedDomains:                      stringListForState(msa.AssociatedDomains),
		AssociatedDomainsEnableDirectDownloads: types.BoolValue(msa.AssociatedDomainsEnableDirectDownloads),
		Configuration:                          optionalStringForState(msa.Configuration),
		ContentFilterUUID:                      optionalStringForState(msa.ContentFilterUUID),
		DNSProxyUUID:                           optionalStringForState(msa.DNSProxyUUID),
		VPNUUID:                                optionalStringForState(msa.VPNUUID),
		PreventBackup:                          types.BoolValue(msa.PreventBackup),
		Removable:                              types.BoolValue(msa.Removable),
		RemoveOnUnenroll:                       types.BoolValue(msa.RemoveOnUnenroll),
		ArtifactID:                             types.StringValue(msa.ArtifactID),
		IOS:                                    types.BoolValue(msa.IOS),
		IOSMaxVersion:                          types.StringValue(msa.IOSMaxVersion),
		IOSMinVersion:                          types.StringValue(msa.IOSMinVersion),
		IPadOS:                                 types.BoolValue(msa.IPadOS),
		IPadOSMaxVersion:                       types.StringValue(msa.IPadOSMaxVersion),
		IPadOSMinVersion:                       types.StringValue(msa.IPadOSMinVersion),
		MacOS:                                  types.BoolValue(msa.MacOS),
		MacOSMaxVersion:                        types.StringValue(msa.MacOSMaxVersion),
		MacOSMinVersion:                        types.StringValue(msa.MacOSMinVersion),
		TVOS:                                   types.BoolValue(msa.TVOS),
		TVOSMaxVersion:                         types.StringValue(msa.TVOSMaxVersion),
		TVOSMinVersion:                         types.StringValue(msa.TVOSMinVersion),
		DefaultShard:                           types.Int64Value(int64(msa.DefaultShard)),
		ShardModulo:                            types.Int64Value(int64(msa.ShardModulo)),
		ExcludedTagIDs:                         int64SetForState(msa.MDMArtifactVersion.ExcludedTagIDs),
		TagShards:                              tagShardsForState(msa.MDMArtifactVersion),
		Version:                                types.Int64Value(int64(msa.Version)),
	}
}

func mdmStoreAppRequestWithState(data mdmStoreApp) *goztl.MDMStoreAppRequest {
	return &goztl.MDMStoreAppRequest{
		LocationAssetID:                        int(data.LocationAssetID.ValueInt64()),
		AssociatedDomains:                      stringListWithStateList(data.AssociatedDomains),
		AssociatedDomainsEnableDirectDownloads: data.AssociatedDomainsEnableDirectDownloads.ValueBool(),
		Configuration:                          optionalStringWithState(data.Configuration),
		ContentFilterUUID:                      optionalStringWithState(data.ContentFilterUUID),
		DNSProxyUUID:                           optionalStringWithState(data.DNSProxyUUID),
		VPNUUID:                                optionalStringWithState(data.VPNUUID),
		PreventBackup:                          data.PreventBackup.ValueBool(),
		Removable:                              data.Removable.ValueBool(),
		RemoveOnUnenroll:                       data.RemoveOnUnenroll.ValueBool(),
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
