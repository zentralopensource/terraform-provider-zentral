package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type osqueryQuery struct {
	ID                     types.Int64  `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	SQL                    types.String `tfsdk:"sql"`
	Platforms              types.Set    `tfsdk:"platforms"`
	MinOsqueryVersion      types.String `tfsdk:"minimum_osquery_version"`
	Description            types.String `tfsdk:"description"`
	Value                  types.String `tfsdk:"value"`
	Version                types.Int64  `tfsdk:"version"`
	ComplianceCheckEnabled types.Bool   `tfsdk:"compliance_check_enabled"`
}

func osqueryQueryForState(oq *goztl.OsqueryQuery) osqueryQuery {
	platforms := make([]attr.Value, 0)
	for _, platform := range oq.Platforms {
		platforms = append(platforms, types.StringValue(platform))
	}

	var minOsqueryVersion types.String
	if oq.MinOsqueryVersion != nil {
		minOsqueryVersion = types.StringValue(*oq.MinOsqueryVersion)
	} else {
		minOsqueryVersion = types.StringNull()
	}

	return osqueryQuery{
		ID:                     types.Int64Value(int64(oq.ID)),
		Name:                   types.StringValue(oq.Name),
		SQL:                    types.StringValue(oq.SQL),
		Platforms:              types.SetValueMust(types.StringType, platforms),
		MinOsqueryVersion:      minOsqueryVersion,
		Description:            types.StringValue(oq.Description),
		Value:                  types.StringValue(oq.Value),
		Version:                types.Int64Value(int64(oq.Version)),
		ComplianceCheckEnabled: types.BoolValue(oq.ComplianceCheckEnabled),
	}
}

func osqueryQueryRequestWithState(data osqueryQuery) *goztl.OsqueryQueryRequest {
	platforms := make([]string, 0)
	for _, platform := range data.Platforms.Elements() { // nil if null or unknown → no iterations
		platforms = append(platforms, platform.(types.String).ValueString())
	}

	var minOsqueryVersion *string
	if !data.MinOsqueryVersion.IsNull() {
		minOsqueryVersion = goztl.String(data.MinOsqueryVersion.ValueString())
	}

	return &goztl.OsqueryQueryRequest{
		Name:                   data.Name.ValueString(),
		SQL:                    data.SQL.ValueString(),
		Platforms:              platforms,
		MinOsqueryVersion:      minOsqueryVersion,
		Description:            data.Description.ValueString(),
		Value:                  data.Value.ValueString(),
		ComplianceCheckEnabled: data.ComplianceCheckEnabled.ValueBool(),
	}
}
