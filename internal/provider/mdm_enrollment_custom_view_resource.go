package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMDEPEnrollmentResource{}
var _ resource.ResourceWithImportState = &MDMDEPEnrollmentResource{}

func NewMDMEnrollmentCustomViewResource() resource.Resource {
	return &MDMEnrollmentCustomViewResource{}
}

// MDMEnrollmentCustomViewResource defines the resource implementation.
type MDMEnrollmentCustomViewResource struct {
	client *goztl.Client
}

func (r *MDMEnrollmentCustomViewResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_enrollment_custom_view"
}

func (r *MDMEnrollmentCustomViewResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM enrollments custom views.",
		MarkdownDescription: "The resource `zentral_mdm_enrollment_custom_view` manages MDM enrollments custom views.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the MDM enrollment custom view.",
				MarkdownDescription: "ID of the MDM enrollment custom view.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the MDM enrollment custom view.",
				MarkdownDescription: "Name of the MDM enrollment custom view.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the MDM enrollment custom view.",
				MarkdownDescription: "Description of the MDM enrollment custom view.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"html": schema.StringAttribute{
				Description:         "HTML template.",
				MarkdownDescription: "HTML template.",
				Required:            true,
			},
			"requires_authentication": schema.BoolAttribute{
				Description:         "If true, the custom view will be displayed after the realm authentication. Default is true.",
				MarkdownDescription: "If `true`, the custom view will be displayed after the realm authentication. Default is `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func (r *MDMEnrollmentCustomViewResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMEnrollmentCustomViewResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmEnrollmentCustomView

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlCustomView, _, err := r.client.MDMEnrollmentCustomViews.Create(ctx, mdmEnrollmentCustomViewRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM enrollment custom view, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a MDM enrollment custom view")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmEnrollmentCustomViewForState(ztlCustomView))...)
}

func (r *MDMEnrollmentCustomViewResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmEnrollmentCustomView

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlCustomView, _, err := r.client.MDMEnrollmentCustomViews.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM enrollment custom view %s, got error: %s", data.ID, err),
		)
		return
	}

	tflog.Trace(ctx, "read a MDM enrollment custom view")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmEnrollmentCustomViewForState(ztlCustomView))...)
}

func (r *MDMEnrollmentCustomViewResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmEnrollmentCustomView

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlCustomView, _, err := r.client.MDMEnrollmentCustomViews.Update(ctx, data.ID.ValueString(), mdmEnrollmentCustomViewRequestWithState(data))

	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM enrollment custom view %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a MDM enrollment custom view")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmEnrollmentCustomViewForState(ztlCustomView))...)
}

func (r *MDMEnrollmentCustomViewResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmEnrollmentCustomView

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMEnrollmentCustomViews.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM enrollment custom view %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a MDM enrollment custom view")
}

func (r *MDMEnrollmentCustomViewResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "MDM enrollment custom view", req, resp)
}
