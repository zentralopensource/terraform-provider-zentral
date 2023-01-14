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
var _ resource.Resource = &SantaEnrollmentResource{}
var _ resource.ResourceWithImportState = &SantaEnrollmentResource{}

func NewSantaEnrollmentResource() resource.Resource {
	return &SantaEnrollmentResource{}
}

// SantaEnrollmentResource defines the resource implementation.
type SantaEnrollmentResource struct {
	client *goztl.Client
}

func (r *SantaEnrollmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_santa_enrollment"
}

func (r *SantaEnrollmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Santa enrollments.",
		MarkdownDescription: "The resource `zentral_santa_enrollment` manages Santa enrollments.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Santa enrollment.",
				MarkdownDescription: "`ID` of the Santa enrollment.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"configuration_id": schema.Int64Attribute{
				Description:         "ID of the Santa configuration.",
				MarkdownDescription: "`ID` of the Santa configuration.",
				Required:            true,
			},
			"configuration_profile_url": schema.StringAttribute{
				Description:         "Configuration profile download URL.",
				MarkdownDescription: "Configuration profile download URL.",
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

func (r *SantaEnrollmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SantaEnrollmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data santaEnrollment

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlSE, _, err := r.client.SantaEnrollments.Create(ctx, santaEnrollmentRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Santa enrollment, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Santa enrollment")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, santaEnrollmentForState(ztlSE))...)
}

func (r *SantaEnrollmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data santaEnrollment

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlSE, _, err := r.client.SantaEnrollments.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Santa enrollment %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Santa enrollment")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, santaEnrollmentForState(ztlSE))...)
}

func (r *SantaEnrollmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data santaEnrollment

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlSE, _, err := r.client.SantaEnrollments.Update(ctx, int(data.ID.ValueInt64()), santaEnrollmentRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Santa enrollment %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Santa enrollment")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, santaEnrollmentForState(ztlSE))...)
}

func (r *SantaEnrollmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data santaEnrollment

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.SantaEnrollments.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Santa enrollment %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Santa enrollment")
}

func (r *SantaEnrollmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Santa enrollment", req, resp)
}
