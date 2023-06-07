package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmBlueprint struct {
	ID                  types.Int64  `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	InventoryInterval   types.Int64  `tfsdk:"inventory_interval"`
	CollectApps         types.String `tfsdk:"collect_apps"`
	CollectCertificates types.String `tfsdk:"collect_certificates"`
	CollectProfiles     types.String `tfsdk:"collect_profiles"`
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
	return mdmBlueprint{
		ID:                  types.Int64Value(int64(mb.ID)),
		Name:                types.StringValue(mb.Name),
		InventoryInterval:   types.Int64Value(int64(mb.InventoryInterval)),
		CollectApps:         collectionOptForState(mb.CollectApps),
		CollectCertificates: collectionOptForState(mb.CollectCertificates),
		CollectProfiles:     collectionOptForState(mb.CollectProfiles),
	}
}

func mdmBlueprintRequestWithState(data mdmBlueprint) *goztl.MDMBlueprintRequest {
	return &goztl.MDMBlueprintRequest{
		Name:                data.Name.ValueString(),
		InventoryInterval:   int(data.InventoryInterval.ValueInt64()),
		CollectApps:         collectionOptWithState(data.CollectApps),
		CollectCertificates: collectionOptWithState(data.CollectCertificates),
		CollectProfiles:     collectionOptWithState(data.CollectProfiles),
	}
}
