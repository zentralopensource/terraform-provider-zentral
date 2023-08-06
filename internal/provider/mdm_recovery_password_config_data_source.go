package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MDMRecoveryPasswordConfigDataSource{}

func NewMDMRecoveryPasswordConfigDataSource() datasource.DataSource {
	return &MDMRecoveryPasswordConfigDataSource{}
}

// MDMRecoveryPasswordConfigDataSource defines the data source implementation.
type MDMRecoveryPasswordConfigDataSource struct {
	client *goztl.Client
}

func (d *MDMRecoveryPasswordConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_recovery_password_config"
}

func (d *MDMRecoveryPasswordConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM recovery password configuration to be retrieved by its ID or its name.",
		MarkdownDescription: "The data source `zentral_mdm_recovery_password_config` allows details of a MDM recovery password configuration to be retrieved by its `ID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the MDM recovery password configuration.",
				MarkdownDescription: "`ID` of the MDM recovery password configuration.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the recovery password configuration.",
				MarkdownDescription: "Name of the recovery password configuration.",
				Optional:            true,
			},
			"dynamic_password": schema.BoolAttribute{
				Description:         "If true, a unique password is generated for each device. Defaults to true.",
				MarkdownDescription: "If `true`, a unique password is generated for each device. Defaults to `true`.",
				Computed:            true,
			},
			"static_password": schema.StringAttribute{
				Description:         "The static password to set for all devices.",
				MarkdownDescription: "The  static password to set for all devices.",
				Computed:            true,
			},
			"rotation_interval_days": schema.Int64Attribute{
				Description:         "The automatic recovery password rotation interval in days. It has a maximum value of 365. Defaults to 0 (no automatic rotation).",
				MarkdownDescription: "The automatic recovery password rotation interval in days. It has a maximum value of `365`. Defaults to `0` (no automatic rotation).",
				Computed:            true,
			},
			"rotate_firmware_password": schema.BoolAttribute{
				Description:         "Set to true to rotate the firmware passwords. Defaults to false.",
				MarkdownDescription: "Set to `true` to rotate the firmware passwords. Defaults to `false`.",
				Computed:            true,
			},
		},
	}
}

func (d *MDMRecoveryPasswordConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MDMRecoveryPasswordConfigDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmRecoveryPasswordConfig
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_recovery_password_config` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_recovery_password_config` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MDMRecoveryPasswordConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmRecoveryPasswordConfig

	// Read Terraform recovery password configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMRPC *goztl.MDMRecoveryPasswordConfig
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlMRPC, _, err = d.client.MDMRecoveryPasswordConfigs.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM recovery password configuration '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlMRPC, _, err = d.client.MDMRecoveryPasswordConfigs.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM recovery password configuration '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMRPC != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmRecoveryPasswordConfigForState(ztlMRPC))...)
	}
}
