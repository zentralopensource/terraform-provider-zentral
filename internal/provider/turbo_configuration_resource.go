package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &TurboConfigurationResource{}
var _ resource.ResourceWithImportState = &TurboConfigurationResource{}

func NewTurboConfigurationResource() resource.Resource {
	return &TurboConfigurationResource{}
}

// TurboConfigurationResource defines the resource implementation.
type TurboConfigurationResource struct {
	client *goztl.Client
}

func (r *TurboConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_turbo_configuration"
}

func (r *TurboConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Turbo configurations.",
		MarkdownDescription: "The resource `zentral_turbo_configuration` manages Turbo configurations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the Turbo configuration.",
				MarkdownDescription: "`ID` of the Turbo configuration.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the Turbo configuration.",
				MarkdownDescription: "Name of the Turbo configuration.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the Turbo configuration.",
				MarkdownDescription: "Description of the Turbo configuration.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"collect_inventory": schema.BoolAttribute{
				Description:         "If true, the agent posts a full machine inventory snapshot on the inventory interval. Defaults to true.",
				MarkdownDescription: "If `true`, the agent posts a full machine inventory snapshot on the inventory interval. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"inventory_interval": schema.Int64Attribute{
				Description:         "Inventory refresh interval in seconds. Must be between 60 and 604800. Defaults to 86400.",
				MarkdownDescription: "Inventory refresh interval in seconds. Must be between `60` and `604800`. Defaults to `86400`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(86400),
				Validators: []validator.Int64{
					int64validator.Between(60, 604800),
				},
			},
			"default_check_interval": schema.Int64Attribute{
				Description:         "Default run interval in seconds for recurring jobs that don't set their own. Must be between 60 and 604800. Defaults to 86400.",
				MarkdownDescription: "Default run interval in seconds for recurring jobs that don't set their own. Must be between `60` and `604800`. Defaults to `86400`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(86400),
				Validators: []validator.Int64{
					int64validator.Between(60, 604800),
				},
			},
			"config_refresh_interval": schema.Int64Attribute{
				Description:         "How long the agent may trust a cached configuration before refreshing it, in seconds. Must be between 60 and 604800. Defaults to 600.",
				MarkdownDescription: "How long the agent may trust a cached configuration before refreshing it, in seconds. Must be between `60` and `604800`. Defaults to `600`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(600),
				Validators: []validator.Int64{
					int64validator.Between(60, 604800),
				},
			},
			"results_batch_size": schema.Int64Attribute{
				Description:         "Maximum number of results the agent uploads per request. Must be between 1 and 1000. Defaults to 100.",
				MarkdownDescription: "Maximum number of results the agent uploads per request. Must be between `1` and `1000`. Defaults to `100`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(100),
				Validators: []validator.Int64{
					int64validator.Between(1, 1000),
				},
			},
		},
	}
}

func (r *TurboConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TurboConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data turboConfiguration

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlTC, _, err := r.client.TurboConfigurations.Create(ctx, turboConfigurationRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Turbo configuration, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Turbo configuration")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, turboConfigurationForState(ztlTC))...)
}

func (r *TurboConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data turboConfiguration

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlTC, _, err := r.client.TurboConfigurations.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Turbo configuration %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Turbo configuration")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, turboConfigurationForState(ztlTC))...)
}

func (r *TurboConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data turboConfiguration

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlTC, _, err := r.client.TurboConfigurations.Update(ctx, data.ID.ValueString(), turboConfigurationRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Turbo configuration %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Turbo configuration")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, turboConfigurationForState(ztlTC))...)
}

func (r *TurboConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data turboConfiguration

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.TurboConfigurations.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Turbo configuration %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Turbo configuration")
}

func (r *TurboConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "Turbo configuration", req, resp)
}
