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
	RevealRotationDelay    types.Int64  `tfsdk:"reveal_rotation_delay"`
	RotateFirmwarePassword types.Bool   `tfsdk:"rotate_firmware_password"`
}

func mdmRecoveryPasswordConfigForState(mrpc *goztl.MDMRecoveryPasswordConfig) mdmRecoveryPasswordConfig {
	return mdmRecoveryPasswordConfig{
		ID:                     types.Int64Value(int64(mrpc.ID)),
		Name:                   types.StringValue(mrpc.Name),
		DynamicPassword:        types.BoolValue(mrpc.DynamicPassword),
		StaticPassword:         optionalStringForState(mrpc.StaticPassword),
		RotationIntervalDays:   types.Int64Value(int64(mrpc.RotationIntervalDays)),
		RevealRotationDelay:    types.Int64Value(int64(mrpc.RevealRotationDelay)),
		RotateFirmwarePassword: types.BoolValue(mrpc.RotateFirmwarePassword),
	}
}

func mdmRecoveryPasswordConfigRequestWithState(data mdmRecoveryPasswordConfig) *goztl.MDMRecoveryPasswordConfigRequest {
	return &goztl.MDMRecoveryPasswordConfigRequest{
		Name:                   data.Name.ValueString(),
		DynamicPassword:        data.DynamicPassword.ValueBool(),
		StaticPassword:         optionalStringWithState(data.StaticPassword),
		RotationIntervalDays:   int(data.RotationIntervalDays.ValueInt64()),
		RevealRotationDelay:    int(data.RevealRotationDelay.ValueInt64()),
		RotateFirmwarePassword: data.RotateFirmwarePassword.ValueBool(),
	}
}
