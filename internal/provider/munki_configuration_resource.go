package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MunkiConfigurationResource{}
var _ resource.ResourceWithImportState = &MunkiConfigurationResource{}

func NewMunkiConfigurationResource() resource.Resource {
	return &MunkiConfigurationResource{}
}

// MunkiConfigurationResource defines the resource implementation.
type MunkiConfigurationResource struct {
	client *goztl.Client
}

func (r *MunkiConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_munki_configuration"
}

func (r *MunkiConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Munki configurations.",
		MarkdownDescription: "The resource `zentral_munki_configuration` manages Munki configurations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Munki configuration.",
				MarkdownDescription: "`ID` of the Munki configuration.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the Munki configuration.",
				MarkdownDescription: "Name of the Munki configuration.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the Munki configuration.",
				MarkdownDescription: "Description of the Munki configuration.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"inventory_apps_full_info_shard": schema.Int64Attribute{
				Description:         "Percentage of machines configured to collect the full inventory apps information. Defaults to 100.",
				MarkdownDescription: "Percentage of machines configured to collect the full inventory apps information. Defaults to `100`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(100),
			},
			"principal_user_detection_sources": schema.ListAttribute{
				Description:         "List of principal user detection sources.",
				MarkdownDescription: "List of principal user detection sources.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"principal_user_detection_domains": schema.SetAttribute{
				Description:         "Set of principal user detection domains.",
				MarkdownDescription: "Set of principal user detection domains.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"collected_condition_keys": schema.SetAttribute{
				Description:         "Set of the condition keys to collect.",
				MarkdownDescription: "Set of the condition keys to collect.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"managed_installs_sync_interval_days": schema.Int64Attribute{
				Description:         "Interval in days between full managed installs sync. Defaults to 7 days.",
				MarkdownDescription: "Interval in days between full managed installs sync. Defaults to 7 days.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(7),
			},
			"auto_reinstall_incidents": schema.BoolAttribute{
				Description:         "If true, incidents will be managed automatically when package reinstalls are observed. Defaults to false.",
				MarkdownDescription: "If `true`, incidents will be managed automatically when package reinstalls are observed. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"auto_failed_install_incidents": schema.BoolAttribute{
				Description:         "If true, incidents will be managed automatically when package failed installs are observed. Defaults to false.",
				MarkdownDescription: "If `true`, incidents will be managed automatically when package failed installs are observed. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"version": schema.Int64Attribute{
				Description:         "Version of the Munki configuration.",
				MarkdownDescription: "Version of the Munki configuration.",
				Computed:            true,
			},
		},
	}
}

func (r *MunkiConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MunkiConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data munkiConfiguration

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMC, _, err := r.client.MunkiConfigurations.Create(ctx, munkiConfigurationRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Munki configuration, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Munki configuration")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, munkiConfigurationForState(ztlMC))...)
}

func (r *MunkiConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data munkiConfiguration

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMC, _, err := r.client.MunkiConfigurations.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Munki configuration %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Munki configuration")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, munkiConfigurationForState(ztlMC))...)
}

func (r *MunkiConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data munkiConfiguration

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMC, _, err := r.client.MunkiConfigurations.Update(ctx, int(data.ID.ValueInt64()), munkiConfigurationRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Munki configuration %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Munki configuration")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, munkiConfigurationForState(ztlMC))...)
}

func (r *MunkiConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data munkiConfiguration

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MunkiConfigurations.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Munki configuration %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Munki configuration")
}

func (r *MunkiConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Munki configuration", req, resp)
}
