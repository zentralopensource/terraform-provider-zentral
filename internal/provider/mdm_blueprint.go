package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmBlueprint struct {
	ID                       types.Int64  `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	InventoryInterval        types.Int64  `tfsdk:"inventory_interval"`
	CollectApps              types.String `tfsdk:"collect_apps"`
	CollectCertificates      types.String `tfsdk:"collect_certificates"`
	CollectProfiles          types.String `tfsdk:"collect_profiles"`
	FileVaultConfigID        types.Int64  `tfsdk:"filevault_config_id"`
	RecoveryPasswordConfigID types.Int64  `tfsdk:"recovery_password_config_id"`
}

func collectionOptForState(collectionOpt int) types.String {
	switch collectionOpt {
	case 2:
		return types.StringValue("ALL")
	case 1:
		return types.StringValue("MANAGED_ONLY")
	default:
		return types.StringValue("NO")
	}
}

func collectionOptWithState(collectionOpt types.String) int {
	switch collectionOpt.ValueString() {
	case "ALL":
		return 2
	case "MANAGED_ONLY":
		return 1
	default:
		return 0
	}
}

func mdmBlueprintForState(mb *goztl.MDMBlueprint) mdmBlueprint {
	var fvCfgID types.Int64
	if mb.FileVaultConfigID != nil {
		fvCfgID = types.Int64Value(int64(*mb.FileVaultConfigID))
	} else {
		fvCfgID = types.Int64Null()
	}
	var rpCfgID types.Int64
	if mb.RecoveryPasswordConfigID != nil {
		rpCfgID = types.Int64Value(int64(*mb.RecoveryPasswordConfigID))
	} else {
		rpCfgID = types.Int64Null()
	}
	return mdmBlueprint{
		ID:                       types.Int64Value(int64(mb.ID)),
		Name:                     types.StringValue(mb.Name),
		InventoryInterval:        types.Int64Value(int64(mb.InventoryInterval)),
		CollectApps:              collectionOptForState(mb.CollectApps),
		CollectCertificates:      collectionOptForState(mb.CollectCertificates),
		CollectProfiles:          collectionOptForState(mb.CollectProfiles),
		FileVaultConfigID:        fvCfgID,
		RecoveryPasswordConfigID: rpCfgID,
	}
}

func mdmBlueprintRequestWithState(data mdmBlueprint) *goztl.MDMBlueprintRequest {
	var fvCfgID *int
	if !data.FileVaultConfigID.IsNull() {
		fvCfgID = goztl.Int(int(data.FileVaultConfigID.ValueInt64()))
	}
	var rpCfgID *int
	if !data.RecoveryPasswordConfigID.IsNull() {
		rpCfgID = goztl.Int(int(data.RecoveryPasswordConfigID.ValueInt64()))
	}
	return &goztl.MDMBlueprintRequest{
		Name:                     data.Name.ValueString(),
		InventoryInterval:        int(data.InventoryInterval.ValueInt64()),
		CollectApps:              collectionOptWithState(data.CollectApps),
		CollectCertificates:      collectionOptWithState(data.CollectCertificates),
		CollectProfiles:          collectionOptWithState(data.CollectProfiles),
		FileVaultConfigID:        fvCfgID,
		RecoveryPasswordConfigID: rpCfgID,
	}
}
