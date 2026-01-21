package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmEnrollmentCustomView struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Description            types.String `tfsdk:"description"`
	HTML                   types.String `tfsdk:"html"`
	RequiresAuthentication types.Bool   `tfsdk:"requires_authentication"`
}

func mdmEnrollmentCustomViewForState(customView *goztl.MDMEnrollmentCustomView) mdmEnrollmentCustomView {
	return mdmEnrollmentCustomView{
		ID:                     types.StringValue(customView.ID),
		Name:                   types.StringValue(customView.Name),
		Description:            types.StringValue(customView.Description),
		HTML:                   types.StringValue(customView.HTML),
		RequiresAuthentication: types.BoolValue(customView.RequiresAuthentication),
	}
}

func mdmEnrollmentCustomViewRequestWithState(data mdmEnrollmentCustomView) *goztl.MDMEnrollmentCustomViewRequest {
	return &goztl.MDMEnrollmentCustomViewRequest{
		Name:                   data.Name.ValueString(),
		Description:            data.Description.ValueString(),
		HTML:                   data.HTML.ValueString(),
		RequiresAuthentication: data.RequiresAuthentication.ValueBool(),
	}
}
