package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	columns := make([]attr.Value, 0)
	for _, column := range oa.Columns {
		columns = append(columns, types.StringValue(column))
	}

	platforms := make([]attr.Value, 0)
	for _, platform := range oa.Platforms {
		platforms = append(platforms, types.StringValue(platform))
	}

	return osqueryATC{
		ID:          types.Int64Value(int64(oa.ID)),
		Name:        types.StringValue(oa.Name),
		Description: types.StringValue(oa.Description),
		TableName:   types.StringValue(oa.TableName),
		Query:       types.StringValue(oa.Query),
		Path:        types.StringValue(oa.Path),
		Columns:     types.ListValueMust(types.StringType, columns),
		Platforms:   types.SetValueMust(types.StringType, platforms),
	}
}

func osqueryATCRequestWithState(data osqueryATC) *goztl.OsqueryATCRequest {
	columns := make([]string, 0)
	for _, column := range data.Columns.Elements() { // nil if null or unknown → no iterations
		columns = append(columns, column.(types.String).ValueString())
	}

	platforms := make([]string, 0)
	for _, platform := range data.Platforms.Elements() { // nil if null or unknown → no iterations
		platforms = append(platforms, platform.(types.String).ValueString())
	}

	return &goztl.OsqueryATCRequest{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		TableName:   data.TableName.ValueString(),
		Query:       data.Query.ValueString(),
		Path:        data.Path.ValueString(),
		Columns:     columns,
		Platforms:   platforms,
	}
}
