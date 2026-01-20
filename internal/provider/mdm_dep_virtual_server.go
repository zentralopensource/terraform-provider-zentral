package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmDEPVirtualServer struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func mdmDEPVirtualServerForState(virtualServer *goztl.MDMDEPVirtualServer) mdmDEPVirtualServer {
	return mdmDEPVirtualServer{
		ID:   types.Int64Value(int64(virtualServer.ID)),
		Name: types.StringValue(virtualServer.Name),
	}
}
