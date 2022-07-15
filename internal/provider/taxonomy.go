package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type taxonomy struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func taxonomyForState(t *goztl.Taxonomy) taxonomy {
	return taxonomy{
		ID:   types.Int64{Value: int64(t.ID)},
		Name: types.String{Value: t.Name},
	}
}
