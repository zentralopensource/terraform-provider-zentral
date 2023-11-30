package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMSoftwareUpdateEnforcementResource{}
var _ resource.ResourceWithImportState = &MDMSoftwareUpdateEnforcementResource{}

func NewMDMSoftwareUpdateEnforcementResource() resource.Resource {
	return &MDMSoftwareUpdateEnforcementResource{}
}

// MDMSoftwareUpdateEnforcementResource defines the resource implementation.
type MDMSoftwareUpdateEnforcementResource struct {
	client *goztl.Client
}

func (r *MDMSoftwareUpdateEnforcementResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_software_update_enforcement"
}

func (r *MDMSoftwareUpdateEnforcementResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM software update enforcements.",
		MarkdownDescription: "The resource `zentral_mdm_software_update_enforcement` manages MDM software update enforcements.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the software update enforcement.",
				MarkdownDescription: "`ID` of the software update enforcement.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the software update enforcement.",
				MarkdownDescription: "Name of the software update enforcement.",
				Required:            true,
			},
			"details_url": schema.StringAttribute{
				Description:         "The URL of a web page that shows details that the organization provides about the enforced update.",
				MarkdownDescription: "The URL of a web page that shows details that the organization provides about the enforced update.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags used to scope the software update enforcement.",
				MarkdownDescription: "The `ID`s of the tags used to scope the software update enforcement.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"os_version": schema.StringAttribute{
				Description:         "The target OS version to update the device to by the appropriate time.",
				MarkdownDescription: "The target OS version to update the device to by the appropriate time.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"build_version": schema.StringAttribute{
				Description:         "The target build version to update the device to by the appropriate time.",
				MarkdownDescription: "The target build version to update the device to by the appropriate time.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"local_datetime": schema.StringAttribute{
				Description:         "The local date time value that specifies when to force install the software update.",
				MarkdownDescription: "The local date time value that specifies when to force install the software update.",
				Optional:            true,
			},
			"max_os_version": schema.StringAttribute{
				Description:         "The maximum (excluded) target OS version to update the device to by the appropriate time.",
				MarkdownDescription: "The maximum (excluded) target OS version to update the device to by the appropriate time.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"delay_days": schema.Int64Attribute{
				Description:         "Number of days after a software update release before the device force installs it.",
				MarkdownDescription: "Number of days after a software update release before the device force installs it.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 120),
				},
			},
			"local_time": schema.StringAttribute{
				Description:         "The local time value that specifies when to force install the software update.",
				MarkdownDescription: "The local time value that specifies when to force install the software update.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *MDMSoftwareUpdateEnforcementResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMSoftwareUpdateEnforcementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmSoftwareUpdateEnforcement

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSUE, _, err := r.client.MDMSoftwareUpdateEnforcements.Create(ctx, mdmSoftwareUpdateEnforcementRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM software update enforcement, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a MDM software update enforcement")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmSoftwareUpdateEnforcementForState(ztlMSUE))...)
}

func (r *MDMSoftwareUpdateEnforcementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmSoftwareUpdateEnforcement

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSUE, _, err := r.client.MDMSoftwareUpdateEnforcements.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM software update enforcement %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a MDM software update enforcement")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmSoftwareUpdateEnforcementForState(ztlMSUE))...)
}

func (r *MDMSoftwareUpdateEnforcementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmSoftwareUpdateEnforcement

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSUE, _, err := r.client.MDMSoftwareUpdateEnforcements.Update(ctx, int(data.ID.ValueInt64()), mdmSoftwareUpdateEnforcementRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM software update enforcement %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a MDM software update enforcement")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmSoftwareUpdateEnforcementForState(ztlMSUE))...)
}

func (r *MDMSoftwareUpdateEnforcementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmSoftwareUpdateEnforcement

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMSoftwareUpdateEnforcements.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM software update enforcement %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a MDM software update enforcement")
}

func (r *MDMSoftwareUpdateEnforcementResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "MDM software update enforcement", req, resp)
}
