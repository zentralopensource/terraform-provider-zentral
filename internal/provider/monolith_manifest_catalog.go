package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type monolithManifestCatalog struct {
	ID         types.Int64 `tfsdk:"id"`
	ManifestID types.Int64 `tfsdk:"manifest_id"`
	CatalogID  types.Int64 `tfsdk:"catalog_id"`
	TagIDs     types.Set   `tfsdk:"tag_ids"`
}

func monolithManifestCatalogForState(mmc *goztl.MonolithManifestCatalog) monolithManifestCatalog {
	tagIDs := make([]attr.Value, 0)
	for _, tagID := range mmc.TagIDs {
		tagIDs = append(tagIDs, types.Int64Value(int64(tagID)))
	}

	return monolithManifestCatalog{
		ID:         types.Int64Value(int64(mmc.ID)),
		ManifestID: types.Int64Value(int64(mmc.ManifestID)),
		CatalogID:  types.Int64Value(int64(mmc.CatalogID)),
		TagIDs:     types.SetValueMust(types.Int64Type, tagIDs),
	}
}

func monolithManifestCatalogRequestWithState(data monolithManifestCatalog) *goztl.MonolithManifestCatalogRequest {
	tagIDs := make([]int, 0)
	for _, tagID := range data.TagIDs.Elements() { // nil if null or unknown → no iterations
		tagIDs = append(tagIDs, int(tagID.(types.Int64).ValueInt64()))
	}

	return &goztl.MonolithManifestCatalogRequest{
		ManifestID: int(data.ManifestID.ValueInt64()),
		CatalogID:  int(data.CatalogID.ValueInt64()),
		TagIDs:     tagIDs,
	}
}
