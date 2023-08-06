package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MDMBlueprintDataSource{}

func NewMDMBlueprintDataSource() datasource.DataSource {
	return &MDMBlueprintDataSource{}
}

// MDMBlueprintDataSource defines the data source implementation.
type MDMBlueprintDataSource struct {
	client *goztl.Client
}

func (d *MDMBlueprintDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_blueprint"
}

func (d *MDMBlueprintDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM blueprint to be retrieved by its ID or its name.",
		MarkdownDescription: "The data source `zentral_mdm_blueprint` allows details of a MDM blueprint to be retrieved by its `ID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the MDM blueprint.",
				MarkdownDescription: "`ID` of the MDM blueprint.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the blueprint.",
				MarkdownDescription: "Name of the blueprint.",
				Optional:            true,
			},
			"inventory_interval": schema.Int64Attribute{
				Description:         "In seconds, the minimum interval between two inventory collection. Minimum 4h, maximum 7d, default 1d.",
				MarkdownDescription: "In seconds, the minimum interval between two inventory collection. Minimum 4h, maximum 7d, default 1d.",
				Computed:            true,
			},
			"collect_apps": schema.StringAttribute{
				Description:         "Inventory apps collection setting. Possible values: NO, MANAGED_ONLY, ALL.",
				MarkdownDescription: "Possible values: `NO`, `MANAGED_ONLY`, `ALL`.",
				Computed:            true,
			},
			"collect_certificates": schema.StringAttribute{
				Description:         "Inventory certificates collection setting. Possible values: NO, MANAGED_ONLY, ALL.",
				MarkdownDescription: "Possible values: `NO`, `MANAGED_ONLY`, `ALL`.",
				Computed:            true,
			},
			"collect_profiles": schema.StringAttribute{
				Description:         "Inventory profiles collection setting. Possible values: NO, MANAGED_ONLY, ALL.",
				MarkdownDescription: "Possible values: `NO`, `MANAGED_ONLY`, `ALL`.",
				Computed:            true,
			},
			"filevault_config_id": schema.Int64Attribute{
				Description:         "The ID of the attached FileVault configuration.",
				MarkdownDescription: "The `ID` of the attached FileVault configuration.",
				Computed:            true,
			},
			"recovery_password_config_id": schema.Int64Attribute{
				Description:         "The ID of the attached recovery password configuration.",
				MarkdownDescription: "The `ID` of the attached recovery password configuration.",
				Computed:            true,
			},
		},
	}
}

func (d *MDMBlueprintDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MDMBlueprintDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmBlueprint
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_blueprint` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_blueprint` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MDMBlueprintDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmBlueprint

	// Read Terraform blueprint data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMB *goztl.MDMBlueprint
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlMB, _, err = d.client.MDMBlueprints.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM blueprint '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlMB, _, err = d.client.MDMBlueprints.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM blueprint '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMB != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmBlueprintForState(ztlMB))...)
	}
}
