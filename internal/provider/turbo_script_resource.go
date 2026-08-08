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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &TurboScriptResource{}
var _ resource.ResourceWithImportState = &TurboScriptResource{}
var _ resource.ResourceWithValidateConfig = &TurboScriptResource{}

func NewTurboScriptResource() resource.Resource {
	return &TurboScriptResource{}
}

// TurboScriptResource defines the resource implementation.
type TurboScriptResource struct {
	client *goztl.Client
}

func (r *TurboScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_turbo_script"
}

func (r *TurboScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Turbo scripts.",
		MarkdownDescription: "The resource `zentral_turbo_script` manages Turbo scripts.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the Turbo script.",
				MarkdownDescription: "`ID` of the Turbo script.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the Turbo script.",
				MarkdownDescription: "Name of the Turbo script.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the Turbo script.",
				MarkdownDescription: "Description of the Turbo script.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"source": schema.StringAttribute{
				Description:         "Source of the Turbo script. A zsh script; exit 0 means OK, exit greater than 0 means FAIL.",
				MarkdownDescription: "Source of the Turbo script. A zsh script; exit `0` means OK, exit greater than `0` means FAIL.",
				Required:            true,
			},
			"tag_id": schema.Int64Attribute{
				Description:         "The ID of the tag added on exit 0 and removed on exit greater than 0.",
				MarkdownDescription: "The `ID` of the tag added on exit `0` and removed on exit greater than `0`.",
				Optional:            true,
			},
			"arch_amd64": schema.BoolAttribute{
				Description:         "If true, this Turbo script can run on Intel machines. Defaults to true.",
				MarkdownDescription: "If `true`, this Turbo script can run on Intel machines. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"arch_arm64": schema.BoolAttribute{
				Description:         "If true, this Turbo script can run on Apple Silicon machines. Defaults to true.",
				MarkdownDescription: "If `true`, this Turbo script can run on Apple Silicon machines. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"min_os_version": schema.StringAttribute{
				Description:         "This Turbo script will run on machines with an OS version higher or equal to this value.",
				MarkdownDescription: "This Turbo script will run on machines with an OS version higher or equal to this value.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"max_os_version": schema.StringAttribute{
				Description:         "This Turbo script will run on machines with an OS version lower than this value.",
				MarkdownDescription: "This Turbo script will run on machines with an OS version lower than this value.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"compliance_check_enabled": schema.BoolAttribute{
				Description:         "If true, the script is a compliance check. Defaults to false.",
				MarkdownDescription: "If `true`, the script is a compliance check. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"compliance_check_id": schema.Int64Attribute{
				Description:         "ID of the compliance check backing this script, when compliance_check_enabled is true.",
				MarkdownDescription: "`ID` of the compliance check backing this script, when `compliance_check_enabled` is `true`.",
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				Description:         "Version of the Turbo script.",
				MarkdownDescription: "Version of the Turbo script.",
				Computed:            true,
			},
			"job_id": schema.StringAttribute{
				Description:         "ID of the job this script is scheduled with. Use it to reference the script from a recurring or one-time job.",
				MarkdownDescription: "`ID` of the job this script is scheduled with. Use it to reference the script from a recurring or one-time job.",
				Computed:            true,
				// the job outlives the script updates (they only bump its version). Without this, an update
				// makes job_id unknown, and the jobs referencing it are planned for replacement.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *TurboScriptResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data turboScript

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// an unset attribute is null in the configuration: the true defaults are only applied to the plan
	disabled := func(b types.Bool) bool {
		return !b.IsNull() && !b.IsUnknown() && !b.ValueBool()
	}

	// a script that runs on neither architecture never runs anywhere, and the agent silently skips it
	if disabled(data.ArchAMD64) && disabled(data.ArchARM64) {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"At least one of arch_amd64 or arch_arm64 must be true.",
		)
	}
}

func (r *TurboScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TurboScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data turboScript

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlTS, _, err := r.client.TurboScripts.Create(ctx, turboScriptRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Turbo script, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Turbo script")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, turboScriptForState(ztlTS))...)
}

func (r *TurboScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data turboScript

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlTS, _, err := r.client.TurboScripts.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Turbo script %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Turbo script")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, turboScriptForState(ztlTS))...)
}

func (r *TurboScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data turboScript

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlTS, _, err := r.client.TurboScripts.Update(ctx, data.ID.ValueString(), turboScriptRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Turbo script %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Turbo script")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, turboScriptForState(ztlTS))...)
}

func (r *TurboScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data turboScript

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.TurboScripts.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Turbo script %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Turbo script")
}

func (r *TurboScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "Turbo script", req, resp)
}
