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

func NewMDMEnrollmentCustomViewDataSource() datasource.DataSource {
	return &MDMEnrollmentCustomViewDataSource{}
}

// MDMEnrollmentCustomViewDataSource defines the data source implementation.
type MDMEnrollmentCustomViewDataSource struct {
	client *goztl.Client
}

func (d *MDMEnrollmentCustomViewDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_enrollment_custom_view"
}

func (d *MDMEnrollmentCustomViewDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM DEP enrollment custom view to be retrieved by its ID and name.",
		MarkdownDescription: "The data source `zentral_mdm_dep_enrollment_custom_view` allows details of a MDM DEPenrollment custom view to be retrieved by its `ID` and `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the MDM enrollment custom view.",
				MarkdownDescription: "ID of the MDM enrollment custom view.",
				Computed:            true,
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the MDM enrollment custom view.",
				MarkdownDescription: "Name of the MDM enrollment custom view.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the MDM enrollment custom view.",
				MarkdownDescription: "Description of the MDM enrollment custom view.",
				Optional:            true,
			},
			"html": schema.StringAttribute{
				Description:         "HTML of the MDM enrollment custom view.",
				MarkdownDescription: "HTML of the MDM enrollment custom view.",
				Optional:            true,
			},
			"requires_authentication": schema.BoolAttribute{
				Description:         "Toggles if custom view requires authentication.",
				MarkdownDescription: "Toggles if custom view requires authentication.",
				Optional:            true,
			},
		},
	}
}

func (d *MDMEnrollmentCustomViewDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MDMEnrollmentCustomViewDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmEnrollmentCustomViewDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_dep_enrollment` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_dep_enrollment` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MDMEnrollmentCustomViewDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmEnrollmentCustomViewDataSourceModel

	// Read Terraform enrollment data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlEnrollment *goztl.MDMEnrollmentCustomView
	var err error
	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		ztlEnrollment, _, err = d.client.MDMEnrollmentCustomViews.GetByID(ctx, data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM enrollment custom view '%s' by ID, got error: %s", data.ID.ValueString(), err),
			)
		}
	} else {
		ztlEnrollment, _, err = d.client.MDMEnrollmentCustomViews.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM DEP enrollment '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlEnrollment != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmEnrollmentCustomViewDSForState(ztlEnrollment))...)
	}
}
