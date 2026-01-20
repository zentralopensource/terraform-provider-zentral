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
var _ datasource.DataSource = &MDMDEPEnrollmentDataSource{}

func NewMDMDEPEnrollmentDataSource() datasource.DataSource {
	return &MDMDEPEnrollmentDataSource{}
}

// MDMDEPEnrollmentDataSource defines the data source implementation.
type MDMDEPEnrollmentDataSource struct {
	client *goztl.Client
}

func (d *MDMDEPEnrollmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_dep_enrollment"
}

func (d *MDMDEPEnrollmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM DEP enrollment to be retrieved by its ID and name.",
		MarkdownDescription: "The data source `zentral_mdm_ota_enrollment` allows details of a MDM DEP enrollment to be retrieved by its `ID` and `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the MDM DEP enrollment.",
				MarkdownDescription: "`ID` of the MDM DEP enrollment.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the MDM DEP enrollment.",
				MarkdownDescription: "Name of the MDM DEP enrollment.",
				Optional:            true,
			},
			"display_name": schema.StringAttribute{
				Description:         "Name of the MDM DEP enrollment as displayed on the device.",
				MarkdownDescription: "Name of the MDM DEP enrollment as displayed on the device.",
				Computed:            true,
			},
			"use_realm_user": schema.BoolAttribute{
				Description:         "Toggles if the realm user usage is assigned with the DEP enrollment.",
				MarkdownDescription: "Toggles if the realm user usage is assigned with the DEP enrollment.",
				Computed:            true,
			},
			"username_pattern": schema.StringAttribute{
				Description:         "Allowed username pattern assigned the DEP enrollment.",
				MarkdownDescription: "Allowed username pattern assigned the DEP enrollment.",
				Computed:            true,
			},
			"realm_user_is_admin": schema.BoolAttribute{
				Description:         "Toggles if the created user is admin with the DEP enrollment.",
				MarkdownDescription: "Toggles if the created user is admin with the DEP enrollment.",
				Computed:            true,
			},
			"admin_full_name": schema.StringAttribute{
				Description:         "Full name of the admin created by the DEP enrollment.",
				MarkdownDescription: "Full name of the admin created by the DEP enrollment.",
				Optional:            true,
			},
			"admin_short_name": schema.StringAttribute{
				Description:         "Short name of the admin created by the DEP enrollment.",
				MarkdownDescription: "Short name of the admin created by the DEP enrollment.",
				Optional:            true,
			},
			"hidden_admin": schema.BoolAttribute{
				Description:         "Toggles if the created admin is hidden with the DEP enrollment.",
				MarkdownDescription: "Toggles if the created admin is hidden with the DEP enrollment.",
				Computed:            true,
			},
			"admin_password_complexity": schema.Int64Attribute{
				Description:         "Required password complexity for the DEP enrollment.",
				MarkdownDescription: "Required password complexity for the DEP enrollment.",
				Computed:            true,
			},
			"admin_password_rotation_delay": schema.Int64Attribute{
				Description:         "Delay for the password rotation for the DEP enrollment.",
				MarkdownDescription: "Delay for the password rotation for the DEP enrollment.",
				Computed:            true,
			},
			"allow_pairing": schema.BoolAttribute{
				Description:         "Toggles if the pairing is allowed with the DEP enrollment.",
				MarkdownDescription: "Toggles if the pairing is allowed with the DEP enrollment.",
				Computed:            true,
			},
			"auto_advance_setup": schema.BoolAttribute{
				Description:         "Toggles if the auto advance setup is used with the DEP enrollment.",
				MarkdownDescription: "Toggles if the auto advance setup is used with the DEP enrollment.",
				Computed:            true,
			},
			"await_device_configured": schema.BoolAttribute{
				Description:         "Toggles if await device configured is used with the DEP enrollment.",
				MarkdownDescription: "Toggles if await device configured is used with the DEP enrollment.",
				Computed:            true,
			},
			"department": schema.StringAttribute{
				Description:         "Department asigned with the DEP enrollment.",
				MarkdownDescription: "Department asigned with the DEP enrollment.",
				Computed:            true,
			},
			"is_mandatory": schema.BoolAttribute{
				Description:         "Toggles if the DEP enrollment is mandatory.",
				MarkdownDescription: "Toggles if the DEP enrollment is mandatory.",
				Computed:            true,
			},
			"is_mdm_removable": schema.BoolAttribute{
				Description:         "Toggles if the DEP enrollment is removable via mdm.",
				MarkdownDescription: "Toggles if the DEP enrollment is removable via mdm.",
				Computed:            true,
			},
			"is_multi_user": schema.BoolAttribute{
				Description:         "Toggles if the DEP enrollment allows multi users.",
				MarkdownDescription: "Toggles if the DEP enrollment allows multi users.",
				Computed:            true,
			},
			"is_supervised": schema.BoolAttribute{
				Description:         "Toggles if the DEP enrollment is supervised.",
				MarkdownDescription: "Toggles if the DEP enrollment is supervised.",
				Computed:            true,
			},
			"language": schema.StringAttribute{
				Description:         "Language asigned with the DEP enrollment.",
				MarkdownDescription: "Language asigned with the DEP enrollment.",
				Computed:            true,
			},
			"org_magic": schema.StringAttribute{
				Description:         "Org magic asigned with the DEP enrollment.",
				MarkdownDescription: "Org magic asigned with the DEP enrollment.",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				Description:         "Region asigned with the DEP enrollment.",
				MarkdownDescription: "Region asigned with the DEP enrollment.",
				Computed:            true,
			},
			"skip_setup_items": schema.SetAttribute{
				Description:         "Set of to be skiped setups assigned with the DEP enrollment.",
				MarkdownDescription: "Set of to be skiped setups assigned with the DEP enrollment.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"support_email_address": schema.StringAttribute{
				Description:         "Support email address asigned with the DEP enrollment.",
				MarkdownDescription: "Support email address asigned with the DEP enrollment.",
				Computed:            true,
			},
			"support_phone_number": schema.StringAttribute{
				Description:         "Support phone number asigned with the DEP enrollment.",
				MarkdownDescription: "Support phone number asigned with the DEP enrollment.",
				Computed:            true,
			},
			"include_tls_certificates": schema.BoolAttribute{
				Description:         "Toggles if the DEP enrollment include tls certificates.",
				MarkdownDescription: "Toggles if the DEP enrollment include tls certificates.",
				Computed:            true,
			},
			"ios_max_version": schema.StringAttribute{
				Description:         "Max iOS version enforced with this DEP enrollment.",
				MarkdownDescription: "Max iOS version enforced with this DEP enrollment.",
				Computed:            true,
			},
			"ios_min_version": schema.StringAttribute{
				Description:         "Min iOS version enforced with this DEP enrollment.",
				MarkdownDescription: "Min iOS version enforced with this DEP enrollment.",
				Computed:            true,
			},
			"macos_max_version": schema.StringAttribute{
				Description:         "Max macOS version enforced with this DEP enrollment.",
				MarkdownDescription: "Max macOS version enforced with this DEP enrollment.",
				Computed:            true,
			},
			"macos_min_version": schema.StringAttribute{
				Description:         "Min macOS version enforced with this DEP enrollment.",
				MarkdownDescription: "Min macOS version enforced with this DEP enrollment.",
				Computed:            true,
			},
			"blueprint_id": schema.Int64Attribute{
				Description:         "ID of the MDM blueprint linked to the DEP enrollment.",
				MarkdownDescription: "`ID` of the MDM blueprint linked to the DEP enrollment.",
				Computed:            true,
			},
			"push_certificate_id": schema.Int64Attribute{
				Description:         "ID of the MDM push certificate linked to the DEP enrollment.",
				MarkdownDescription: "`ID` of the MDM push certificate linked to the DEP enrollment.",
				Computed:            true,
			},
			"realm_uuid": schema.StringAttribute{
				Description:         "UUID of the identity realm linked to the DEP enrollment.",
				MarkdownDescription: "`UUID` of the identity realm linked to the DEP enrollment.",
				Optional:            true,
			},
			"acme_issuer_id": schema.StringAttribute{
				Description:         "ID of the optional MDM ACME issuer linked to the DEP enrollment.",
				MarkdownDescription: "`ID` of the optional MDM ACME issuer linked to the ODEPenrollment.",
				Optional:            true,
			},
			"scep_issuer_id": schema.StringAttribute{
				Description:         "ID of the MDM SCEP issuer linked to the DEP enrollment.",
				MarkdownDescription: "`ID` of the MDM SCEP issuer linked to the DEP enrollment.",
				Computed:            true,
			},
			"virtual_server_id": schema.Int64Attribute{
				Description:         "ID of the MDM virtual server linked to the DEP enrollment.",
				MarkdownDescription: "`ID` of the MDM virtual server linked to the DEP enrollment.",
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
				Optional:            true,
			},
		},
	}
}

func (d *MDMDEPEnrollmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MDMDEPEnrollmentDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmDEPEnrollment
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

func (d *MDMDEPEnrollmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmDEPEnrollment

	// Read Terraform enrollment data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlEnrollment *goztl.MDMDEPEnrollment
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlEnrollment, _, err = d.client.MDMDEPEnrollments.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM DEP enrollment '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlEnrollment, _, err = d.client.MDMDEPEnrollments.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM DEP enrollment '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlEnrollment != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmDEPEnrollmentForState(ztlEnrollment))...)
	}
}
