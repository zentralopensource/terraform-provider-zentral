package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type osqueryConfiguration struct {
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Inventory         types.Bool   `tfsdk:"inventory"`
	InventoryApps     types.Bool   `tfsdk:"inventory_apps"`
	InventoryEC2      types.Bool   `tfsdk:"inventory_ec2"`
	InventoryInterval types.Int64  `tfsdk:"inventory_interval"`
	Options           types.Map    `tfsdk:"options"`
}

func osqueryConfigurationForState(oc *goztl.OsqueryConfiguration) osqueryConfiguration {
	options := make(map[string]attr.Value)
	for k, v := range oc.Options {
		options[k] = types.StringValue(fmt.Sprintf("%v", v))
	}

	return osqueryConfiguration{
		ID:                types.Int64Value(int64(oc.ID)),
		Name:              types.StringValue(oc.Name),
		Description:       types.StringValue(oc.Description),
		Inventory:         types.BoolValue(oc.Inventory),
		InventoryApps:     types.BoolValue(oc.InventoryApps),
		InventoryEC2:      types.BoolValue(oc.InventoryEC2),
		InventoryInterval: types.Int64Value(int64(oc.InventoryInterval)),
		Options:           types.MapValueMust(types.StringType, options),
	}
}

func osqueryConfigurationRequestWithState(data osqueryConfiguration) *goztl.OsqueryConfigurationRequest {
	options := make(map[string]interface{})
	for k, v := range data.Options.Elements() {
		options[k] = v.(types.String).ValueString()
	}

	return &goztl.OsqueryConfigurationRequest{
		Name:              data.Name.ValueString(),
		Description:       data.Description.ValueString(),
		Inventory:         data.Inventory.ValueBool(),
		InventoryApps:     data.InventoryApps.ValueBool(),
		InventoryEC2:      data.InventoryEC2.ValueBool(),
		InventoryInterval: int(data.InventoryInterval.ValueInt64()),
		Options:           options,
	}
}
