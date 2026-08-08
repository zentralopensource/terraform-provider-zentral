package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type osqueryATC struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	TableName   types.String `tfsdk:"table_name"`
	Query       types.String `tfsdk:"query"`
	Path        types.String `tfsdk:"path"`
	Columns     types.List   `tfsdk:"columns"`
	Platforms   types.Set    `tfsdk:"platforms"`
}

func osqueryATCForState(oa *goztl.OsqueryATC) osqueryATC {
	return osqueryATC{
		ID:          types.Int64Value(int64(oa.ID)),
		Name:        types.StringValue(oa.Name),
		Description: types.StringValue(oa.Description),
		TableName:   types.StringValue(oa.TableName),
		Query:       types.StringValue(oa.Query),
		Path:        types.StringValue(oa.Path),
		Columns:     stringListForState(oa.Columns),
		Platforms:   stringSetForState(oa.Platforms),
	}
}

func osqueryATCRequestWithState(data osqueryATC) *goztl.OsqueryATCRequest {
	return &goztl.OsqueryATCRequest{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		TableName:   data.TableName.ValueString(),
		Query:       data.Query.ValueString(),
		Path:        data.Path.ValueString(),
		Columns:     stringListWithStateList(data.Columns),
		Platforms:   stringListWithStateSet(data.Platforms),
	}
}
