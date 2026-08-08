package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type tag struct {
	ID         types.Int64  `tfsdk:"id"`
	TaxonomyID types.Int64  `tfsdk:"taxonomy_id"`
	Name       types.String `tfsdk:"name"`
	Color      types.String `tfsdk:"color"`
}

func tagForState(t *goztl.Tag) tag {
	return tag{
		ID:         types.Int64Value(int64(t.ID)),
		TaxonomyID: optionalInt64ForState(t.TaxonomyID),
		Name:       types.StringValue(t.Name),
		Color:      types.StringValue(t.Color),
	}
}
