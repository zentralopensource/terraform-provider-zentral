package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMRecoveryPasswordConfigResource{}
var _ resource.ResourceWithImportState = &MDMRecoveryPasswordConfigResource{}

func NewMDMRecoveryPasswordConfigResource() resource.Resource {
	return &MDMRecoveryPasswordConfigResource{}
}

// MDMRecoveryPasswordConfigResource defines the resource implementation.
type MDMRecoveryPasswordConfigResource struct {
	client *goztl.Client
}

func (r *MDMRecoveryPasswordConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_recovery_password_config"
}

func (r *MDMRecoveryPasswordConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM recovery password configurations.",
		MarkdownDescription: "The resource `zentral_mdm_recovery_password_config` manages MDM recovery password configurations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the recovery password configuration.",
				MarkdownDescription: "`ID` of the recovery password configuration.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the recovery password configuration.",
				MarkdownDescription: "Name of the recovery password configuration.",
				Required:            true,
			},
			"dynamic_password": schema.BoolAttribute{
				Description:         "If true, a unique password is generated for each device. Defaults to true.",
				MarkdownDescription: "If `true`, a unique password is generated for each device. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"static_password": schema.StringAttribute{
				Description:         "The static password to set for all devices.",
				MarkdownDescription: "The  static password to set for all devices.",
				Optional:            true,
			},
			"rotation_interval_days": schema.Int64Attribute{
				Description:         "The automatic recovery password rotation interval in days. It has a maximum value of 365. Defaults to 0 (no automatic rotation).",
				MarkdownDescription: "The automatic recovery password rotation interval in days. It has a maximum value of `365`. Defaults to `0` (no automatic rotation).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 365),
				},
			},
			"rotate_firmware_password": schema.BoolAttribute{
				Description:         "Set to true to rotate the firmware passwords. Defaults to false.",
				MarkdownDescription: "Set to `true` to rotate the firmware passwords. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func (r *MDMRecoveryPasswordConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMRecoveryPasswordConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmRecoveryPasswordConfig

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMRPC, _, err := r.client.MDMRecoveryPasswordConfigs.Create(ctx, mdmRecoveryPasswordConfigRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM recovery password configuration, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a MDM recovery password configuration")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmRecoveryPasswordConfigForState(ztlMRPC))...)
}

func (r *MDMRecoveryPasswordConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmRecoveryPasswordConfig

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMRPC, _, err := r.client.MDMRecoveryPasswordConfigs.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM recovery password configuration %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a MDM recovery password configuration")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmRecoveryPasswordConfigForState(ztlMRPC))...)
}

func (r *MDMRecoveryPasswordConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmRecoveryPasswordConfig

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMRPC, _, err := r.client.MDMRecoveryPasswordConfigs.Update(ctx, int(data.ID.ValueInt64()), mdmRecoveryPasswordConfigRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM recovery password configuration %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a MDM recovery password configuration")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmRecoveryPasswordConfigForState(ztlMRPC))...)
}

func (r *MDMRecoveryPasswordConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmRecoveryPasswordConfig

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMRecoveryPasswordConfigs.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM recovery password configuration %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a MDM recovery password configuration")
}

func (r *MDMRecoveryPasswordConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "MDM recovery password configuration", req, resp)
}
