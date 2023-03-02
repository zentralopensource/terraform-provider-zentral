package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &OsqueryPackQueryDataSource{}

func NewOsqueryPackQueryDataSource() datasource.DataSource {
	return &OsqueryPackQueryDataSource{}
}

// OsqueryPackQueryDataSource defines the data source implementation.
type OsqueryPackQueryDataSource struct {
	client *goztl.Client
}

func (d *OsqueryPackQueryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_osquery_pack_query"
}

func (d *OsqueryPackQueryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Osquery pack query to be retrieved by its ID.",
		MarkdownDescription: "The data source `zentral_osquery_pack_query` allows details of a Osquery pack query to be retrieved by its `ID`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the pack query.",
				MarkdownDescription: "`ID` of the pack query.",
				Required:            true,
			},
			"pack_id": schema.Int64Attribute{
				Description:         "ID of the pack.",
				MarkdownDescription: "ID of the pack.",
				Computed:            true,
			},
			"query_id": schema.Int64Attribute{
				Description:         "ID of the query.",
				MarkdownDescription: "ID of the query.",
				Computed:            true,
			},
			"slug": schema.StringAttribute{
				Description:         "Slug of the pack query.",
				MarkdownDescription: "Slug of the pack query.",
				Computed:            true,
			},
			"interval": schema.Int64Attribute{
				Description:         "Query frequency, in seconds.",
				MarkdownDescription: "Query frequency, in seconds.",
				Computed:            true,
			},
			"log_removed_actions": schema.BoolAttribute{
				Description:         "If true, 'removed' actions should be logged.",
				MarkdownDescription: "If `true`, `removed` actions should be logged.",
				Computed:            true,
			},
			"snapshot_mode": schema.BoolAttribute{
				Description:         "If true, differentials will not be stored and this query will not emulate an event stream.",
				MarkdownDescription: "If `true`, differentials will not be stored and this query will not emulate an event stream.",
				Computed:            true,
			},
			"shard": schema.Int64Attribute{
				Description:         "Restrict this query to a percentage (1-100) of target hosts.",
				MarkdownDescription: "Restrict this query to a percentage (1-100) of target hosts.",
				Computed:            true,
			},
			"can_be_denylisted": schema.BoolAttribute{
				Description:         "If true, this query can be denylisted when stopped by the watchdog for excessive resource consumption.",
				MarkdownDescription: "If `true`, this query can be denylisted when stopped by the watchdog for excessive resource consumption.",
				Computed:            true,
			},
		},
	}
}

func (d *OsqueryPackQueryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OsqueryPackQueryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data osqueryPackQuery

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOPQ, _, err := d.client.OsqueryPackQueries.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to get Osquery pack query '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
		)
	}

	if ztlOPQ != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, osqueryPackQueryForState(ztlOPQ))...)
	}
}
