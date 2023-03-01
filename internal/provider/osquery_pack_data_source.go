package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &OsqueryPackDataSource{}

func NewOsqueryPackDataSource() datasource.DataSource {
	return &OsqueryPackDataSource{}
}

// OsqueryPackDataSource defines the data source implementation.
type OsqueryPackDataSource struct {
	client *goztl.Client
}

func (d *OsqueryPackDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_osquery_pack"
}

func (d *OsqueryPackDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Osquery pack to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `zentral_osquery_pack` allows details of a Osquery pack to be retrieved by its `ID` or name.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the pack.",
				MarkdownDescription: "`ID` of the pack.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the pack.",
				MarkdownDescription: "Name of the pack.",
				Optional:            true,
			},
			"slug": schema.StringAttribute{
				Description:         "Slug of the pack.",
				MarkdownDescription: "Slug of the pack.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the pack.",
				MarkdownDescription: "Description of the pack.",
				Computed:            true,
			},
			"discovery_queries": schema.ListAttribute{
				Description:         "List of osquery queries which control whether or not the pack will execute.",
				MarkdownDescription: "List of osquery queries which control whether or not the pack will execute.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"shard": schema.Int64Attribute{
				Description:         "Restrict the pack to a percentage (1-100) of target hosts.",
				MarkdownDescription: "Restrict the pack to a percentage (1-100) of target hosts.",
				Computed:            true,
			},
			"event_routing_key": schema.StringAttribute{
				Description:         "Routing key added to the metadata of the events that the queries of this pack generate.",
				MarkdownDescription: "Routing key added to the metadata of the events that the queries of this pack generate.",
				Computed:            true,
			},
		},
	}
}

func (d *OsqueryPackDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OsqueryPackDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data osqueryPack
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_osquery_pack` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_osquery_pack` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *OsqueryPackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data osqueryPack

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlOP *goztl.OsqueryPack
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlOP, _, err = d.client.OsqueryPacks.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Osquery pack '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlOP, _, err = d.client.OsqueryPacks.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Osquery pack '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlOP != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, osqueryPackForState(ztlOP))...)
	}
}
