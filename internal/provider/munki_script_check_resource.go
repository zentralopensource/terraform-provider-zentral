package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MunkiScriptCheckResource{}
var _ resource.ResourceWithImportState = &MunkiScriptCheckResource{}

func NewMunkiScriptCheckResource() resource.Resource {
	return &MunkiScriptCheckResource{}
}

// MunkiScriptCheckResource defines the resource implementation.
type MunkiScriptCheckResource struct {
	client *goztl.Client
}

func (r *MunkiScriptCheckResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_munki_script_check"
}

func (r *MunkiScriptCheckResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Munki script checks.",
		MarkdownDescription: "The resource `zentral_munki_script_check` manages Munki script checks.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Munki script check.",
				MarkdownDescription: "`ID` of the Munki script check.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the Munki script check.",
				MarkdownDescription: "Name of the Munki script check.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the Munki script check.",
				MarkdownDescription: "Description of the Munki script check.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"type": schema.StringAttribute{
				Description:         "Type of the script check. Can be ZSH_STR, ZSH_INT or ZSH_BOOL. Defaults to ZSH_STR.",
				MarkdownDescription: "Type of the script check. Can be `ZSH_STR`, `ZSH_INT` or `ZSH_BOOL`. Defaults to `ZSH_STR`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("ZSH_STR"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"ZSH_STR", "ZSH_INT", "ZSH_BOOL"}...),
				},
			},
			"source": schema.StringAttribute{
				Description:         "Source of the Munki script check.",
				MarkdownDescription: "Source of the Munki script check.",
				Required:            true,
			},
			"expected_result": schema.StringAttribute{
				Description:         "Expected result of the Munki script check.",
				MarkdownDescription: "Expected result of the Munki script check.",
				Required:            true,
			},
			"arch_amd64": schema.BoolAttribute{
				Description:         "If true, this Munki script check will be scheduled on Intel machines. Defaults to true.",
				MarkdownDescription: "If `true`, this Munki script check will be scheduled on Intel machines. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"arch_arm64": schema.BoolAttribute{
				Description:         "If true, this Munki script check will be scheduled on Apple Silicon machines. Defaults to true.",
				MarkdownDescription: "If `true`, this Munki script check will be scheduled on Apple Silicon machines. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"min_os_version": schema.StringAttribute{
				Description:         "This Munki script check will be scheduled on machines with an OS version higher or equal to this value.",
				MarkdownDescription: "This Munki script check will be scheduled on machines with an OS version higher or equal to this value.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"max_os_version": schema.StringAttribute{
				Description:         "This Munki script check will be scheduled on machines with an OS version lower than this value.",
				MarkdownDescription: "This Munki script check will be scheduled on machines with an OS version lower than this value.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags this Munki script check is restricted to.",
				MarkdownDescription: "The IDs of the tags this Munki script check is restricted to.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				Description:         "Version of the Munki script check.",
				MarkdownDescription: "Version of the Munki script check.",
				Computed:            true,
			},
		},
	}
}

func (r *MunkiScriptCheckResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MunkiScriptCheckResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data munkiScriptCheck

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSC, _, err := r.client.MunkiScriptChecks.Create(ctx, munkiScriptCheckRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Munki script check, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Munki script check")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, munkiScriptCheckForState(ztlMSC))...)
}

func (r *MunkiScriptCheckResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data munkiScriptCheck

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSC, _, err := r.client.MunkiScriptChecks.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Munki script check %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Munki script check")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, munkiScriptCheckForState(ztlMSC))...)
}

func (r *MunkiScriptCheckResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data munkiScriptCheck

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSC, _, err := r.client.MunkiScriptChecks.Update(ctx, int(data.ID.ValueInt64()), munkiScriptCheckRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Munki script check %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Munki script check")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, munkiScriptCheckForState(ztlMSC))...)
}

func (r *MunkiScriptCheckResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data munkiScriptCheck

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MunkiScriptChecks.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Munki script check %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Munki script check")
}

func (r *MunkiScriptCheckResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Munki script check", req, resp)
}
