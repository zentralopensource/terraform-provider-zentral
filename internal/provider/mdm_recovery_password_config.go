package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmRecoveryPasswordConfig struct {
	ID                     types.Int64  `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	DynamicPassword        types.Bool   `tfsdk:"dynamic_password"`
	StaticPassword         types.String `tfsdk:"static_password"`
	RotationIntervalDays   types.Int64  `tfsdk:"rotation_interval_days"`
	RotateFirmwarePassword types.Bool   `tfsdk:"rotate_firmware_password"`
}

func mdmRecoveryPasswordConfigForState(mrpc *goztl.MDMRecoveryPasswordConfig) mdmRecoveryPasswordConfig {
	var staticPassword types.String
	if mrpc.StaticPassword != nil {
		staticPassword = types.StringValue(*mrpc.StaticPassword)
	} else {
		staticPassword = types.StringNull()
	}

	return mdmRecoveryPasswordConfig{
		ID:                     types.Int64Value(int64(mrpc.ID)),
		Name:                   types.StringValue(mrpc.Name),
		DynamicPassword:        types.BoolValue(mrpc.DynamicPassword),
		StaticPassword:         staticPassword,
		RotationIntervalDays:   types.Int64Value(int64(mrpc.RotationIntervalDays)),
		RotateFirmwarePassword: types.BoolValue(mrpc.RotateFirmwarePassword),
	}
}

func mdmRecoveryPasswordConfigRequestWithState(data mdmRecoveryPasswordConfig) *goztl.MDMRecoveryPasswordConfigRequest {
	var staticPassword *string
	if !data.StaticPassword.IsNull() {
		staticPassword = goztl.String(data.StaticPassword.ValueString())
	}

	return &goztl.MDMRecoveryPasswordConfigRequest{
		Name:                   data.Name.ValueString(),
		DynamicPassword:        data.DynamicPassword.ValueBool(),
		StaticPassword:         staticPassword,
		RotationIntervalDays:   int(data.RotationIntervalDays.ValueInt64()),
		RotateFirmwarePassword: data.RotateFirmwarePassword.ValueBool(),
	}
}
