package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const (
	tfCertAssetAccessibilityDefault          string = "Default"
	tfCertAssetAccessibilityAfterFirstUnlock        = "AfterFirstUnlock"
)

type mdmCertAsset struct {
	ID               types.String `tfsdk:"id"`
	Accessible       types.String `tfsdk:"accessible"`
	ACMEIssuerUUID   types.String `tfsdk:"acme_issuer_id"`
	SCEPIssuerUUID   types.String `tfsdk:"scep_issuer_id"`
	Subject          types.Set    `tfsdk:"subject"`
	SubjectAltName   types.Object `tfsdk:"subject_alt_name"`
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

var rdnAttrTypes = map[string]attr.Type{
	"type":  types.StringType,
	"value": types.StringType,
}

func subjectForState(mca *goztl.MDMCertAsset) types.Set {
	rdns := make([]attr.Value, 0)
	for _, rdn := range mca.Subject {
		rdns = append(
			rdns,
			types.ObjectValueMust(
				rdnAttrTypes,
				map[string]attr.Value{
					"type":  types.StringValue(rdn.Type),
					"value": types.StringValue(rdn.Value),
				},
			),
		)
	}
	return types.SetValueMust(types.ObjectType{AttrTypes: rdnAttrTypes}, rdns)
}

func subjectWithState(s types.Set) []goztl.MDMCertAssetRDN {
	subject := make([]goztl.MDMCertAssetRDN, 0)
	for _, rdn := range s.Elements() {
		rdnMap := rdn.(types.Object).Attributes()
		if rdnMap != nil {
			subject = append(
				subject,
				goztl.MDMCertAssetRDN{
					Type:  rdnMap["type"].(types.String).ValueString(),
					Value: rdnMap["value"].(types.String).ValueString(),
				},
			)
		}
	}
	return subject
}

var sanAttrTypes = map[string]attr.Type{
	"rfc822_name":       types.StringType,
	"uri":               types.StringType,
	"dns_name":          types.StringType,
	"nt_principal_name": types.StringType,
}

func sanForState(san goztl.MDMCertAssetSubjectAltName) types.Object {
	return types.ObjectValueMust(
		sanAttrTypes,
		map[string]attr.Value{
			"rfc822_name":       optionalStringForState(san.RFC822Name),
			"uri":               optionalStringForState(san.URI),
			"dns_name":          optionalStringForState(san.DNSName),
			"nt_principal_name": optionalStringForState(san.NTPrincipalName),
		},
	)
}

func defaultSAN() types.Object {
	return types.ObjectValueMust(
		sanAttrTypes,
		map[string]attr.Value{
			"rfc822_name":       types.StringNull(),
			"uri":               types.StringNull(),
			"dns_name":          types.StringNull(),
			"nt_principal_name": types.StringNull(),
		},
	)
}

func sanWithState(o types.Object) goztl.MDMCertAssetSubjectAltName {
	oMap := o.Attributes()
	return goztl.MDMCertAssetSubjectAltName{
		RFC822Name:      optionalStringWithState(oMap["rfc822_name"].(types.String)),
		URI:             optionalStringWithState(oMap["uri"].(types.String)),
		DNSName:         optionalStringWithState(oMap["dns_name"].(types.String)),
		NTPrincipalName: optionalStringWithState(oMap["nt_principal_name"].(types.String)),
	}
}

func mdmCertAssetForState(mca *goztl.MDMCertAsset) mdmCertAsset {
	return mdmCertAsset{
		ID:               types.StringValue(mca.ID),
		Accessible:       types.StringValue(mca.Accessible),
		ACMEIssuerUUID:   optionalStringForState(mca.ACMEIssuerUUID),
		SCEPIssuerUUID:   optionalStringForState(mca.SCEPIssuerUUID),
		Subject:          subjectForState(mca),
		SubjectAltName:   sanForState(mca.SubjectAltName),
		ArtifactID:       types.StringValue(mca.ArtifactID),
		IOS:              types.BoolValue(mca.IOS),
		IOSMaxVersion:    types.StringValue(mca.IOSMaxVersion),
		IOSMinVersion:    types.StringValue(mca.IOSMinVersion),
		IPadOS:           types.BoolValue(mca.IPadOS),
		IPadOSMaxVersion: types.StringValue(mca.IPadOSMaxVersion),
		IPadOSMinVersion: types.StringValue(mca.IPadOSMinVersion),
		MacOS:            types.BoolValue(mca.MacOS),
		MacOSMaxVersion:  types.StringValue(mca.MacOSMaxVersion),
		MacOSMinVersion:  types.StringValue(mca.MacOSMinVersion),
		TVOS:             types.BoolValue(mca.TVOS),
		TVOSMaxVersion:   types.StringValue(mca.TVOSMaxVersion),
		TVOSMinVersion:   types.StringValue(mca.TVOSMinVersion),
		DefaultShard:     types.Int64Value(int64(mca.DefaultShard)),
		ShardModulo:      types.Int64Value(int64(mca.ShardModulo)),
		ExcludedTagIDs:   int64SetForState(mca.MDMArtifactVersion.ExcludedTagIDs),
		TagShards:        tagShardsForState(mca.MDMArtifactVersion),
		Version:          types.Int64Value(int64(mca.Version)),
	}
}

func mdmCertAssetRequestWithState(data mdmCertAsset) *goztl.MDMCertAssetRequest {
	return &goztl.MDMCertAssetRequest{
		Accessible:     data.Accessible.ValueString(),
		ACMEIssuerUUID: optionalStringWithState(data.ACMEIssuerUUID),
		SCEPIssuerUUID: optionalStringWithState(data.SCEPIssuerUUID),
		Subject:        subjectWithState(data.Subject),
		SubjectAltName: sanWithState(data.SubjectAltName),
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
