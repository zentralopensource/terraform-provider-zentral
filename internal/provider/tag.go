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
	var taxonomyID types.Int64
	if t.TaxonomyID != nil {
		taxonomyID = types.Int64Value(int64(*t.TaxonomyID))
	} else {
		taxonomyID = types.Int64Null()
	}
	return tag{
		ID:         types.Int64Value(int64(t.ID)),
		TaxonomyID: taxonomyID,
		Name:       types.StringValue(t.Name),
		Color:      types.StringValue(t.Color),
	}
}
