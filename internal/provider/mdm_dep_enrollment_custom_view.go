package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmDEPEnrollmentCutomView struct {
	ID              types.String `tfsdk:"id"`
	DEPEnrollmentID types.Int64  `tfsdk:"dep_enrollment"`
	CustomViewID    types.String `tfsdk:"custom_view"`
	Weight          types.Int64  `tfsdk:"weight"`
}

func mdmDEPEnrollmentCustomViewForState(customView *goztl.MDMDEPEnrollmentCustomView) mdmDEPEnrollmentCutomView {
	return mdmDEPEnrollmentCutomView{
		ID:              types.StringValue(customView.ID),
		DEPEnrollmentID: types.Int64Value(int64(customView.DEPEnrollmentID)),
		CustomViewID:    types.StringValue(customView.CustomViewID),
		Weight:          types.Int64Value(int64(customView.Weight)),
	}
}

func mdmDEPEnrollmentCustomViewRequestWithState(data mdmDEPEnrollmentCutomView) *goztl.MDMDEPEnrollmentCustomViewRequest {
	mdmEnrollmentCustomViewRequest := &goztl.MDMDEPEnrollmentCustomViewRequest{
		DEPEnrollmentID: int(data.DEPEnrollmentID.ValueInt64()),
		CustomViewID:    data.CustomViewID.ValueString(),
		Weight:          int(data.Weight.ValueInt64()),
	}
	return mdmEnrollmentCustomViewRequest
}
