package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MonolithCatalogDataSource{}

func NewMonolithCatalogDataSource() datasource.DataSource {
	return &MonolithCatalogDataSource{}
}

// MonolithCatalogDataSource defines the data source implementation.
type MonolithCatalogDataSource struct {
	client *goztl.Client
}

func (d *MonolithCatalogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monolith_catalog"
}

func (d *MonolithCatalogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Monolith catalog to be retrieved by its ID or its name.",
		MarkdownDescription: "The data source `zentral_monolith_catalog` allows details of a Monolith catalog to be retrieved by its `ID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Monolith catalog.",
				MarkdownDescription: "`ID` of the Monolith catalog.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the catalog.",
				MarkdownDescription: "Name of the catalog.",
				Optional:            true,
			},
			"priority": schema.Int64Attribute{
				Description:         "Priority of the catalog.",
				MarkdownDescription: "Priority of the catalog.",
				Computed:            true,
			},
		},
	}
}

func (d *MonolithCatalogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonolithCatalogDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data monolithCatalog
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_monolith_catalog` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_monolith_catalog` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MonolithCatalogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data monolithCatalog

	// Read Terraform catalog data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMC *goztl.MonolithCatalog
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlMC, _, err = d.client.MonolithCatalogs.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Monolith catalog '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlMC, _, err = d.client.MonolithCatalogs.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Monolith catalog '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMC != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, monolithCatalogForState(ztlMC))...)
	}
}
