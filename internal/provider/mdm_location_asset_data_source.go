package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MDMLocationAssetDataSource{}

func NewMDMLocationAssetDataSource() datasource.DataSource {
	return &MDMLocationAssetDataSource{}
}

// MDMLocationAssetDataSource defines the data source implementation.
type MDMLocationAssetDataSource struct {
	client *goztl.Client
}

func (d *MDMLocationAssetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_location_asset"
}

func (d *MDMLocationAssetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM location asset to be retrieved by its location ID, Adam ID, and its pricing param.",
		MarkdownDescription: "The data source `zentral_mdm_location_asset` allows details of a MDM location asset to be retrieved by its location `ID`, `Adam ID`, and its pricing param.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the MDM location asset.",
				MarkdownDescription: "`ID` of the MDM location asset.",
				Computed:            true,
			},
			"location_id": schema.Int64Attribute{
				Description:         "ID of the location.",
				MarkdownDescription: "`ID` of the location.",
				Required:            true,
			},
			"asset_id": schema.Int64Attribute{
				Description:         "ID of the asset.",
				MarkdownDescription: "ID of the asset.",
				Computed:            true,
			},
			"adam_id": schema.StringAttribute{
				Description:         "Adam ID of the asset.",
				MarkdownDescription: "Adam ID of the asset.",
				Required:            true,
			},
			"pricing_param": schema.StringAttribute{
				Description:         "Pricing param of the asset.",
				MarkdownDescription: "Pricing param of the asset.",
				Required:            true,
			},
		},
	}
}

func (d *MDMLocationAssetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*goztl.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *goztl.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *MDMLocationAssetDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmLocationAsset
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *MDMLocationAssetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmLocationAsset

	// Read Terraform location data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMLA *goztl.MDMLocationAsset
	var err error
	ztlMLA, _, err = d.client.MDMLocationAssets.Get(ctx, int(data.LocationID.ValueInt64()), data.AdamID.ValueString(), data.PricingParam.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to get MDM location asset, got error: %s", err),
		)
	}

	if ztlMLA != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmLocationAssetForState(ztlMLA))...)
	}
}
