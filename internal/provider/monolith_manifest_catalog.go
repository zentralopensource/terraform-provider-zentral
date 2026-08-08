package provider

import (
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
	return monolithManifestCatalog{
		ID:         types.Int64Value(int64(mmc.ID)),
		ManifestID: types.Int64Value(int64(mmc.ManifestID)),
		CatalogID:  types.Int64Value(int64(mmc.CatalogID)),
		TagIDs:     int64SetForState(mmc.TagIDs),
	}
}

func monolithManifestCatalogRequestWithState(data monolithManifestCatalog) *goztl.MonolithManifestCatalogRequest {
	return &goztl.MonolithManifestCatalogRequest{
		ManifestID: int(data.ManifestID.ValueInt64()),
		CatalogID:  int(data.CatalogID.ValueInt64()),
		TagIDs:     intListWithState(data.TagIDs),
	}
}
