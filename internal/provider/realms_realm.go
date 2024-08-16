package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type realmsRealm struct {
	ID   types.String `tfsdk:"id"`
	UUID types.String `tfsdk:"uuid"`
	Name types.String `tfsdk:"name"`
}

func realmsRealmForState(r *goztl.RealmsRealm) realmsRealm {
	return realmsRealm{
		ID:   types.StringValue(r.UUID),
		UUID: types.StringValue(r.UUID),
		Name: types.StringValue(r.Name),
	}
}
