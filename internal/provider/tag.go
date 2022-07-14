package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type tag struct {
	ID    types.Int64  `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Color types.String `tfsdk:"color"`
}

func tagForState(t *goztl.Tag) tag {
	return tag{
		ID:    types.Int64{Value: int64(t.ID)},
		Name:  types.String{Value: t.Name},
		Color: types.String{Value: t.Color},
	}
}
