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
var _ resource.Resource = &OsqueryConfigurationResource{}
var _ resource.ResourceWithImportState = &OsqueryConfigurationResource{}

func NewOsqueryConfigurationResource() resource.Resource {
	return &OsqueryConfigurationResource{}
}

// OsqueryConfigurationResource defines the resource implementation.
type OsqueryConfigurationResource struct {
	client *goztl.Client
}

func (r *OsqueryConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_osquery_configuration"
}

func (r *OsqueryConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Osquery configurations.",
		MarkdownDescription: "The resource `zentral_osquery_configuration` manages Osquery configurations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Osquery configuration.",
				MarkdownDescription: "`ID` of the Osquery configuration.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the Osquery configuration.",
				MarkdownDescription: "Name of the Osquery configuration.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the Osquery configuration.",
				MarkdownDescription: "Description of the Osquery configuration.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.StringDefault(""),
				},
			},
			"inventory": schema.BoolAttribute{
				Description:         "If true, Osquery is configured to collect inventory data. Defaults to true.",
				MarkdownDescription: "If `true`, Osquery is configured to collect inventory data. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefault(true),
				},
			},
			"inventory_apps": schema.BoolAttribute{
				Description:         "If true, Osquery is configured to collect the applications. Defaults to false.",
				MarkdownDescription: "If `true`, Osquery is configured to collect the applications. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefault(false),
				},
			},
			"inventory_ec2": schema.BoolAttribute{
				Description:         "If true, Osquery is configured to collect the EC2 metadata. Defaults to false.",
				MarkdownDescription: "If `true`, Osquery is configured to collect the EC2 metadata. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefault(false),
				},
			},
			"inventory_interval": schema.Int64Attribute{
				Description:         "Number of seconds to wait between collecting the inventory data.",
				MarkdownDescription: "Number of seconds to wait between collecting the inventory data.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					planmodifiers.Int64Default(int64(86400)),
				},
			},
			"options": schema.MapAttribute{
				Description:         "A map of extra options to pass to Osquery in the flag file.",
				MarkdownDescription: "A map of extra options to pass to Osquery in the flag file.",
				// Options ElementType is types.StringType
				// This is much easier this way, and since the options are serialized in the flag file, this is not restrictive.
				// Non-string elements coming from the server will be converted to strings.
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"automatic_table_constructions": schema.SetAttribute{
				Description:         "List of ATC IDs to include in this configuration.",
				MarkdownDescription: "List of ATC IDs to include in this configuration.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *OsqueryConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OsqueryConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data osqueryConfiguration

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOC, _, err := r.client.OsqueryConfigurations.Create(ctx, osqueryConfigurationRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Osquery configuration, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an Osquery configuration")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryConfigurationForState(ztlOC))...)
}

func (r *OsqueryConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data osqueryConfiguration

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOC, _, err := r.client.OsqueryConfigurations.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Osquery configuration %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read an Osquery configuration")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryConfigurationForState(ztlOC))...)
}

func (r *OsqueryConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data osqueryConfiguration

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOC, _, err := r.client.OsqueryConfigurations.Update(ctx, int(data.ID.ValueInt64()), osqueryConfigurationRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Osquery configuration %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an Osquery configuration")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryConfigurationForState(ztlOC))...)
}

func (r *OsqueryConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data osqueryConfiguration

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.OsqueryConfigurations.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Osquery configuration %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an Osquery configuration")
}

func (r *OsqueryConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Osquery configuration", req, resp)
}
