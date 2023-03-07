package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type monolithManifestSubManifest struct {
	ID            types.Int64 `tfsdk:"id"`
	ManifestID    types.Int64 `tfsdk:"manifest_id"`
	SubManifestID types.Int64 `tfsdk:"sub_manifest_id"`
	TagIDs        types.Set   `tfsdk:"tag_ids"`
}

func monolithManifestSubManifestForState(mmsm *goztl.MonolithManifestSubManifest) monolithManifestSubManifest {
	tagIDs := make([]attr.Value, 0)
	for _, tagID := range mmsm.TagIDs {
		tagIDs = append(tagIDs, types.Int64Value(int64(tagID)))
	}

	return monolithManifestSubManifest{
		ID:            types.Int64Value(int64(mmsm.ID)),
		ManifestID:    types.Int64Value(int64(mmsm.ManifestID)),
		SubManifestID: types.Int64Value(int64(mmsm.SubManifestID)),
		TagIDs:        types.SetValueMust(types.Int64Type, tagIDs),
	}
}

func monolithManifestSubManifestRequestWithState(data monolithManifestSubManifest) *goztl.MonolithManifestSubManifestRequest {
	tagIDs := make([]int, 0)
	for _, tagID := range data.TagIDs.Elements() { // nil if null or unknown → no iterations
		tagIDs = append(tagIDs, int(tagID.(types.Int64).ValueInt64()))
	}

	return &goztl.MonolithManifestSubManifestRequest{
		ManifestID:    int(data.ManifestID.ValueInt64()),
		SubManifestID: int(data.SubManifestID.ValueInt64()),
		TagIDs:        tagIDs,
	}
}
