package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmLocation struct {
	ID               types.Int64  `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	OrganizationName types.String `tfsdk:"organization_name"`
	MDMInfoID        types.String `tfsdk:"mdm_info_id"`
}

func mdmLocationForState(ml *goztl.MDMLocation) mdmLocation {
	return mdmLocation{
		ID:               types.Int64Value(int64(ml.ID)),
		Name:             types.StringValue(ml.Name),
		OrganizationName: types.StringValue(ml.OrganizationName),
		MDMInfoID:        types.StringValue(ml.MDMInfoID),
	}
}
