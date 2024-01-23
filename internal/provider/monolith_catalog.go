package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type monolithCatalog struct {
	ID           types.Int64  `tfsdk:"id"`
	RepositoryID types.Int64  `tfsdk:"repository_id"`
	Name         types.String `tfsdk:"name"`
}

func monolithCatalogForState(mc *goztl.MonolithCatalog) monolithCatalog {
	return monolithCatalog{
		ID:           types.Int64Value(int64(mc.ID)),
		RepositoryID: types.Int64Value(int64(mc.RepositoryID)),
		Name:         types.StringValue(mc.Name),
	}
}

func monolithCatalogRequestWithState(data monolithCatalog) *goztl.MonolithCatalogRequest {
	return &goztl.MonolithCatalogRequest{
		RepositoryID: int(data.RepositoryID.ValueInt64()),
		Name:         data.Name.ValueString(),
	}
}
