package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MDMLocationDataSource{}

func NewMDMLocationDataSource() datasource.DataSource {
	return &MDMLocationDataSource{}
}

// MDMLocationDataSource defines the data source implementation.
type MDMLocationDataSource struct {
	client *goztl.Client
}

func (d *MDMLocationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_location"
}

func (d *MDMLocationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM location to be retrieved by its ID or its name.",
		MarkdownDescription: "The data source `zentral_mdm_location` allows details of a MDM location to be retrieved by its `ID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the MDM location.",
				MarkdownDescription: "`ID` of the MDM location.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the location.",
				MarkdownDescription: "Name of the location.",
				Optional:            true,
			},
			"organization_name": schema.StringAttribute{
				Description:         "Name of the location organization.",
				MarkdownDescription: "Name of the location organization.",
				Optional:            true,
			},
			"mdm_info_id": schema.StringAttribute{
				Description:         "MDM info ID of the location.",
				MarkdownDescription: "MDM info `ID` (`UUID`) of the location.",
				Optional:            true,
			},
		},
	}
}

func (d *MDMLocationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MDMLocationDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmLocation
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	nullInputsCount := 0
	if data.ID.IsNull() {
		nullInputsCount += 1
	}
	if data.MDMInfoID.IsNull() {
		nullInputsCount += 1
	}
	if data.Name.IsNull() {
		nullInputsCount += 1
	}
	if nullInputsCount == 3 {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_location` data source",
			"`id`, `mdm_info_id`, or `name` missing",
		)
	} else if nullInputsCount < 2 {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_location` data source",
			"Only one of `id`, `mdm_info_id`, and `name` can be set",
		)
	}
}

func (d *MDMLocationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmLocation

	// Read Terraform location data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlML *goztl.MDMLocation
	var err error
	if !data.ID.IsNull() {
		ztlML, _, err = d.client.MDMLocations.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM location '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else if !data.MDMInfoID.IsNull() {
		ztlML, _, err = d.client.MDMLocations.GetByMDMInfoID(ctx, data.MDMInfoID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM location '%s' by MDM info ID, got error: %s", data.MDMInfoID.ValueString(), err),
			)
		}
	} else if !data.Name.IsNull() {
		ztlML, _, err = d.client.MDMLocations.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM location '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlML != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmLocationForState(ztlML))...)
	}
}
