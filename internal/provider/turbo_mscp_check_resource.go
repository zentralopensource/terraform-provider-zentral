package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &TurboMSCPCheckResource{}
var _ resource.ResourceWithImportState = &TurboMSCPCheckResource{}
var _ resource.ResourceWithValidateConfig = &TurboMSCPCheckResource{}

func NewTurboMSCPCheckResource() resource.Resource {
	return &TurboMSCPCheckResource{}
}

// TurboMSCPCheckResource defines the resource implementation.
type TurboMSCPCheckResource struct {
	client *goztl.Client
}

func (r *TurboMSCPCheckResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_turbo_mscp_check"
}

func (r *TurboMSCPCheckResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Turbo mSCP checks.",
		MarkdownDescription: "The resource `zentral_turbo_mscp_check` manages Turbo mSCP checks.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the Turbo mSCP check.",
				MarkdownDescription: "`ID` of the Turbo mSCP check.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rule_id": schema.StringAttribute{
				Description:         "mSCP rule ID of the check.",
				MarkdownDescription: "mSCP rule ID of the check.",
				Required:            true,
			},
			"baseline": schema.StringAttribute{
				Description:         "mSCP baseline key (e.g. cis_lvl1, stig); the agent uses that baseline's default ODV for the rule. Mutually exclusive with an explicit ODV override.",
				MarkdownDescription: "mSCP baseline key (e.g. `cis_lvl1`, `stig`); the agent uses that baseline's default ODV for the rule. Mutually exclusive with an explicit ODV override.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"odv_int": schema.Int64Attribute{
				Description:         "Integer Organization Defined Value override for the rule. Mutually exclusive with the other ODV overrides and with baseline.",
				MarkdownDescription: "Integer Organization Defined Value override for the rule. Mutually exclusive with the other ODV overrides and with `baseline`.",
				Optional:            true,
			},
			"odv_string": schema.StringAttribute{
				Description:         "String Organization Defined Value override for the rule. Mutually exclusive with the other ODV overrides and with baseline.",
				MarkdownDescription: "String Organization Defined Value override for the rule. Mutually exclusive with the other ODV overrides and with `baseline`.",
				Optional:            true,
				Validators: []validator.String{
					// an empty override is no override for the API, which stores it as null and hands
					// back a value the configuration does not match
					stringvalidator.LengthAtLeast(1),
				},
			},
			"odv_bool": schema.BoolAttribute{
				Description:         "Boolean Organization Defined Value override for the rule. Mutually exclusive with the other ODV overrides and with baseline.",
				MarkdownDescription: "Boolean Organization Defined Value override for the rule. Mutually exclusive with the other ODV overrides and with `baseline`.",
				Optional:            true,
			},
			"version": schema.Int64Attribute{
				Description:         "Version of the Turbo mSCP check.",
				MarkdownDescription: "Version of the Turbo mSCP check.",
				Computed:            true,
			},
			"compliance_check_id": schema.Int64Attribute{
				Description:         "ID of the compliance check backing this mSCP check.",
				MarkdownDescription: "`ID` of the compliance check backing this mSCP check.",
				Computed:            true,
			},
			"job_id": schema.StringAttribute{
				Description:         "ID of the job this mSCP check is scheduled with. Use it to reference the check from a recurring or one-time job.",
				MarkdownDescription: "`ID` of the job this mSCP check is scheduled with. Use it to reference the check from a recurring or one-time job.",
				Computed:            true,
				// the job outlives the check updates (they only bump its version). Without this, an update
				// makes job_id unknown, and the jobs referencing it are planned for replacement.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *TurboMSCPCheckResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data turboMSCPCheck

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	setODVs := 0
	if !data.ODVInt.IsNull() {
		setODVs++
	}
	if !data.ODVString.IsNull() {
		setODVs++
	}
	if !data.ODVBool.IsNull() {
		setODVs++
	}

	if setODVs > 1 {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"Only one of odv_int, odv_string or odv_bool can be set.",
		)
		return
	}

	if setODVs > 0 && !data.Baseline.IsNull() && data.Baseline.ValueString() != "" {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"Set a baseline or an ODV override, not both.",
		)
	}
}

func (r *TurboMSCPCheckResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TurboMSCPCheckResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data turboMSCPCheck

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlTMC, _, err := r.client.TurboMSCPChecks.Create(ctx, turboMSCPCheckRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Turbo mSCP check, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Turbo mSCP check")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, turboMSCPCheckForState(ztlTMC))...)
}

func (r *TurboMSCPCheckResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data turboMSCPCheck

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlTMC, _, err := r.client.TurboMSCPChecks.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Turbo mSCP check %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Turbo mSCP check")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, turboMSCPCheckForState(ztlTMC))...)
}

func (r *TurboMSCPCheckResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data turboMSCPCheck

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlTMC, _, err := r.client.TurboMSCPChecks.Update(ctx, data.ID.ValueString(), turboMSCPCheckRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Turbo mSCP check %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Turbo mSCP check")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, turboMSCPCheckForState(ztlTMC))...)
}

func (r *TurboMSCPCheckResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data turboMSCPCheck

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.TurboMSCPChecks.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Turbo mSCP check %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Turbo mSCP check")
}

func (r *TurboMSCPCheckResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "Turbo mSCP check", req, resp)
}
