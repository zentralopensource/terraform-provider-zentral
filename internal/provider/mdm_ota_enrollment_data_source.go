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
var _ datasource.DataSource = &MDMOTAEnrollmentDataSource{}

func NewMDMOTAEnrollmentDataSource() datasource.DataSource {
	return &MDMOTAEnrollmentDataSource{}
}

// MDMOTAEnrollmentDataSource defines the data source implementation.
type MDMOTAEnrollmentDataSource struct {
	client *goztl.Client
}

func (d *MDMOTAEnrollmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_ota_enrollment"
}

func (d *MDMOTAEnrollmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM OTA enrollment to be retrieved by its ID and name.",
		MarkdownDescription: "The data source `zentral_mdm_ota_enrollment` allows details of a MDM OTA enrollment to be retrieved by its `ID` and `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the MDM OTA enrollment.",
				MarkdownDescription: "`ID` of the MDM OTA enrollment.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the MDM OTA enrollment.",
				MarkdownDescription: "Name of the MDM OTA enrollment.",
				Optional:            true,
			},
			"display_name": schema.StringAttribute{
				Description:         "Name of the MDM OTA enrollment as displayed on the device.",
				MarkdownDescription: "Name of the MDM OTA enrollment as displayed on the device.",
				Computed:            true,
			},
			"blueprint_id": schema.Int64Attribute{
				Description:         "ID of the MDM blueprint linked to the OTA enrollment.",
				MarkdownDescription: "`ID` of the MDM blueprint linked to the OTA enrollment.",
				Computed:            true,
			},
			"push_certificate_id": schema.Int64Attribute{
				Description:         "ID of the MDM push certificate linked to the OTA enrollment.",
				MarkdownDescription: "`ID` of the MDM push certificate linked to the OTA enrollment.",
				Computed:            true,
			},
			"realm_id": schema.Int64Attribute{
				Description:         "ID of the identity realm linked to the OTA enrollment.",
				MarkdownDescription: "`ID` of the identity realm linked to the OTA enrollment.",
				Computed:            true,
			},
			"scep_config_id": schema.Int64Attribute{
				Description:         "ID of the MDM SCEP configuration linked to the OTA enrollment.",
				MarkdownDescription: "`ID` of the MDM SCEP configuration linked to the OTA enrollment.",
				Computed:            true,
			},
			"scep_verification": schema.BoolAttribute{
				Description:         "Indicates if a SCEP verification is expected during the enrollment.",
				MarkdownDescription: "Indicates if a SCEP verification is expected during the enrollment.",
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

func (d *MDMOTAEnrollmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MDMOTAEnrollmentDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmOTAEnrollment
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_ota_enrollment` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_ota_enrollment` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MDMOTAEnrollmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmOTAEnrollment

	// Read Terraform enrollment data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMOE *goztl.MDMOTAEnrollment
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlMOE, _, err = d.client.MDMOTAEnrollments.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM OTA enrollment '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlMOE, _, err = d.client.MDMOTAEnrollments.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM OTA enrollment '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMOE != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmOTAEnrollmentForState(ztlMOE))...)
	}
}
