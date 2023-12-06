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
var _ datasource.DataSource = &MDMSoftwareUpdateEnforcementDataSource{}

func NewMDMSoftwareUpdateEnforcementDataSource() datasource.DataSource {
	return &MDMSoftwareUpdateEnforcementDataSource{}
}

// MDMSoftwareUpdateEnforcementDataSource defines the data source implementation.
type MDMSoftwareUpdateEnforcementDataSource struct {
	client *goztl.Client
}

func (d *MDMSoftwareUpdateEnforcementDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_software_update_enforcement"
}

func (d *MDMSoftwareUpdateEnforcementDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM software update enforcement to be retrieved by its ID or its name.",
		MarkdownDescription: "The data source `zentral_mdm_software_update_enforcement` allows details of a MDM software update enforcement to be retrieved by its `ID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the MDM software update enforcement.",
				MarkdownDescription: "`ID` of the MDM software update enforcement.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the software update enforcement.",
				MarkdownDescription: "Name of the software update enforcement.",
				Optional:            true,
			},
			"details_url": schema.StringAttribute{
				Description:         "The URL of a web page that shows details that the organization provides about the enforced update.",
				MarkdownDescription: "The URL of a web page that shows details that the organization provides about the enforced update.",
				Computed:            true,
			},
			"platforms": schema.SetAttribute{
				Description:         "The platforms this software update enforcement is scoped to.",
				MarkdownDescription: "The platforms this software update enforcement is scoped to.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags used to scope the software update enforcement.",
				MarkdownDescription: "The `ID`s of the tags used to scope the software update enforcement.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"os_version": schema.StringAttribute{
				Description:         "The target OS version to update the device to by the appropriate time.",
				MarkdownDescription: "The target OS version to update the device to by the appropriate time.",
				Computed:            true,
			},
			"build_version": schema.StringAttribute{
				Description:         "The target build version to update the device to by the appropriate time.",
				MarkdownDescription: "The target build version to update the device to by the appropriate time.",
				Computed:            true,
			},
			"local_datetime": schema.StringAttribute{
				Description:         "The local date time value that specifies when to force install the software update.",
				MarkdownDescription: "The local date time value that specifies when to force install the software update.",
				Computed:            true,
			},
			"max_os_version": schema.StringAttribute{
				Description:         "The maximum (excluded) target OS version to update the device to by the appropriate time.",
				MarkdownDescription: "The maximum (excluded) target OS version to update the device to by the appropriate time.",
				Computed:            true,
			},
			"delay_days": schema.Int64Attribute{
				Description:         "Number of days after a software update release before the device force installs it.",
				MarkdownDescription: "Number of days after a software update release before the device force installs it.",
				Computed:            true,
			},
			"local_time": schema.StringAttribute{
				Description:         "The local time value that specifies when to force install the software update.",
				MarkdownDescription: "The local time value that specifies when to force install the software update.",
				Computed:            true,
			},
		},
	}
}

func (d *MDMSoftwareUpdateEnforcementDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MDMSoftwareUpdateEnforcementDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmSoftwareUpdateEnforcement
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_software_update_enforcement` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_software_update_enforcement` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MDMSoftwareUpdateEnforcementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmSoftwareUpdateEnforcement

	// Read Terraform software update enforcement data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMSUE *goztl.MDMSoftwareUpdateEnforcement
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlMSUE, _, err = d.client.MDMSoftwareUpdateEnforcements.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM software update enforcement '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlMSUE, _, err = d.client.MDMSoftwareUpdateEnforcements.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM software update enforcement '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMSUE != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmSoftwareUpdateEnforcementForState(ztlMSUE))...)
	}
}
