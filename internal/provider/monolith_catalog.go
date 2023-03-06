package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type monolithCatalog struct {
	ID       types.Int64  `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Priority types.Int64  `tfsdk:"priority"`
}

func monolithCatalogForState(mc *goztl.MonolithCatalog) monolithCatalog {
	return monolithCatalog{
		ID:       types.Int64Value(int64(mc.ID)),
		Name:     types.StringValue(mc.Name),
		Priority: types.Int64Value(int64(mc.Priority)),
	}
}

func monolithCatalogRequestWithState(data monolithCatalog) *goztl.MonolithCatalogRequest {
	return &goztl.MonolithCatalogRequest{
		Name:     data.Name.ValueString(),
		Priority: int(data.Priority.ValueInt64()),
	}
}
