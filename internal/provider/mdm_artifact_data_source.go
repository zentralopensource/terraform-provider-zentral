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
var _ datasource.DataSource = &MDMArtifactDataSource{}

func NewMDMArtifactDataSource() datasource.DataSource {
	return &MDMArtifactDataSource{}
}

// MDMArtifactDataSource defines the data source implementation.
type MDMArtifactDataSource struct {
	client *goztl.Client
}

func (d *MDMArtifactDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_artifact"
}

func (d *MDMArtifactDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM artifact to be retrieved by its ID or its name.",
		MarkdownDescription: "The data source `zentral_mdm_artifact` allows details of a MDM artifact to be retrieved by its `ID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the MDM artifact.",
				MarkdownDescription: "`ID` of the MDM artifact.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the artifact.",
				MarkdownDescription: "Name of the artifact.",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				Description:         "Type of the artifact.",
				MarkdownDescription: "Type of the artifact.",
				Computed:            true,
			},
			"channel": schema.StringAttribute{
				Description:         "Channel of the artifact.",
				MarkdownDescription: "Channel of the artifact.",
				Computed:            true,
			},
			"platforms": schema.SetAttribute{
				Description:         "Platforms of the artifact.",
				MarkdownDescription: "Platforms of the artifact.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"install_during_setup_assistant": schema.BoolAttribute{
				Description:         "If true, this artifact will be installed during the setup assistant.",
				MarkdownDescription: "If `true`, this artifact will be installed during the setup assistant.",
				Computed:            true,
			},
			"auto_update": schema.BoolAttribute{
				Description:         "If true, new version of this artifact will be automatically installed.",
				MarkdownDescription: "If `true`, new version of this artifact will be automatically installed.",
				Computed:            true,
			},
			"reinstall_interval": schema.Int64Attribute{
				Description:         "In days, the time interval after which the artifact will be reinstalled. If 0, the artifact will not be reinstalled.",
				MarkdownDescription: "In days, the time interval after which the artifact will be reinstalled. If `0`, the artifact will not be reinstalled.",
				Computed:            true,
			},
			"reinstall_on_os_update": schema.StringAttribute{
				Description:         "Possible values: No, Major, Minor, Patch. Defaults to No.",
				MarkdownDescription: "Possible values: `No`, `Major`, `Minor`, `Patch`. Defaults to `No`.",
				Computed:            true,
			},
			"requires": schema.SetAttribute{
				Description:         "IDs of the artifacts required by this artifact.",
				MarkdownDescription: "`ID`s of the artifacts required by this artifact.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *MDMArtifactDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MDMArtifactDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmArtifact
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_artifact` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_artifact` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MDMArtifactDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmArtifact

	// Read Terraform artifact data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMA *goztl.MDMArtifact
	var err error
	if !data.ID.IsNull() {
		ztlMA, _, err = d.client.MDMArtifacts.GetByID(ctx, data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM artifact '%s' by ID, got error: %s", data.ID.ValueString(), err),
			)
		}
	} else {
		ztlMA, _, err = d.client.MDMArtifacts.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM artifact '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMA != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmArtifactForState(ztlMA))...)
	}
}
