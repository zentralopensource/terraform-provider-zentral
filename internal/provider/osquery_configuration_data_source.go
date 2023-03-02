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
var _ datasource.DataSource = &OsqueryConfigurationDataSource{}

func NewOsqueryConfigurationDataSource() datasource.DataSource {
	return &OsqueryConfigurationDataSource{}
}

// OsqueryConfigurationDataSource defines the data source implementation.
type OsqueryConfigurationDataSource struct {
	client *goztl.Client
}

func (d *OsqueryConfigurationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_osquery_configuration"
}

func (d *OsqueryConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Osquery configuration to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `zentral_osquery_configuration` allows details of a Osquery configuration to be retrieved by its `ID` or name.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Osquery configuration.",
				MarkdownDescription: "`ID` of the Osquery configuration.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the Osquery configuration.",
				MarkdownDescription: "Name of the Osquery configuration.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description or the Osquery configuration.",
				MarkdownDescription: "Description of the Osquery configuration.",
				Computed:            true,
			},
			"inventory": schema.BoolAttribute{
				Description:         "If true, Osquery is configured to collect inventory data. Defaults to true.",
				MarkdownDescription: "If `true`, Osquery is configured to collect inventory data. Defaults to `true`.",
				Computed:            true,
			},
			"inventory_apps": schema.BoolAttribute{
				Description:         "If true, Osquery is configured to collect the applications. Defaults to false.",
				MarkdownDescription: "If `true`, Osquery is configured to collect the applications. Defaults to `false`.",
				Computed:            true,
			},
			"inventory_ec2": schema.BoolAttribute{
				Description:         "If true, Osquery is configured to collect the EC2 metadata. Defaults to false.",
				MarkdownDescription: "If `true`, Osquery is configured to collect the EC2 metadata. Defaults to `false`.",
				Computed:            true,
			},
			"inventory_interval": schema.Int64Attribute{
				Description:         "Number of seconds to wait between collecting the inventory data.",
				MarkdownDescription: "Number of seconds to wait between collecting the inventory data.",
				Computed:            true,
			},
			"options": schema.MapAttribute{
				Description:         "A map of extra options to pass to Osquery in the flag file.",
				MarkdownDescription: "A map of extra options to pass to Osquery in the flag file.",
				// Options ElementType is types.StringType
				// This is much easier this way, and since the options are serialized in the flag file, this is not restrictive.
				// Non-string elements coming from the server will be converted to strings.
				ElementType: types.StringType,
				Computed:    true,
			},
			"atc_ids": schema.SetAttribute{
				Description:         "List of the IDs of the ATCs to include in this configuration.",
				MarkdownDescription: "List of the IDs of the ATCs to include in this configuration.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"file_category_ids": schema.SetAttribute{
				Description:         "List of the IDs of the file categories to include in this configuration.",
				MarkdownDescription: "List of the IDs of the file categories to include in this configuration.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
		},
	}
}

func (d *OsqueryConfigurationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OsqueryConfigurationDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data osqueryConfiguration
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

func (d *OsqueryConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data osqueryConfiguration

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlOC *goztl.OsqueryConfiguration
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlOC, _, err = d.client.OsqueryConfigurations.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Osquery configuration '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlOC, _, err = d.client.OsqueryConfigurations.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Osquery configuration '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlOC != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, osqueryConfigurationForState(ztlOC))...)
	}
}
