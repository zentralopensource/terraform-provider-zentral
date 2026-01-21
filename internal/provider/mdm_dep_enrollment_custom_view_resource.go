package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMDEPEnrollmentResource{}
var _ resource.ResourceWithImportState = &MDMDEPEnrollmentResource{}

func NewMDMDEPEnrollmentCutomViewResource() resource.Resource {
	return &MDMDEPEnrollmentCutomViewResource{}
}

// MDMEnrollmentCutomViewResource defines the resource implementation.
type MDMDEPEnrollmentCutomViewResource struct {
	client *goztl.Client
}

func (r *MDMDEPEnrollmentCutomViewResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_dep_enrollment_custom_view"
}

func (r *MDMDEPEnrollmentCutomViewResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM DEP enrollments custom views.",
		MarkdownDescription: "The resource `zentral_mdm_dep_enrollment_custom_view` manages MDM DEP enrollments custom views.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the MDM DEP enrollment custom view.",
				MarkdownDescription: "ID of the MDM DEP enrollment custom view.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dep_enrollment": schema.Int64Attribute{
				Description:         "ID of MDM DEP enrollment of the MDM dep enrollment custom view.",
				MarkdownDescription: "ID of MDM DEP enrollment of the MDM dep enrollment custom view.",
				Required:            true,
			},
			"custom_view": schema.StringAttribute{
				Description:         "ID of MDM custom view of the MDM dep enrollment custom view.",
				MarkdownDescription: "ID of MDM custom view of the MDM dep enrollment custom view.",
				Required:            true,
			},
			"weight": schema.Int64Attribute{
				Description:         "Weight of the MDM DEP enrollment custom view.",
				MarkdownDescription: "Weight of the MDM DEP enrollment custom view.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *MDMDEPEnrollmentCutomViewResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMDEPEnrollmentCutomViewResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmDEPEnrollmentCutomView

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlCustomView, _, err := r.client.MDMDEPEnrollmentCustomViews.Create(ctx, mdmDEPEnrollmentCustomViewRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM DEP enrollment custom view, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a MDM DEPenrollment custom view")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmDEPEnrollmentCustomViewForState(ztlCustomView))...)
}

func (r *MDMDEPEnrollmentCutomViewResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmDEPEnrollmentCutomView

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlCutomView, _, err := r.client.MDMDEPEnrollmentCustomViews.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM DEPenrollment custom view %s, got error: %s", data.ID, err),
		)
		return
	}

	tflog.Trace(ctx, "read a MDM DEPenrollment custom view")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmDEPEnrollmentCustomViewForState(ztlCutomView))...)
}

func (r *MDMDEPEnrollmentCutomViewResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmDEPEnrollmentCutomView
	var state mdmDEPEnrollmentCutomView

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	// Read prior state to get computed ID reliably
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlCustomView, _, err := r.client.MDMDEPEnrollmentCustomViews.Update(ctx, state.ID.ValueString(), mdmDEPEnrollmentCustomViewRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM dep enrollment custom view %s, got error: %s", state.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a MDM DEP enrollment custom view")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmDEPEnrollmentCustomViewForState(ztlCustomView))...)
}

func (r *MDMDEPEnrollmentCutomViewResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmDEPEnrollmentCutomView

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMDEPEnrollmentCustomViews.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM DEP enrollment custom view %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a MDM DEP enrollment custom view")
}

func (r *MDMDEPEnrollmentCutomViewResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "MDM DEP enrollment custom view", req, resp)
}
