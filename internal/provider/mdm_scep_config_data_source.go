package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MDMSCEPConfigDataSource{}

func NewMDMSCEPConfigDataSource() datasource.DataSource {
	return &MDMSCEPConfigDataSource{}
}

// MDMSCEPConfigDataSource defines the data source implementation.
type MDMSCEPConfigDataSource struct {
	client *goztl.Client
}

func (d *MDMSCEPConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_scep_config"
}

func (d *MDMSCEPConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM SCEP config to be retrieved by its ID or its name.",
		MarkdownDescription: "The data source `zentral_mdm_scep_config` allows details of a MDM SCEP config to be retrieved by its `ID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the MDM SCEP config.",
				MarkdownDescription: "`ID` of the MDM SCEP config.",
				Optional:            true,
			},
			"provisioning_uid": schema.StringAttribute{
				Description:         "Provisioning UID of the SCEP config.",
				MarkdownDescription: "Provisioning `UID` of the SCEP config.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the SCEP config.",
				MarkdownDescription: "Name of the SCEP config.",
				Optional:            true,
			},
			"url": schema.StringAttribute{
				Description:         "URL of the SCEP server.",
				MarkdownDescription: "URL of the SCEP server.",
				Computed:            true,
			},
			"key_usage": schema.Int64Attribute{
				Description:         "Key usage for the SCEP requests.",
				MarkdownDescription: "Key usage for the SCEP requests.",
				Computed:            true,
			},
			"key_is_extractable": schema.BoolAttribute{
				Description:         "Indicates if the private key of the resulting certificate is extractable.",
				MarkdownDescription: "Indicates if the private key of the resulting certificate is extractable.",
				Computed:            true,
			},
			"key_size": schema.Int64Attribute{
				Description:         "The size of the private key in bits.",
				MarkdownDescription: "The size of the private key in bits.",
				Computed:            true,
			},
			"allow_all_apps_access": schema.BoolAttribute{
				Description:         "Indicates if the private key can be used by all apps on the device.",
				MarkdownDescription: "Indicates if the private key can be used by all apps on the device.",
				Computed:            true,
			},
		},
	}
}

func (d *MDMSCEPConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MDMSCEPConfigDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmSCEPConfig
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_scep_config` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_scep_config` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MDMSCEPConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmSCEPConfig

	// Read Terraform SCEP config data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMSC *goztl.MDMSCEPConfig
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlMSC, _, err = d.client.MDMSCEPConfigs.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM SCEP config '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlMSC, _, err = d.client.MDMSCEPConfigs.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM SCEP config '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMSC != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmSCEPConfigForState(ztlMSC))...)
	}
}
