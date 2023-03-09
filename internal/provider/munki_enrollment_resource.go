package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MunkiEnrollmentResource{}
var _ resource.ResourceWithImportState = &MunkiEnrollmentResource{}

func NewMunkiEnrollmentResource() resource.Resource {
	return &MunkiEnrollmentResource{}
}

// MunkiEnrollmentResource defines the resource implementation.
type MunkiEnrollmentResource struct {
	client *goztl.Client
}

func (r *MunkiEnrollmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_munki_enrollment"
}

func (r *MunkiEnrollmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Munki enrollments.",
		MarkdownDescription: "The resource `zentral_munki_enrollment` manages Munki enrollments.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Munki enrollment.",
				MarkdownDescription: "`ID` of the Munki enrollment.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"configuration_id": schema.Int64Attribute{
				Description:         "ID of the Munki configuration.",
				MarkdownDescription: "`ID` of the Munki configuration.",
				Required:            true,
			},
			"package_url": schema.StringAttribute{
				Description:         "Configuration package download URL.",
				MarkdownDescription: "Configuration package download URL.",
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

func (r *MunkiEnrollmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MunkiEnrollmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data munkiEnrollment

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlME, _, err := r.client.MunkiEnrollments.Create(ctx, munkiEnrollmentRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Munki enrollment, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Munki enrollment")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, munkiEnrollmentForState(ztlME))...)
}

func (r *MunkiEnrollmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data munkiEnrollment

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlME, _, err := r.client.MunkiEnrollments.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Munki enrollment %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Munki enrollment")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, munkiEnrollmentForState(ztlME))...)
}

func (r *MunkiEnrollmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data munkiEnrollment

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlME, _, err := r.client.MunkiEnrollments.Update(ctx, int(data.ID.ValueInt64()), munkiEnrollmentRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Munki enrollment %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Munki enrollment")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, munkiEnrollmentForState(ztlME))...)
}

func (r *MunkiEnrollmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data munkiEnrollment

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MunkiEnrollments.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Munki enrollment %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Munki enrollment")
}

func (r *MunkiEnrollmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Munki enrollment", req, resp)
}
