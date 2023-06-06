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
var _ datasource.DataSource = &OsqueryQueryDataSource{}

func NewOsqueryQueryDataSource() datasource.DataSource {
	return &OsqueryQueryDataSource{}
}

// OsqueryQueryDataSource defines the data source implementation.
type OsqueryQueryDataSource struct {
	client *goztl.Client
}

func (d *OsqueryQueryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_osquery_query"
}

func (d *OsqueryQueryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Osquery query to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `zentral_osquery_query` allows details of a Osquery query to be retrieved by its `ID` or name.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the query.",
				MarkdownDescription: "`ID` of the query.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the query.",
				MarkdownDescription: "Name of the query.",
				Optional:            true,
			},
			"sql": schema.StringAttribute{
				Description:         "The SQL query to run.",
				MarkdownDescription: "The SQL query to run.",
				Computed:            true,
			},
			"platforms": schema.SetAttribute{
				Description:         "Restrict the query to some platforms, default is 'all' platforms",
				MarkdownDescription: "Restrict the query to some platforms, default is 'all' platforms",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"minimum_osquery_version": schema.StringAttribute{
				Description:         "Only run on Osquery versions greater than or equal-to this version string",
				MarkdownDescription: "Only run on Osquery versions greater than or equal-to this version string",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the query.",
				MarkdownDescription: "Description of the query.",
				Computed:            true,
			},
			"value": schema.StringAttribute{
				Description:         "Description of the results returned by the query.",
				MarkdownDescription: "Description of the results returned by the query.",
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				Description:         "Version of the query.",
				MarkdownDescription: "Version of the query.",
				Computed:            true,
			},
			"compliance_check_enabled": schema.BoolAttribute{
				Description:         "If true, the query will be used as compliance check. Defaults to false.",
				MarkdownDescription: "If `true`, the query will be used as compliance check. Defaults to `false`.",
				Computed:            true,
			},
			"scheduling": schema.SingleNestedAttribute{
				Description:         "Attributes to link a query to a pack for scheduling.",
				MarkdownDescription: "Attributes to link a query to a pack for scheduling.",
				Attributes: map[string]schema.Attribute{
					"pack_id": schema.Int64Attribute{
						Description:         "The ID of the pack.",
						MarkdownDescription: "The `ID` of the pack.",
						Computed:            true,
					},
					"interval": schema.Int64Attribute{
						Description:         "the query frequency, in seconds. It has a maximum value of 604,800 (1 week).",
						MarkdownDescription: "the query frequency, in seconds. It has a maximum value of 604,800 (1 week).",
						Computed:            true,
					},
					"log_removed_actions": schema.BoolAttribute{
						Description:         "If true, removed actions should be logged. Default to true.",
						MarkdownDescription: "If `true`, remove actions should be logged. Default to `true`.",
						Computed:            true,
					},
					"snapshot_mode": schema.BoolAttribute{
						Description:         "If true, differentials will not be stored and this query will not emulate an event stream. Defaults to false.",
						MarkdownDescription: "If `true`, differentials will not be stored and this query will not emulate an event stream. Defaults to `false`.",
						Computed:            true,
					},
					"shard": schema.Int64Attribute{
						Description:         "Restrict this query to a percentage (1-100) of target hosts.",
						MarkdownDescription: "Restrict this query to a percentage (1-100) of target hosts.",
						Computed:            true,
						Optional:            true,
					},
					"can_be_denylisted": schema.BoolAttribute{
						Description:         "If true, this query can be denylisted when stopped by the watchdog for excessive resource consumption. Defaults to true.",
						MarkdownDescription: "If `true`, this query can be denylisted when stopped by the watchdog for excessive resource consumption. Defaults to `true`.",
						Computed:            true,
					},
				},
				Computed: true,
				Optional: true,
			},
		},
	}
}

func (d *OsqueryQueryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OsqueryQueryDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data osqueryQuery
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_osquery_query` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_osquery_query` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *OsqueryQueryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data osqueryQuery

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlOQ *goztl.OsqueryQuery
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlOQ, _, err = d.client.OsqueryQueries.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Osquery query '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlOQ, _, err = d.client.OsqueryQueries.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Osquery query '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlOQ != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, osqueryQueryForState(ztlOQ))...)
	}
}
