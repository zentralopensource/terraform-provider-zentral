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
var _ datasource.DataSource = &MonolithEnrollmentDataSource{}

func NewMonolithEnrollmentDataSource() datasource.DataSource {
	return &MonolithEnrollmentDataSource{}
}

// MonolithEnrollmentDataSource defines the data source implementation.
type MonolithEnrollmentDataSource struct {
	client *goztl.Client
}

func (d *MonolithEnrollmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monolith_enrollment"
}

func (d *MonolithEnrollmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Monolith enrollment to be retrieved by its ID.",
		MarkdownDescription: "The data source `zentral_monolith_enrollment` allows details of a Monolith enrollment to be retrieved by its `ID`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Monolith enrollment.",
				MarkdownDescription: "`ID` of the Monolith enrollment.",
				Optional:            true,
			},
			"manifest_id": schema.Int64Attribute{
				Description:         "ID of the Monolith manifest.",
				MarkdownDescription: "`ID` of the Monolith manifest.",
				Computed:            true,
			},
			"configuration_profile_url": schema.StringAttribute{
				Description:         "Manifest profile download URL.",
				MarkdownDescription: "Manifest profile download URL.",
				Computed:            true,
			},
			"plist_url": schema.StringAttribute{
				Description:         "Plist download URL.",
				MarkdownDescription: "Plist download URL.",
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				Description:         "Enrollment version.",
				MarkdownDescription: "Enrollment version.",
				Computed:            true,
			},
			"secret": schema.StringAttribute{
				Description:         "Enrollment secret.",
				MarkdownDescription: "Enrollment secret.",
				Computed:            true,
			},
			"meta_business_unit_id": schema.Int64Attribute{
				Description:         "The ID of the meta business unit the machine will be assigned to at enrollment.",
				MarkdownDescription: "The `ID` of the meta business unit the machine will be assigned to at enrollment.",
				Computed:            true,
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags that the machine will get at enrollment.",
				MarkdownDescription: "The `ID`s of the tags that the machine will get at enrollment.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"serial_numbers": schema.SetAttribute{
				Description:         "The serial numbers the enrollment is restricted to.",
				MarkdownDescription: "The serial numbers the enrollment is restricted to.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"udids": schema.SetAttribute{
				Description:         "The UDIDs the enrollment is restricted to.",
				MarkdownDescription: "The `UDID`s the enrollment is restricted to.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"quota": schema.Int64Attribute{
				Description:         "The number of time the enrollment can be used.",
				MarkdownDescription: "The number of time the enrollment can be used.",
				Computed:            true,
			},
		},
	}
}

func (d *MonolithEnrollmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonolithEnrollmentDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data monolithEnrollment
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_monolith_enrollment` data source",
			"`id` missing",
		)
	}
}

func (d *MonolithEnrollmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data monolithEnrollment

	// Read Terraform enrollment data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlME *goztl.MonolithEnrollment
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlME, _, err = d.client.MonolithEnrollments.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Monolith enrollment '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	}

	if ztlME != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, monolithEnrollmentForState(ztlME))...)
	}
}
