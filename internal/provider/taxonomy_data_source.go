package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &TaxonomyDataSource{}

func NewTaxonomyDataSource() datasource.DataSource {
	return &TaxonomyDataSource{}
}

// TaxonomyDataSource defines the data source implementation.
type TaxonomyDataSource struct {
	client *goztl.Client
}

func (d *TaxonomyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_taxonomy"
}

func (d *TaxonomyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a taxonomy to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `zentral_taxonomy` allows details of a taxonomy to be retrieved by its `ID` or name.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Taxonomy.",
				MarkdownDescription: "`ID` of the Taxonomy.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the Taxonomy.",
				MarkdownDescription: "Name of the Taxonomy.",
				Optional:            true,
			},
		},
	}
}

func (d *TaxonomyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TaxonomyDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data taxonomy
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_taxonomy` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_taxonomy` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *TaxonomyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data taxonomy

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlTaxonomy *goztl.Taxonomy
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlTaxonomy, _, err = d.client.Taxonomies.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get taxonomy '%d', got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlTaxonomy, _, err = d.client.Taxonomies.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get taxonomy '%s', got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlTaxonomy != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, taxonomyForState(ztlTaxonomy))...)
	}
}
