package provider

import (
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
	return monolithManifestSubManifest{
		ID:            types.Int64Value(int64(mmsm.ID)),
		ManifestID:    types.Int64Value(int64(mmsm.ManifestID)),
		SubManifestID: types.Int64Value(int64(mmsm.SubManifestID)),
		TagIDs:        int64SetForState(mmsm.TagIDs),
	}
}

func monolithManifestSubManifestRequestWithState(data monolithManifestSubManifest) *goztl.MonolithManifestSubManifestRequest {
	return &goztl.MonolithManifestSubManifestRequest{
		ManifestID:    int(data.ManifestID.ValueInt64()),
		SubManifestID: int(data.SubManifestID.ValueInt64()),
		TagIDs:        intListWithState(data.TagIDs),
	}
}
