package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type monolithManifest struct {
	ID                 types.Int64  `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	MetaBusinessUnitID types.Int64  `tfsdk:"meta_business_unit_id"`
	Version            types.Int64  `tfsdk:"version"`
}

func monolithManifestForState(mm *goztl.MonolithManifest) monolithManifest {
	return monolithManifest{
		ID:                 types.Int64Value(int64(mm.ID)),
		Name:               types.StringValue(mm.Name),
		MetaBusinessUnitID: types.Int64Value(int64(mm.MetaBusinessUnitID)),
		Version:            types.Int64Value(int64(mm.Version)),
	}
}

func monolithManifestRequestWithState(data monolithManifest) *goztl.MonolithManifestRequest {
	return &goztl.MonolithManifestRequest{
		Name:               data.Name.ValueString(),
		MetaBusinessUnitID: int(data.MetaBusinessUnitID.ValueInt64()),
	}
}
