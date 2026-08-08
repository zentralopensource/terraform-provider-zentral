package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &TurboOneTimeJobResource{}
var _ resource.ResourceWithImportState = &TurboOneTimeJobResource{}

func NewTurboOneTimeJobResource() resource.Resource {
	return &TurboOneTimeJobResource{}
}

// TurboOneTimeJobResource defines the resource implementation.
type TurboOneTimeJobResource struct {
	client *goztl.Client
}

func (r *TurboOneTimeJobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_turbo_one_time_job"
}

func (r *TurboOneTimeJobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Turbo one-time jobs.",
		MarkdownDescription: "The resource `zentral_turbo_one_time_job` schedules a Turbo job (script or mSCP check) to run once in a configuration.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the Turbo one-time job.",
				MarkdownDescription: "`ID` of the Turbo one-time job.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"configuration_id": schema.StringAttribute{
				Description:         "ID of the Turbo configuration the job is scheduled in.",
				MarkdownDescription: "`ID` of the Turbo configuration the job is scheduled in.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"job_id": schema.StringAttribute{
				Description:         "ID of the job to schedule. Use the job_id attribute of a zentral_turbo_script or zentral_turbo_mscp_check.",
				MarkdownDescription: "`ID` of the job to schedule. Use the `job_id` attribute of a `zentral_turbo_script` or `zentral_turbo_mscp_check`.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"not_before": schema.StringAttribute{
				Description:         "The job is not delivered before this time, as a UTC timestamp without a timezone suffix (e.g. 2026-09-01T09:00:00). Leave unset to deliver immediately.",
				MarkdownDescription: "The job is not delivered before this time, as a UTC timestamp without a timezone suffix (e.g. `2026-09-01T09:00:00`). Leave unset to deliver immediately.",
				Optional:            true,
				Validators: []validator.String{
					naiveTimestampValidator,
				},
			},
			"not_after": schema.StringAttribute{
				Description:         "The job is not delivered after this time, as a UTC timestamp without a timezone suffix (e.g. 2026-09-30T09:00:00). Leave unset for no expiry.",
				MarkdownDescription: "The job is not delivered after this time, as a UTC timestamp without a timezone suffix (e.g. `2026-09-30T09:00:00`). Leave unset for no expiry.",
				Optional:            true,
				Validators: []validator.String{
					naiveTimestampValidator,
				},
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags the job is scoped to. Empty scopes the job to all machines in the configuration.",
				MarkdownDescription: "The `ID`s of the tags the job is scoped to. Empty scopes the job to all machines in the configuration.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"excluded_tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags the job is excluded from.",
				MarkdownDescription: "The `ID`s of the tags the job is excluded from.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"serial_numbers": schema.SetAttribute{
				Description:         "The serial numbers the job is scoped to.",
				MarkdownDescription: "The serial numbers the job is scoped to.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"excluded_serial_numbers": schema.SetAttribute{
				Description:         "The serial numbers the job is excluded from.",
				MarkdownDescription: "The serial numbers the job is excluded from.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *TurboOneTimeJobResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TurboOneTimeJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data turboOneTimeJob

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlTOTJ, _, err := r.client.TurboOneTimeJobs.Create(ctx, turboOneTimeJobRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Turbo one-time job, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Turbo one-time job")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, turboOneTimeJobForState(ztlTOTJ))...)
}

func (r *TurboOneTimeJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data turboOneTimeJob

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlTOTJ, _, err := r.client.TurboOneTimeJobs.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Turbo one-time job %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Turbo one-time job")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, turboOneTimeJobForState(ztlTOTJ))...)
}

func (r *TurboOneTimeJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data turboOneTimeJob

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlTOTJ, _, err := r.client.TurboOneTimeJobs.Update(ctx, data.ID.ValueString(), turboOneTimeJobRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Turbo one-time job %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Turbo one-time job")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, turboOneTimeJobForState(ztlTOTJ))...)
}

func (r *TurboOneTimeJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data turboOneTimeJob

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.TurboOneTimeJobs.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Turbo one-time job %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Turbo one-time job")
}

func (r *TurboOneTimeJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "Turbo one-time job", req, resp)
}
