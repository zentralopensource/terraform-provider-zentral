package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMDEPEnrollmentResource{}
var _ resource.ResourceWithImportState = &MDMDEPEnrollmentResource{}

func NewMDMDEPEnrollmentResource() resource.Resource {
	return &MDMDEPEnrollmentResource{}
}

// MDMDEPEnrollmentResource defines the resource implementation.
type MDMDEPEnrollmentResource struct {
	client *goztl.Client
}

func (r *MDMDEPEnrollmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_dep_enrollment"
}

func (r *MDMDEPEnrollmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM DEP enrollments.",
		MarkdownDescription: "The resource `zentral_mdm_dep_enrollment` manages MDM DEP enrollments.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the MDM DEP enrollment.",
				MarkdownDescription: "`ID` of the MDM DEP enrollment.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the MDM DEP enrollment.",
				MarkdownDescription: "Name of the MDM DEP enrollment.",
				Required:            true,
			},
			"display_name": schema.StringAttribute{
				Description:         "Name of the MDM DEP enrollment as displayed on the device.",
				MarkdownDescription: "Name of the MDM DEP enrollment as displayed on the device.",
				Required:            true,
			},
			"use_realm_user": schema.BoolAttribute{
				Description:         "Toggles if the realm user usage is assigned with the DEP enrollment.",
				MarkdownDescription: "Toggles if the realm user usage is assigned with the DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"username_pattern": schema.StringAttribute{
				Description:         "Allowed username pattern assigned the DEP enrollment.",
				MarkdownDescription: "Allowed username pattern assigned the DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"realm_user_is_admin": schema.BoolAttribute{
				Description:         "Toggles if the created user is admin with the DEP enrollment.",
				MarkdownDescription: "Toggles if the created user is admin with the DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
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
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"admin_password_complexity": schema.Int64Attribute{
				Description:         "Required password complexity for the DEP enrollment.",
				MarkdownDescription: "Required password complexity for the DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"admin_password_rotation_delay": schema.Int64Attribute{
				Description:         "Delay for the password rotation for the DEP enrollment.",
				MarkdownDescription: "Delay for the password rotation for the DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
			"allow_pairing": schema.BoolAttribute{
				Description:         "Toggles if the pairing is allowed with the DEP enrollment.",
				MarkdownDescription: "Toggles if the pairing is allowed with the DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"auto_advance_setup": schema.BoolAttribute{
				Description:         "Toggles if the auto advance setup is used with the DEP enrollment.",
				MarkdownDescription: "Toggles if the auto advance setup is used with the DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"await_device_configured": schema.BoolAttribute{
				Description:         "Toggles if await device configured is used with the DEP enrollment.",
				MarkdownDescription: "Toggles if await device configured is used with the DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"department": schema.StringAttribute{
				Description:         "Department asigned with the DEP enrollment.",
				MarkdownDescription: "Department asigned with the DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"is_mandatory": schema.BoolAttribute{
				Description:         "Toggles if the DEP enrollment is mandatory.",
				MarkdownDescription: "Toggles if the DEP enrollment is mandatory.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_mdm_removable": schema.BoolAttribute{
				Description:         "Toggles if the DEP enrollment is removable via mdm.",
				MarkdownDescription: "Toggles if the DEP enrollment is removable via mdm.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"is_multi_user": schema.BoolAttribute{
				Description:         "Toggles if the DEP enrollment allows multi users.",
				MarkdownDescription: "Toggles if the DEP enrollment allows multi users.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_supervised": schema.BoolAttribute{
				Description:         "Toggles if the DEP enrollment is supervised.",
				MarkdownDescription: "Toggles if the DEP enrollment is supervised.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"language": schema.StringAttribute{
				Description:         "Language asigned with the DEP enrollment.",
				MarkdownDescription: "Language asigned with the DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"org_magic": schema.StringAttribute{
				Description:         "Org magic asigned with the DEP enrollment.",
				MarkdownDescription: "Org magic asigned with the DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"region": schema.StringAttribute{
				Description:         "Region asigned with the DEP enrollment.",
				MarkdownDescription: "Region asigned with the DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"skip_setup_items": schema.SetAttribute{
				Description:         "Set of to be skiped setups assigned with the DEP enrollment.",
				MarkdownDescription: "Set of to be skiped setups assigned with the DEP enrollment.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"support_email_address": schema.StringAttribute{
				Description:         "Support email address asigned with the DEP enrollment.",
				MarkdownDescription: "Support email address asigned with the DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"support_phone_number": schema.StringAttribute{
				Description:         "Support phone number asigned with the DEP enrollment.",
				MarkdownDescription: "Support phone number asigned with the DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"include_tls_certificates": schema.BoolAttribute{
				Description:         "Toggles if the DEP enrollment include tls certificates.",
				MarkdownDescription: "Toggles if the DEP enrollment include tls certificates.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ios_max_version": schema.StringAttribute{
				Description:         "Max iOS version enforced with this DEP enrollment.",
				MarkdownDescription: "Max iOS version enforced with this DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ios_min_version": schema.StringAttribute{
				Description:         "Min iOS version enforced with this DEP enrollment.",
				MarkdownDescription: "Min iOS version enforced with this DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"macos_max_version": schema.StringAttribute{
				Description:         "Max macOS version enforced with this DEP enrollment.",
				MarkdownDescription: "Max macOS version enforced with this DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"macos_min_version": schema.StringAttribute{
				Description:         "Min macOS version enforced with this DEP enrollment.",
				MarkdownDescription: "Min macOS version enforced with this DEP enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"blueprint_id": schema.Int64Attribute{
				Description:         "ID of the MDM blueprint linked to the DEP enrollment.",
				MarkdownDescription: "`ID` of the MDM blueprint linked to the DEP enrollment.",
				Optional:            true,
			},
			"push_certificate_id": schema.Int64Attribute{
				Description:         "ID of the MDM push certificate linked to the DEP enrollment.",
				MarkdownDescription: "`ID` of the MDM push certificate linked to the DEP enrollment.",
				Required:            true,
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
				Required:            true,
			},
			"virtual_server_id": schema.Int64Attribute{
				Description:         "ID of the MDM virtual server linked to the DEP enrollment.",
				MarkdownDescription: "`ID` of the MDM virtual server linked to the DEP enrollment.",
				Required:            true,
			},
			"secret": schema.StringAttribute{
				Description:         "Enrollment secret.",
				MarkdownDescription: "Enrollment secret.",
				Computed:            true,
			},
			"meta_business_unit_id": schema.Int64Attribute{
				Description:         "The ID of the meta business unit the machine will be assigned to at enrollment.",
				MarkdownDescription: "The `ID` of the meta business unit the machine will be assigned to at enrollment.",
				Required:            true,
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags that the machine will get at enrollment.",
				MarkdownDescription: "The `ID`s of the tags that the machine will get at enrollment.",
				ElementType:         types.Int64Type,
				Required:            true,
			},
			"serial_numbers": schema.SetAttribute{
				Description:         "The serial numbers the enrollment is restricted to.",
				MarkdownDescription: "The serial numbers the enrollment is restricted to.",
				ElementType:         types.StringType,
				Required:            true,
			},
			"udids": schema.SetAttribute{
				Description:         "The UDIDs the enrollment is restricted to.",
				MarkdownDescription: "The `UDID`s the enrollment is restricted to.",
				ElementType:         types.StringType,
				Required:            true,
			},
			"quota": schema.Int64Attribute{
				Description:         "The number of time the enrollment can be used.",
				MarkdownDescription: "The number of time the enrollment can be used.",
				Optional:            true,
			},
		},
	}
}

func (r *MDMDEPEnrollmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*goztl.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *goztl.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *MDMDEPEnrollmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmDEPEnrollment

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlEnrollment, _, err := r.client.MDMDEPEnrollments.Create(ctx, mdmDEPEnrollmentCreateRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM DEP enrollment, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a MDM DEP enrollment")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmDEPEnrollmentForState(ztlEnrollment))...)
}

func (r *MDMDEPEnrollmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmDEPEnrollment

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlEnrollment, _, err := r.client.MDMDEPEnrollments.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM DEP enrollment %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a MDM DEP enrollment")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmDEPEnrollmentForState(ztlEnrollment))...)
}

func (r *MDMDEPEnrollmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmDEPEnrollment

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlEnrollment, _, err := r.client.MDMDEPEnrollments.Update(ctx, int(data.ID.ValueInt64()), mdmDEPEnrollmentUpdateRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM DEP enrollment %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a MDM DEP enrollment")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmDEPEnrollmentForState(ztlEnrollment))...)
}

func (r *MDMDEPEnrollmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmDEPEnrollment

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMDEPEnrollments.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM DEP enrollment %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a MDM DEP enrollment")
}

func (r *MDMDEPEnrollmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "MDM DEP enrollment", req, resp)
}
