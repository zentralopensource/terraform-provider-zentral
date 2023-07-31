package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MDMFileVaultConfigDataSource{}

func NewMDMFileVaultConfigDataSource() datasource.DataSource {
	return &MDMFileVaultConfigDataSource{}
}

// MDMFileVaultConfigDataSource defines the data source implementation.
type MDMFileVaultConfigDataSource struct {
	client *goztl.Client
}

func (d *MDMFileVaultConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_filevault_config"
}

func (d *MDMFileVaultConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM FileVault configuration to be retrieved by its ID or its name.",
		MarkdownDescription: "The data source `zentral_mdm_filevault_config` allows details of a MDM FileVault configuration to be retrieved by its `ID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the MDM FileVault configuration.",
				MarkdownDescription: "`ID` of the MDM FileVault configuration.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the FileVault configuration.",
				MarkdownDescription: "Name of the FileVault configuration.",
				Optional:            true,
			},
			"escrow_location_display_name": schema.StringAttribute{
				Description:         "Description of the location where the FDE PRK will be escrowed. This text will be inserted into the message the user sees when enabling FileVault.",
				MarkdownDescription: "Description of the location where the FDE PRK will be escrowed. This text will be inserted into the message the user sees when enabling FileVault.",
				Computed:            true,
			},
			"at_login_only": schema.BoolAttribute{
				Description:         "If true, prevents requests for enabling FileVault at user logout time. Defaults to false.",
				MarkdownDescription: "If `true`, prevents requests for enabling FileVault at user logout time. Defaults to `false`.",
				Computed:            true,
			},
			"bypass_attempts": schema.Int64Attribute{
				Description:         "The maximum number of times users can bypass enabling FileVault before being required to enable it to log in.",
				MarkdownDescription: "The maximum number of times users can bypass enabling FileVault before being required to enable it to log in.",
				Computed:            true,
			},
			"show_recovery_key": schema.BoolAttribute{
				Description:         "If false, prevents display of the personal recovery key to the user after FileVault is enabled. Defaults to false.",
				MarkdownDescription: "If `false`, prevents display of the personal recovery key to the user after FileVault is enabled. Defaults to `false`.",
				Computed:            true,
			},
			"destroy_key_on_standby": schema.BoolAttribute{
				Description:         "Set to true to prevent storing the FileVault key across restarts. Defaults to false.",
				MarkdownDescription: "Set to `true` to prevent storing the FileVault key across restarts. Defaults to `false`.",
				Computed:            true,
			},
			"prk_rotation_interval_days": schema.Int64Attribute{
				Description:         "The automatic PRK rotation interval in days. It has a maximum value of 365. Defaults to 0 (no automatic rotation).",
				MarkdownDescription: "The automatic PRK rotation interval in days. It has a maximum value of `365`. Defaults to `0` (no automatic rotation).",
				Computed:            true,
			},
		},
	}
}

func (d *MDMFileVaultConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MDMFileVaultConfigDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmFileVaultConfig
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_filevault_config` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_filevault_config` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MDMFileVaultConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmFileVaultConfig

	// Read Terraform FileVault configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMFC *goztl.MDMFileVaultConfig
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlMFC, _, err = d.client.MDMFileVaultConfigs.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM FileVault configuration '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlMFC, _, err = d.client.MDMFileVaultConfigs.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM FileVault configuration '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMFC != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmFileVaultConfigForState(ztlMFC))...)
	}
}
