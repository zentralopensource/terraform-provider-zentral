package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type gwsConnection struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func gwsConnectionForState(gwsConn *goztl.GWSConnection) gwsConnection {
	return gwsConnection{
		ID:   types.StringValue((gwsConn.ID)),
		Name: types.StringValue(gwsConn.Name),
	}
}
