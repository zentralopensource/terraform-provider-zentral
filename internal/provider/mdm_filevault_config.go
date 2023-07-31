package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmFileVaultConfig struct {
	ID                        types.Int64  `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	EscrowLocationDisplayName types.String `tfsdk:"escrow_location_display_name"`
	AtLoginOnly               types.Bool   `tfsdk:"at_login_only"`
	BypassAttempts            types.Int64  `tfsdk:"bypass_attempts"`
	ShowRecoveryKey           types.Bool   `tfsdk:"show_recovery_key"`
	DestroyKeyOnStandby       types.Bool   `tfsdk:"destroy_key_on_standby"`
	PRKRotationIntervalDays   types.Int64  `tfsdk:"prk_rotation_interval_days"`
}

func mdmFileVaultConfigForState(mfc *goztl.MDMFileVaultConfig) mdmFileVaultConfig {
	return mdmFileVaultConfig{
		ID:                        types.Int64Value(int64(mfc.ID)),
		Name:                      types.StringValue(mfc.Name),
		EscrowLocationDisplayName: types.StringValue(mfc.EscrowLocationDisplayName),
		AtLoginOnly:               types.BoolValue(mfc.AtLoginOnly),
		BypassAttempts:            types.Int64Value(int64(mfc.BypassAttempts)),
		ShowRecoveryKey:           types.BoolValue(mfc.ShowRecoveryKey),
		DestroyKeyOnStandby:       types.BoolValue(mfc.DestroyKeyOnStandby),
		PRKRotationIntervalDays:   types.Int64Value(int64(mfc.PRKRotationIntervalDays)),
	}
}

func mdmFileVaultConfigRequestWithState(data mdmFileVaultConfig) *goztl.MDMFileVaultConfigRequest {
	return &goztl.MDMFileVaultConfigRequest{
		Name:                      data.Name.ValueString(),
		EscrowLocationDisplayName: data.EscrowLocationDisplayName.ValueString(),
		AtLoginOnly:               data.AtLoginOnly.ValueBool(),
		BypassAttempts:            int(data.BypassAttempts.ValueInt64()),
		ShowRecoveryKey:           data.ShowRecoveryKey.ValueBool(),
		DestroyKeyOnStandby:       data.DestroyKeyOnStandby.ValueBool(),
		PRKRotationIntervalDays:   int(data.PRKRotationIntervalDays.ValueInt64()),
	}
}
