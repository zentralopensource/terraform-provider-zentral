package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMBlueprintResource{}
var _ resource.ResourceWithImportState = &MDMBlueprintResource{}

func NewMDMBlueprintResource() resource.Resource {
	return &MDMBlueprintResource{}
}

// MDMBlueprintResource defines the resource implementation.
type MDMBlueprintResource struct {
	client *goztl.Client
}

func (r *MDMBlueprintResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_blueprint"
}

func (r *MDMBlueprintResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM blueprints.",
		MarkdownDescription: "The resource `zentral_mdm_blueprint` manages MDM blueprints.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the blueprint.",
				MarkdownDescription: "`ID` of the blueprint.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the blueprint.",
				MarkdownDescription: "Name of the blueprint.",
				Required:            true,
			},
			"inventory_interval": schema.Int64Attribute{
				Description:         "In seconds, the minimum interval between two inventory collection. Minimum 4h, maximum 7d, default 1d.",
				MarkdownDescription: "In seconds, the minimum interval between two inventory collection. Minimum 4h, maximum 7d, default 1d.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(86400),
				Validators: []validator.Int64{
					int64validator.Between(14400, 604800),
				},
			},
			"collect_apps": schema.StringAttribute{
				Description:         "Inventory apps collection setting. Possible values: NO, MANAGED_ONLY, ALL.",
				MarkdownDescription: "Possible values: `NO`, `MANAGED_ONLY`, `ALL`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"NO", "MANAGED_ONLY", "ALL"}...),
				},
				Default: stringdefault.StaticString("NO"),
			},
			"collect_certificates": schema.StringAttribute{
				Description:         "Inventory certificates collection setting. Possible values: NO, MANAGED_ONLY, ALL.",
				MarkdownDescription: "Possible values: `NO`, `MANAGED_ONLY`, `ALL`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"NO", "MANAGED_ONLY", "ALL"}...),
				},
				Default: stringdefault.StaticString("NO"),
			},
			"collect_profiles": schema.StringAttribute{
				Description:         "Inventory profiles collection setting. Possible values: NO, MANAGED_ONLY, ALL.",
				MarkdownDescription: "Possible values: `NO`, `MANAGED_ONLY`, `ALL`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"NO", "MANAGED_ONLY", "ALL"}...),
				},
				Default: stringdefault.StaticString("NO"),
			},
			"legacy_profiles_via_ddm": schema.BoolAttribute{
				Description:         "If true, legacy profiles are distributed via DDM. Defaults to true.",
				MarkdownDescription: "If `true`, legady profiles are distributed via DDM. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"default_location_id": schema.Int64Attribute{
				Description:         "The ID of the default apps & books location.",
				MarkdownDescription: "The `ID` of the default apps & books location.",
				Optional:            true,
			},
			"filevault_config_id": schema.Int64Attribute{
				Description:         "The ID of the attached FileVault configuration.",
				MarkdownDescription: "The `ID` of the attached FileVault configuration.",
				Optional:            true,
			},
			"recovery_password_config_id": schema.Int64Attribute{
				Description:         "The ID of the attached recovery password configuration.",
				MarkdownDescription: "The `ID` of the attached recovery password configuration.",
				Optional:            true,
			},
			"software_update_enforcement_ids": schema.SetAttribute{
				Description:         "The IDs of the software update enforcements.",
				MarkdownDescription: "The `ID`s of the software update enforcements.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *MDMBlueprintResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMBlueprintResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmBlueprint

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMB, _, err := r.client.MDMBlueprints.Create(ctx, mdmBlueprintRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM blueprint, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a MDM blueprint")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmBlueprintForState(ztlMB))...)
}

func (r *MDMBlueprintResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmBlueprint

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMB, _, err := r.client.MDMBlueprints.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM blueprint %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a MDM blueprint")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmBlueprintForState(ztlMB))...)
}

func (r *MDMBlueprintResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmBlueprint

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMB, _, err := r.client.MDMBlueprints.Update(ctx, int(data.ID.ValueInt64()), mdmBlueprintRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM blueprint %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a MDM blueprint")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmBlueprintForState(ztlMB))...)
}

func (r *MDMBlueprintResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmBlueprint

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMBlueprints.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM blueprint %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a MDM blueprint")
}

func (r *MDMBlueprintResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "MDM blueprint", req, resp)
}
