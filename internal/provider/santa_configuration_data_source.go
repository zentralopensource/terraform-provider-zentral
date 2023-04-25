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
var _ datasource.DataSource = &SantaConfigurationDataSource{}

func NewSantaConfigurationDataSource() datasource.DataSource {
	return &SantaConfigurationDataSource{}
}

// SantaConfigurationDataSource defines the data source implementation.
type SantaConfigurationDataSource struct {
	client *goztl.Client
}

func (d *SantaConfigurationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_santa_configuration"
}

func (d *SantaConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Santa configuration to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `zentral_santa_configuration` allows details of a Santa configuration to be retrieved by its `ID` or name.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Santa configuration.",
				MarkdownDescription: "`ID` of the Santa configuration.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the Santa configuration.",
				MarkdownDescription: "Name of the Santa configuration.",
				Optional:            true,
			},
			"client_mode": schema.StringAttribute{
				Description:         "Client mode of the Santa configuration.",
				MarkdownDescription: "Client mode of the Santa configuration.",
				Computed:            true,
			},
			"client_certificate_auth": schema.BoolAttribute{
				Description:         "If `true`, mTLS is required between Santa and Zentral.",
				MarkdownDescription: "If `true`, mTLS is required between Santa and Zentral.",
				Computed:            true,
			},
			"batch_size": schema.Int64Attribute{
				Description:         "The number of rules to download or events to upload per request.",
				MarkdownDescription: "The number of rules to download or events to upload per request.",
				Computed:            true,
			},
			"full_sync_interval": schema.Int64Attribute{
				Description:         "The max time to wait before performing a full sync with the server.",
				MarkdownDescription: "The max time to wait before performing a full sync with the server.",
				Computed:            true,
			},
			"enable_bundles": schema.BoolAttribute{
				Description:         "If set to true the bundle scanning feature is enabled.",
				MarkdownDescription: "If set to `true` the bundle scanning feature is enabled.",
				Computed:            true,
			},
			"enable_transitive_rules": schema.BoolAttribute{
				Description:         "If set to true the transitive rule feature is enabled.",
				MarkdownDescription: "If set to `true` the transitive rule feature is enabled.",
				Computed:            true,
			},
			"allowed_path_regex": schema.StringAttribute{
				Description:         "A regex to allow if the binary, certificate, or Team ID scopes did not allow/block execution.",
				MarkdownDescription: "A regex to allow if the binary, certificate, or Team ID scopes did not allow/block execution.",
				Computed:            true,
			},
			"blocked_path_regex": schema.StringAttribute{
				Description:         "A regex to block if the binary, certificate, or Team ID scopes did not allow/block an execution.",
				MarkdownDescription: "A regex to block if the binary, certificate, or Team ID scopes did not allow/block an execution.",
				Computed:            true,
			},
			"block_usb_mount": schema.BoolAttribute{
				Description:         "If set to true blocking USB Mass storage feature is enabled.",
				MarkdownDescription: "If set to `true` blocking USB Mass storage feature is enabled.",
				Computed:            true,
			},
			"remount_usb_mode": schema.SetAttribute{
				Description:         "Array of strings for arguments to pass to mount -o.",
				MarkdownDescription: "Array of strings for arguments to pass to `mount -o`.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"allow_unknown_shard": schema.Int64Attribute{
				Description:         "Restrict the reporting of 'Allow Unknown' events to a percentage (0-100) of hosts.",
				MarkdownDescription: "Restrict the reporting of 'Allow Unknown' events to a percentage (0-100) of hosts.",
				Computed:            true,
			},
			"enable_all_event_upload_shard": schema.Int64Attribute{
				Description:         "Restrict the upload of all execution events to Zentral, including those that were explicitly allowed, to a percentage (0-100) of hosts",
				MarkdownDescription: "Restrict the upload of all execution events to Zentral, including those that were explicitly allowed, to a percentage (0-100) of hosts",
				Computed:            true,
			},
			"sync_incident_severity": schema.Int64Attribute{
				Description:         "If 100, 200, 300, incidents will be automatically opened and closed when the santa agent rules are out of sync.",
				MarkdownDescription: "If 100, 200, 300, incidents will be automatically opened and closed when the santa agent rules are out of sync.",
				Computed:            true,
			},
		},
	}
}

func (d *SantaConfigurationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SantaConfigurationDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data santaConfiguration
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_santa_configuration` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_santa_configuration` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *SantaConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data santaConfiguration

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlSC *goztl.SantaConfiguration
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlSC, _, err = d.client.SantaConfigurations.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Santa configuration '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlSC, _, err = d.client.SantaConfigurations.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Santa configuration '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlSC != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, santaConfigurationForState(ztlSC))...)
	}
}
