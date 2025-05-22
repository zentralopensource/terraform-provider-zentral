package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmDeclaration struct {
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

// source serialization / deserialization
// we try to minimize the risk of Terraform thinking that the state has changed
// because the source could be serialized differently

func serializeMDMDeclarationSource(mds goztl.MDMDeclarationSource) (string, error) {
	b, err := json.Marshal(mds)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func deserializeMDMDeclarationSource(src string) (goztl.MDMDeclarationSource, error) {
	var mds goztl.MDMDeclarationSource
	err := json.Unmarshal([]byte(src), &mds)
	if err != nil {
		return mds, err
	}
	return mds, nil
}

// conversion between TF state and goztl

func mdmDeclarationForState(mda *goztl.MDMDeclaration) (mdmDeclaration, error) {
	source, err := serializeMDMDeclarationSource(mda.Source)
	if err != nil {
		return mdmDeclaration{}, err
	}

	exTagIDs := exTagIDsForState(mda.MDMArtifactVersion)
	tagShards := tagShardsForState(mda.MDMArtifactVersion)

	return mdmDeclaration{
		ID:               types.StringValue(mda.ID),
		Source:           types.StringValue(source),
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
	}, nil
}

func mdmDeclarationRequestWithState(data mdmDeclaration) (*goztl.MDMDeclarationRequest, error) {
	source, err := deserializeMDMDeclarationSource(data.Source.ValueString())
	if err != nil {
		return nil, err
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

	return &goztl.MDMDeclarationRequest{
		Source: source,
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
	}, nil
}

// MDM declaration source validator

var _ validator.String = mdmDeclarationSourceValidator{}

type mdmDeclarationSourceValidator struct{}

func (v mdmDeclarationSourceValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v mdmDeclarationSourceValidator) MarkdownDescription(_ context.Context) string {
	return "value must be a valid DDM declaration"
}

func (v mdmDeclarationSourceValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	_, err := deserializeMDMDeclarationSource(value.ValueString())
	if err == nil {
		return
	}

	response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
		request.Path,
		v.Description(ctx),
		value.String(),
	))
}
