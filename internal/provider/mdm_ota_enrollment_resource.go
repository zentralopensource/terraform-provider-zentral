package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMOTAEnrollmentResource{}
var _ resource.ResourceWithImportState = &MDMOTAEnrollmentResource{}

func NewMDMOTAEnrollmentResource() resource.Resource {
	return &MDMOTAEnrollmentResource{}
}

// MDMOTAEnrollmentResource defines the resource implementation.
type MDMOTAEnrollmentResource struct {
	client *goztl.Client
}

func (r *MDMOTAEnrollmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_ota_enrollment"
}

func (r *MDMOTAEnrollmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM OTA enrollments.",
		MarkdownDescription: "The resource `zentral_mdm_ota_enrollment` manages MDM OTA enrollments.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the MDM OTA enrollment.",
				MarkdownDescription: "`ID` of the MDM OTA enrollment.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the MDM OTA enrollment.",
				MarkdownDescription: "Name of the MDM OTA enrollment.",
				Required:            true,
			},
			"display_name": schema.StringAttribute{
				Description:         "Name of the MDM OTA enrollment as displayed on the device.",
				MarkdownDescription: "Name of the MDM OTA enrollment as displayed on the device.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Zentral MDM"),
			},
			"blueprint_id": schema.Int64Attribute{
				Description:         "ID of the MDM blueprint linked to the OTA enrollment.",
				MarkdownDescription: "`ID` of the MDM blueprint linked to the OTA enrollment.",
				Optional:            true,
				// not Computed, because it triggers a foreign key error on the server with PK = 0
			},
			"push_certificate_id": schema.Int64Attribute{
				Description:         "ID of the MDM push certificate linked to the OTA enrollment.",
				MarkdownDescription: "`ID` of the MDM push certificate linked to the OTA enrollment.",
				Required:            true,
			},
			"realm_uuid": schema.StringAttribute{
				Description:         "UUID of the identity realm linked to the OTA enrollment.",
				MarkdownDescription: "`UUID` of the identity realm linked to the OTA enrollment.",
				Optional:            true,
				// not Computed, because it triggers a foreign key error on the server with PK = 0
			},
			"scep_config_id": schema.Int64Attribute{
				Description:         "ID of the MDM SCEP configuration linked to the OTA enrollment.",
				MarkdownDescription: "`ID` of the MDM SCEP configuration linked to the OTA enrollment.",
				Required:            true,
			},
			"scep_verification": schema.BoolAttribute{
				Description:         "Indicates if a SCEP verification is expected during the enrollment.",
				MarkdownDescription: "Indicates if a SCEP verification is expected during the enrollment.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
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
				Optional:            true,
				Computed:            true,
			},
			"serial_numbers": schema.SetAttribute{
				Description:         "The serial numbers the enrollment is restricted to.",
				MarkdownDescription: "The serial numbers the enrollment is restricted to.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"udids": schema.SetAttribute{
				Description:         "The UDIDs the enrollment is restricted to.",
				MarkdownDescription: "The `UDID`s the enrollment is restricted to.",
				ElementType:         types.StringType,
				Optional:            true,
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

func (r *MDMOTAEnrollmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMOTAEnrollmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmOTAEnrollment

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMOE, _, err := r.client.MDMOTAEnrollments.Create(ctx, mdmOTAEnrollmentRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM OTA enrollment, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a MDM OTA enrollment")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmOTAEnrollmentForState(ztlMOE))...)
}

func (r *MDMOTAEnrollmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmOTAEnrollment

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMOE, _, err := r.client.MDMOTAEnrollments.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM OTA enrollment %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a MDM OTA enrollment")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmOTAEnrollmentForState(ztlMOE))...)
}

func (r *MDMOTAEnrollmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmOTAEnrollment

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMOE, _, err := r.client.MDMOTAEnrollments.Update(ctx, int(data.ID.ValueInt64()), mdmOTAEnrollmentRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM OTA enrollment %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a MDM OTA enrollment")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmOTAEnrollmentForState(ztlMOE))...)
}

func (r *MDMOTAEnrollmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmOTAEnrollment

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMOTAEnrollments.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM OTA enrollment %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a MDM OTA enrollment")
}

func (r *MDMOTAEnrollmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "MDM OTA enrollment", req, resp)
}
