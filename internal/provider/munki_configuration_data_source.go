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
var _ datasource.DataSource = &MunkiConfigurationDataSource{}

func NewMunkiConfigurationDataSource() datasource.DataSource {
	return &MunkiConfigurationDataSource{}
}

// MunkiConfigurationDataSource defines the data source implementation.
type MunkiConfigurationDataSource struct {
	client *goztl.Client
}

func (d *MunkiConfigurationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_munki_configuration"
}

func (d *MunkiConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Munki configuration to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `zentral_munki_configuration` allows details of a Munki configuration to be retrieved by its `ID` or name.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Munki configuration.",
				MarkdownDescription: "`ID` of the Munki configuration.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the Munki configuration.",
				MarkdownDescription: "Name of the Munki configuration.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the Munki configuration.",
				MarkdownDescription: "Description of the Munki configuration.",
				Computed:            true,
			},
			"inventory_apps_full_info_shard": schema.Int64Attribute{
				Description:         "Percentage of machines configured to collect the full inventory apps information. Defaults to 100.",
				MarkdownDescription: "Percentage of machines configured to collect the full inventory apps information. Defaults to `100`.",
				Computed:            true,
			},
			"principal_user_detection_sources": schema.ListAttribute{
				Description:         "List of principal user detection sources.",
				MarkdownDescription: "List of principal user detection sources.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"principal_user_detection_domains": schema.SetAttribute{
				Description:         "Set of principal user detection domains.",
				MarkdownDescription: "Set of principal user detection domains.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"collected_condition_keys": schema.SetAttribute{
				Description:         "Set of the condition keys to collect.",
				MarkdownDescription: "Set of the condition keys to collect.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"managed_installs_sync_interval_days": schema.Int64Attribute{
				Description:         "Interval in days between full managed installs sync. Defaults to 7 days.",
				MarkdownDescription: "Interval in days between full managed installs sync. Defaults to 7 days.",
				Computed:            true,
			},
			"auto_reinstall_incidents": schema.BoolAttribute{
				Description:         "If true, incidents will be managed automatically when package reinstalls are observed. Defaults to false.",
				MarkdownDescription: "If `true`, incidents will be managed automatically when package reinstalls are observed. Defaults to `false`.",
				Computed:            true,
			},
			"auto_failed_install_incidents": schema.BoolAttribute{
				Description:         "If true, incidents will be managed automatically when package failed installs are observed. Defaults to false.",
				MarkdownDescription: "If `true`, incidents will be managed automatically when package failed installs are observed. Defaults to `false`.",
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				Description:         "Version of the Munki configuration.",
				MarkdownDescription: "Version of the Munki configuration.",
				Computed:            true,
			},
		},
	}
}

func (d *MunkiConfigurationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MunkiConfigurationDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data munkiConfiguration
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_munki_configuration` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_munki_configuration` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MunkiConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data munkiConfiguration

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMC *goztl.MunkiConfiguration
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlMC, _, err = d.client.MunkiConfigurations.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Munki configuration '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlMC, _, err = d.client.MunkiConfigurations.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Munki configuration '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMC != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, munkiConfigurationForState(ztlMC))...)
	}
}
