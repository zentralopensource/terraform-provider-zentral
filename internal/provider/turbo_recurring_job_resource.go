package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &TurboRecurringJobResource{}
var _ resource.ResourceWithImportState = &TurboRecurringJobResource{}

func NewTurboRecurringJobResource() resource.Resource {
	return &TurboRecurringJobResource{}
}

// TurboRecurringJobResource defines the resource implementation.
type TurboRecurringJobResource struct {
	client *goztl.Client
}

func (r *TurboRecurringJobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_turbo_recurring_job"
}

func (r *TurboRecurringJobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Turbo recurring jobs.",
		MarkdownDescription: "The resource `zentral_turbo_recurring_job` schedules a Turbo job (script or mSCP check) to run repeatedly in a configuration.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the Turbo recurring job.",
				MarkdownDescription: "`ID` of the Turbo recurring job.",
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
			"interval": schema.Int64Attribute{
				Description:         "Run interval in seconds. Must be between 60 and 604800. Leave unset to use the configuration default.",
				MarkdownDescription: "Run interval in seconds. Must be between `60` and `604800`. Leave unset to use the configuration default.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(60, 604800),
				},
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags the job is scoped to. Empty scopes the job to all machines in the configuration.",
				MarkdownDescription: "The `ID`s of the tags the job is scoped to. Empty scopes the job to all machines in the configuration.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.Int64Type, []attr.Value{})),
			},
			"excluded_tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags the job is excluded from.",
				MarkdownDescription: "The `ID`s of the tags the job is excluded from.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.Int64Type, []attr.Value{})),
			},
			"serial_numbers": schema.SetAttribute{
				Description:         "The serial numbers the job is scoped to.",
				MarkdownDescription: "The serial numbers the job is scoped to.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"excluded_serial_numbers": schema.SetAttribute{
				Description:         "The serial numbers the job is excluded from.",
				MarkdownDescription: "The serial numbers the job is excluded from.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
		},
	}
}

func (r *TurboRecurringJobResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TurboRecurringJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data turboRecurringJob

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlTRJ, _, err := r.client.TurboRecurringJobs.Create(ctx, turboRecurringJobRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Turbo recurring job, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Turbo recurring job")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, turboRecurringJobForState(ztlTRJ))...)
}

func (r *TurboRecurringJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data turboRecurringJob

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlTRJ, _, err := r.client.TurboRecurringJobs.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Turbo recurring job %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Turbo recurring job")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, turboRecurringJobForState(ztlTRJ))...)
}

func (r *TurboRecurringJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data turboRecurringJob

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlTRJ, _, err := r.client.TurboRecurringJobs.Update(ctx, data.ID.ValueString(), turboRecurringJobRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Turbo recurring job %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Turbo recurring job")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, turboRecurringJobForState(ztlTRJ))...)
}

func (r *TurboRecurringJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data turboRecurringJob

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.TurboRecurringJobs.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Turbo recurring job %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Turbo recurring job")
}

func (r *TurboRecurringJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "Turbo recurring job", req, resp)
}
