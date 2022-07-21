package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type metaBusinessUnit struct {
	ID                   types.Int64  `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	APIEnrollmentEnabled types.Bool   `tfsdk:"api_enrollment_enabled"`
}

func metaBusinessUnitForState(mbu *goztl.MetaBusinessUnit) metaBusinessUnit {
	return metaBusinessUnit{
		ID:                   types.Int64{Value: int64(mbu.ID)},
		Name:                 types.String{Value: mbu.Name},
		APIEnrollmentEnabled: types.Bool{Value: mbu.APIEnrollmentEnabled},
	}
}
