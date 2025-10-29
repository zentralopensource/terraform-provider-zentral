package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmLocationAsset struct {
	ID           types.Int64  `tfsdk:"id"`
	LocationID   types.Int64  `tfsdk:"location_id"`
	AssetID      types.Int64  `tfsdk:"asset_id"`
	AdamID       types.String `tfsdk:"adam_id"`
	PricingParam types.String `tfsdk:"pricing_param"`
}

func mdmLocationAssetForState(mla *goztl.MDMLocationAsset) mdmLocationAsset {
	return mdmLocationAsset{
		ID:           types.Int64Value(int64(mla.ID)),
		LocationID:   types.Int64Value(int64(mla.LocationID)),
		AssetID:      types.Int64Value(int64(mla.AssetID)),
		AdamID:       types.StringValue(mla.AdamID),
		PricingParam: types.StringValue(mla.PricingParam),
	}
}
