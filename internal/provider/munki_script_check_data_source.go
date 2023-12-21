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
var _ datasource.DataSource = &MunkiScriptCheckDataSource{}

func NewMunkiScriptCheckDataSource() datasource.DataSource {
	return &MunkiScriptCheckDataSource{}
}

// MunkiScriptCheckDataSource defines the data source implementation.
type MunkiScriptCheckDataSource struct {
	client *goztl.Client
}

func (d *MunkiScriptCheckDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_munki_script_check"
}

func (d *MunkiScriptCheckDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Munki script check to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `zentral_munki_script_check` allows details of a Munki script check to be retrieved by its `ID` or name.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Munki script check.",
				MarkdownDescription: "`ID` of the Munki script check.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the Munki script check.",
				MarkdownDescription: "Name of the Munki script check.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the Munki script check.",
				MarkdownDescription: "Description of the Munki script check.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				Description:         "Type of the script check.",
				MarkdownDescription: "Type of the script check.",
				Computed:            true,
			},
			"source": schema.StringAttribute{
				Description:         "Source of the Munki script check.",
				MarkdownDescription: "Source of the Munki script check.",
				Computed:            true,
			},
			"expected_result": schema.StringAttribute{
				Description:         "Expected result of the Munki script check.",
				MarkdownDescription: "Expected result of the Munki script check.",
				Computed:            true,
			},
			"arch_amd64": schema.BoolAttribute{
				Description:         "If true, this Munki script check is scheduled on Intel machines.",
				MarkdownDescription: "If `true`, this Munki script check is scheduled on Intel machines.",
				Computed:            true,
			},
			"arch_arm64": schema.BoolAttribute{
				Description:         "If true, this Munki script check is scheduled on Apple Silicon machines.",
				MarkdownDescription: "If `true`, this Munki script check is scheduled on Apple Silicon machines.",
				Computed:            true,
			},
			"min_os_version": schema.StringAttribute{
				Description:         "This Munki script check is scheduled on machines with an OS version higher or equal to this value.",
				MarkdownDescription: "This Munki script check is scheduled on machines with an OS version higher or equal to this value.",
				Computed:            true,
			},
			"max_os_version": schema.StringAttribute{
				Description:         "This Munki script check is scheduled on machines with an OS version lower than this value.",
				MarkdownDescription: "This Munki script check is scheduled on machines with an OS version lower than this value.",
				Computed:            true,
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags this Munki script check is restricted to.",
				MarkdownDescription: "The IDs of the tags this Munki script check is restricted to.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"excluded_tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags this Munki script check is not scoped to.",
				MarkdownDescription: "The IDs of the tags this Munki script check is not scoped to.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				Description:         "Version of the Munki script check.",
				MarkdownDescription: "Version of the Munki script check.",
				Computed:            true,
			},
		},
	}
}

func (d *MunkiScriptCheckDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MunkiScriptCheckDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data munkiScriptCheck
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_munki_script_check` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_munki_script_check` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MunkiScriptCheckDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data munkiScriptCheck

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMSC *goztl.MunkiScriptCheck
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlMSC, _, err = d.client.MunkiScriptChecks.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Munki script check '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlMSC, _, err = d.client.MunkiScriptChecks.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Munki script check '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMSC != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, munkiScriptCheckForState(ztlMSC))...)
	}
}
