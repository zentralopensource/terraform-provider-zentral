package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
	"github.com/zentralopensource/terraform-provider-zentral/internal/planmodifiers"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &SantaConfigurationResource{}
var _ resource.ResourceWithImportState = &SantaConfigurationResource{}

func NewSantaConfigurationResource() resource.Resource {
	return &SantaConfigurationResource{}
}

// SantaConfigurationResource defines the resource implementation.
type SantaConfigurationResource struct {
	client *goztl.Client
}

func (r *SantaConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_santa_configuration"
}

func (r *SantaConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Santa configurations.",
		MarkdownDescription: "The resource `zentral_santa_configuration` manages Santa configurations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Santa configuration.",
				MarkdownDescription: "`ID` of the Santa configuration.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the Santa configuration.",
				MarkdownDescription: "Name of the Santa configuration.",
				Required:            true,
			},
			"client_mode": schema.Int64Attribute{
				Description:         "Client mode of the Santa configuration.",
				MarkdownDescription: "Client mode of the Santa configuration.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					planmodifiers.Int64Default(int64(1)), // MONITOR
				},
			},
			"client_certificate_auth": schema.BoolAttribute{
				Description:         "If `true`, mTLS is required between Santa and Zentral.",
				MarkdownDescription: "If `true`, mTLS is required between Santa and Zentral.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefault(false),
				},
			},
			"batch_size": schema.Int64Attribute{
				Description:         "The number of rules to download or events to upload per request.",
				MarkdownDescription: "The number of rules to download or events to upload per request.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					planmodifiers.Int64Default(int64(50)),
				},
			},
			"full_sync_interval": schema.Int64Attribute{
				Description:         "The max time to wait before performing a full sync with the server.",
				MarkdownDescription: "The max time to wait before performing a full sync with the server.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					planmodifiers.Int64Default(int64(600)),
				},
			},
			"enable_bundles": schema.BoolAttribute{
				Description:         "If set to true the bundle scanning feature is enabled.",
				MarkdownDescription: "If set to `true` the bundle scanning feature is enabled.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefault(false),
				},
			},
			"enable_transitive_rules": schema.BoolAttribute{
				Description:         "If set to true the transitive rule feature is enabled.",
				MarkdownDescription: "If set to `true` the transitive rule feature is enabled.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefault(false),
				},
			},
			"allowed_path_regex": schema.StringAttribute{
				Description:         "A regex to allow if the binary, certificate, or Team ID scopes did not allow/block execution.",
				MarkdownDescription: "A regex to allow if the binary, certificate, or Team ID scopes did not allow/block execution.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.StringDefault(""),
				},
			},
			"blocked_path_regex": schema.StringAttribute{
				Description:         "A regex to block if the binary, certificate, or Team ID scopes did not allow/block an execution.",
				MarkdownDescription: "A regex to block if the binary, certificate, or Team ID scopes did not allow/block an execution.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.StringDefault(""),
				},
			},
			"block_usb_mount": schema.BoolAttribute{
				Description:         "If set to true blocking USB Mass storage feature is enabled.",
				MarkdownDescription: "If set to `true` blocking USB Mass storage feature is enabled.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefault(false),
				},
			},
			"remount_usb_mode": schema.SetAttribute{
				Description:         "Array of strings for arguments to pass to mount -o.",
				MarkdownDescription: "Array of strings for arguments to pass to `mount -o`.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"allow_unknown_shard": schema.Int64Attribute{
				Description:         "Restrict the reporting of 'Allow Unknown' events to a percentage (0-100) of hosts.",
				MarkdownDescription: "Restrict the reporting of 'Allow Unknown' events to a percentage (0-100) of hosts.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					planmodifiers.Int64Default(int64(100)),
				},
			},
			"enable_all_event_upload_shard": schema.Int64Attribute{
				Description:         "Restrict the upload of all execution events to Zentral, including those that were explicitly allowed, to a percentage (0-100) of hosts",
				MarkdownDescription: "Restrict the upload of all execution events to Zentral, including those that were explicitly allowed, to a percentage (0-100) of hosts",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					planmodifiers.Int64Default(int64(0)),
				},
			},
			"sync_incident_severity": schema.Int64Attribute{
				Description:         "If 100, 200, 300, incidents will be automatically opened and closed when the santa agent rules are out of sync.",
				MarkdownDescription: "If 100, 200, 300, incidents will be automatically opened and closed when the santa agent rules are out of sync.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					planmodifiers.Int64Default(int64(0)), // Severity.None
				},
			},
		},
	}
}

func (r *SantaConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SantaConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data santaConfiguration

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlSC, _, err := r.client.SantaConfigurations.Create(ctx, santaConfigurationRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Santa configuration, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Santa configuration")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, santaConfigurationForState(ztlSC))...)
}

func (r *SantaConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data santaConfiguration

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlSC, _, err := r.client.SantaConfigurations.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Santa configuration %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Santa configuration")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, santaConfigurationForState(ztlSC))...)
}

func (r *SantaConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data santaConfiguration

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlSC, _, err := r.client.SantaConfigurations.Update(ctx, int(data.ID.ValueInt64()), santaConfigurationRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Santa configuration %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Santa configuration")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, santaConfigurationForState(ztlSC))...)
}

func (r *SantaConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data santaConfiguration

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.SantaConfigurations.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Santa configuration %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Santa configuration")
}

func (r *SantaConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Santa configuration", req, resp)
}
