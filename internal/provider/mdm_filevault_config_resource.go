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
var _ resource.Resource = &MDMFileVaultConfigResource{}
var _ resource.ResourceWithImportState = &MDMFileVaultConfigResource{}

func NewMDMFileVaultConfigResource() resource.Resource {
	return &MDMFileVaultConfigResource{}
}

// MDMFileVaultConfigResource defines the resource implementation.
type MDMFileVaultConfigResource struct {
	client *goztl.Client
}

func (r *MDMFileVaultConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_filevault_config"
}

func (r *MDMFileVaultConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM FileVault configurations.",
		MarkdownDescription: "The resource `zentral_mdm_filevault_config` manages MDM FileVault configurations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the FileVault configuration.",
				MarkdownDescription: "`ID` of the FileVault configuration.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the FileVault configuration.",
				MarkdownDescription: "Name of the FileVault configuration.",
				Required:            true,
			},
			"escrow_location_display_name": schema.StringAttribute{
				Description:         "Description of the location where the FDE PRK will be escrowed. This text will be inserted into the message the user sees when enabling FileVault.",
				MarkdownDescription: "Description of the location where the FDE PRK will be escrowed. This text will be inserted into the message the user sees when enabling FileVault.",
				Required:            true,
			},
			"at_login_only": schema.BoolAttribute{
				Description:         "If true, prevents requests for enabling FileVault at user logout time. Defaults to false.",
				MarkdownDescription: "If `true`, prevents requests for enabling FileVault at user logout time. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"bypass_attempts": schema.Int64Attribute{
				Description:         "The maximum number of times users can bypass enabling FileVault before being required to enable it to log in.",
				MarkdownDescription: "The maximum number of times users can bypass enabling FileVault before being required to enable it to log in.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
				Validators: []validator.Int64{
					int64validator.Between(-1, 9999),
				},
			},
			"show_recovery_key": schema.BoolAttribute{
				Description:         "If false, prevents display of the personal recovery key to the user after FileVault is enabled. Defaults to false.",
				MarkdownDescription: "If `false`, prevents display of the personal recovery key to the user after FileVault is enabled. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"destroy_key_on_standby": schema.BoolAttribute{
				Description:         "Set to true to prevent storing the FileVault key across restarts. Defaults to false.",
				MarkdownDescription: "Set to `true` to prevent storing the FileVault key across restarts. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"prk_rotation_interval_days": schema.Int64Attribute{
				Description:         "The automatic PRK rotation interval in days. It has a maximum value of 365. Defaults to 0 (no automatic rotation).",
				MarkdownDescription: "The automatic PRK rotation interval in days. It has a maximum value of `365`. Defaults to `0` (no automatic rotation).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 365),
				},
			},
		},
	}
}

func (r *MDMFileVaultConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMFileVaultConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmFileVaultConfig

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMFC, _, err := r.client.MDMFileVaultConfigs.Create(ctx, mdmFileVaultConfigRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM FileVault configuration, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a MDM FileVault configuration")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmFileVaultConfigForState(ztlMFC))...)
}

func (r *MDMFileVaultConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmFileVaultConfig

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMFC, _, err := r.client.MDMFileVaultConfigs.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM FileVault configuration %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a MDM FileVault configuration")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmFileVaultConfigForState(ztlMFC))...)
}

func (r *MDMFileVaultConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmFileVaultConfig

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMFC, _, err := r.client.MDMFileVaultConfigs.Update(ctx, int(data.ID.ValueInt64()), mdmFileVaultConfigRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM FileVault configuration %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a MDM FileVault configuration")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmFileVaultConfigForState(ztlMFC))...)
}

func (r *MDMFileVaultConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmFileVaultConfig

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMFileVaultConfigs.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM FileVault configuration %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a MDM FileVault configuration")
}

func (r *MDMFileVaultConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "MDM FileVault configuration", req, resp)
}
