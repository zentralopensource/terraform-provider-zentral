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
var _ datasource.DataSource = &OsqueryATCDataSource{}

func NewOsqueryATCDataSource() datasource.DataSource {
	return &OsqueryATCDataSource{}
}

// OsqueryATCDataSource defines the data source implementation.
type OsqueryATCDataSource struct {
	client *goztl.Client
}

func (d *OsqueryATCDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_osquery_atc"
}

func (d *OsqueryATCDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Osquery automatic table construction to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `zentral_osquery_atc` allows details of a Osquery automatic table construction to be retrieved by its `ID` or name.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Osquery ATC.",
				MarkdownDescription: "`ID` of the Osquery ATC.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the Osquery ATC.",
				MarkdownDescription: "Name of the Osquery ATC.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the Osquery ATC.",
				MarkdownDescription: "Description of the Osquery ATC.",
				Computed:            true,
			},
			"table_name": schema.StringAttribute{
				Description:         "Name of the Osquery ATC table.",
				MarkdownDescription: "Name of the Osquery ATC table.",
				Computed:            true,
			},
			"query": schema.StringAttribute{
				Description:         "Query used to fetch the ATC data.",
				MarkdownDescription: "Query used to fetch the ATC data.",
				Computed:            true,
			},
			"path": schema.StringAttribute{
				Description:         "Path of the SQLite table on the device.",
				MarkdownDescription: "Path of the SQLite table on the device.",
				Computed:            true,
			},
			"columns": schema.ListAttribute{
				Description:         "List of the column names corresponding the the query.",
				MarkdownDescription: "List of the column names corresponding the the query.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"platforms": schema.SetAttribute{
				Description:         "Platform on which this ATC can be activated",
				MarkdownDescription: "Platform on which this ATC can be activated",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *OsqueryATCDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OsqueryATCDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data osqueryATC
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_osquery_configuration` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_osquery_configuration` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *OsqueryATCDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data osqueryATC

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlOA *goztl.OsqueryATC
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlOA, _, err = d.client.OsqueryATC.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Osquery ATC '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlOA, _, err = d.client.OsqueryATC.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Osquery ATC '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlOA != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, osqueryATCForState(ztlOA))...)
	}
}
