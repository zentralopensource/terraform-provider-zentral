package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type metaBusinessUnit struct {
	Name types.String `tfsdk:"name"`
	ID   types.Int64  `tfsdk:"id"`
}

func metaBusinessUnitForState(mbu *goztl.MetaBusinessUnit) metaBusinessUnit {
	return metaBusinessUnit{
		Name: types.String{Value: mbu.Name},
		ID:   types.Int64{Value: int64(mbu.ID)},
	}
}
