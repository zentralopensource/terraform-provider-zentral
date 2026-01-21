package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MDMDEPEnrollmentDataSource{}

func NewMDMDEPEnrollmentCustomViewDataSource() datasource.DataSource {
	return &MDMDEPEnrollmentCustomViewDataSource{}
}

// MDMDEPEnrollmentDataSource defines the data source implementation.
type MDMDEPEnrollmentCustomViewDataSource struct {
	client *goztl.Client
}

func (d *MDMDEPEnrollmentCustomViewDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_dep_enrollment_custom_view"
}

func (d *MDMDEPEnrollmentCustomViewDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM enrollment custom view to be retrieved by its ID and name.",
		MarkdownDescription: "The data source `_mdm_enrollment_custom_view` allows details of a MDM enrollment custom view to be retrieved by its `ID` and `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the MDM DEP enrollment custom view.",
				MarkdownDescription: "ID of the MDM DEP enrollment custom view.",
				Optional:            true,
				Computed:            true,
			},
			"dep_enrollment": schema.Int64Attribute{
				Description:         "ID of MDM DEP enrollment of the MDM dep enrollment custom view.",
				MarkdownDescription: "ID of MDM DEP enrollment of the MDM dep enrollment custom view.",
				Optional:            true,
			},
			"custom_view": schema.StringAttribute{
				Description:         "ID of MDM custom view of the MDM dep enrollment custom view.",
				MarkdownDescription: "ID of MDM custom view of the MDM dep enrollment custom view.",
				Optional:            true,
			},
			"weight": schema.Int64Attribute{
				Description:         "Weight of the MDM DEP enrollment custom view.",
				MarkdownDescription: "Weight of the MDM DEP enrollment custom view.",
				Optional:            true,
			},
		},
	}
}

func (d *MDMDEPEnrollmentCustomViewDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MDMDEPEnrollmentCustomViewDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmDEPEnrollmentCutomView
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_dep_enrollment` data source",
			"`id` missing",
		)
	}
}

func (d *MDMDEPEnrollmentCustomViewDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmDEPEnrollmentCutomView

	// Read Terraform enrollment data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlCustomView *goztl.MDMDEPEnrollmentCustomView
	var err error
	if len(data.ID.ValueString()) > 0 {
		ztlCustomView, _, err = d.client.MDMDEPEnrollmentCustomViews.GetByID(ctx, data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM DEP enrollment custom view '%s' by ID, got error: %s", data.ID.ValueString(), err),
			)
		}
	}

	if ztlCustomView != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmDEPEnrollmentCustomViewForState(ztlCustomView))...)
	}
}
