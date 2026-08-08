package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type turboConfiguration struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	CollectInventory      types.Bool   `tfsdk:"collect_inventory"`
	InventoryInterval     types.Int64  `tfsdk:"inventory_interval"`
	DefaultCheckInterval  types.Int64  `tfsdk:"default_check_interval"`
	ConfigRefreshInterval types.Int64  `tfsdk:"config_refresh_interval"`
	ResultsBatchSize      types.Int64  `tfsdk:"results_batch_size"`
}

func turboConfigurationForState(tc *goztl.TurboConfiguration) turboConfiguration {
	return turboConfiguration{
		ID:                    types.StringValue(tc.ID),
		Name:                  types.StringValue(tc.Name),
		Description:           types.StringValue(tc.Description),
		CollectInventory:      types.BoolValue(tc.CollectInventory),
		InventoryInterval:     types.Int64Value(int64(tc.InventoryInterval)),
		DefaultCheckInterval:  types.Int64Value(int64(tc.DefaultCheckInterval)),
		ConfigRefreshInterval: types.Int64Value(int64(tc.ConfigRefreshInterval)),
		ResultsBatchSize:      types.Int64Value(int64(tc.ResultsBatchSize)),
	}
}

func turboConfigurationRequestWithState(data turboConfiguration) *goztl.TurboConfigurationRequest {
	return &goztl.TurboConfigurationRequest{
		Name:                  data.Name.ValueString(),
		Description:           data.Description.ValueString(),
		CollectInventory:      data.CollectInventory.ValueBool(),
		InventoryInterval:     int(data.InventoryInterval.ValueInt64()),
		DefaultCheckInterval:  int(data.DefaultCheckInterval.ValueInt64()),
		ConfigRefreshInterval: int(data.ConfigRefreshInterval.ValueInt64()),
		ResultsBatchSize:      int(data.ResultsBatchSize.ValueInt64()),
	}
}
