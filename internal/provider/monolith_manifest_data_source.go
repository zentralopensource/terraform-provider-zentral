package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MonolithManifestDataSource{}

func NewMonolithManifestDataSource() datasource.DataSource {
	return &MonolithManifestDataSource{}
}

// MonolithManifestDataSource defines the data source implementation.
type MonolithManifestDataSource struct {
	client *goztl.Client
}

func (d *MonolithManifestDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monolith_manifest"
}

func (d *MonolithManifestDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Monolith manifest to be retrieved by its ID or its name.",
		MarkdownDescription: "The data source `zentral_monolith_manifest` allows details of a Monolith manifest to be retrieved by its `ID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Monolith manifest.",
				MarkdownDescription: "`ID` of the Monolith manifest.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the manifest.",
				MarkdownDescription: "Name of the manifest.",
				Optional:            true,
			},
			"meta_business_unit_id": schema.Int64Attribute{
				Description:         "The ID of the meta business unit this manifest is attached to.",
				MarkdownDescription: "The `ID` of the meta business unit this manifest is attached to.",
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				Description:         "Manifest version.",
				MarkdownDescription: "Manifest version.",
				Computed:            true,
			},
		},
	}
}

func (d *MonolithManifestDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonolithManifestDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data monolithManifest
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_monolith_manifest` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_monolith_manifest` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MonolithManifestDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data monolithManifest

	// Read Terraform manifest data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMM *goztl.MonolithManifest
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlMM, _, err = d.client.MonolithManifests.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Monolith manifest '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlMM, _, err = d.client.MonolithManifests.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Monolith manifest '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMM != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, monolithManifestForState(ztlMM))...)
	}
}
