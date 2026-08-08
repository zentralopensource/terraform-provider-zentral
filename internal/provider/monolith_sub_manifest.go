package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type monolithSubManifest struct {
	ID                 types.Int64  `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	MetaBusinessUnitID types.Int64  `tfsdk:"meta_business_unit_id"`
}

func monolithSubManifestForState(msm *goztl.MonolithSubManifest) monolithSubManifest {
	return monolithSubManifest{
		ID:                 types.Int64Value(int64(msm.ID)),
		Name:               types.StringValue(msm.Name),
		Description:        types.StringValue(msm.Description),
		MetaBusinessUnitID: optionalInt64ForState(msm.MetaBusinessUnitID),
	}
}

func monolithSubManifestRequestWithState(data monolithSubManifest) *goztl.MonolithSubManifestRequest {
	return &goztl.MonolithSubManifestRequest{
		Name:               data.Name.ValueString(),
		Description:        data.Description.ValueString(),
		MetaBusinessUnitID: optionalIntWithState(data.MetaBusinessUnitID),
	}
}
