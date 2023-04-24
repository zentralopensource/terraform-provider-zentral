package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OsqueryPackResource{}
var _ resource.ResourceWithImportState = &OsqueryPackResource{}

func NewOsqueryPackResource() resource.Resource {
	return &OsqueryPackResource{}
}

// OsqueryPackResource defines the resource implementation.
type OsqueryPackResource struct {
	client *goztl.Client
}

func (r *OsqueryPackResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_osquery_pack"
}

func (r *OsqueryPackResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Osquery packs.",
		MarkdownDescription: "The resource `zentral_osquery_pack` manages Osquery packs.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the pack.",
				MarkdownDescription: "`ID` of the pack.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the pack.",
				MarkdownDescription: "Name of the pack.",
				Required:            true,
			},
			"slug": schema.StringAttribute{
				Description:         "Slug of the pack.",
				MarkdownDescription: "Slug of the pack.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the pack.",
				MarkdownDescription: "Description of the pack.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"discovery_queries": schema.ListAttribute{
				Description:         "List of osquery queries which control whether or not the pack will execute.",
				MarkdownDescription: "List of osquery queries which control whether or not the pack will execute.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"shard": schema.Int64Attribute{
				Description:         "Restrict the pack to a percentage (1-100) of target hosts.",
				MarkdownDescription: "Restrict the pack to a percentage (1-100) of target hosts.",
				Optional:            true,
			},
			"event_routing_key": schema.StringAttribute{
				Description:         "Routing key added to the metadata of the events that the queries of this pack generate.",
				MarkdownDescription: "Routing key added to the metadata of the events that the queries of this pack generate.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
		},
	}
}

func (r *OsqueryPackResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OsqueryPackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data osqueryPack

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOP, _, err := r.client.OsqueryPacks.Create(ctx, osqueryPackRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Osquery pack, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an Osquery pack")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryPackForState(ztlOP))...)
}

func (r *OsqueryPackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data osqueryPack

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOP, _, err := r.client.OsqueryPacks.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Osquery pack %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read an Osquery pack")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryPackForState(ztlOP))...)
}

func (r *OsqueryPackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data osqueryPack

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOP, _, err := r.client.OsqueryPacks.Update(ctx, int(data.ID.ValueInt64()), osqueryPackRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Osquery pack %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an Osquery pack")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryPackForState(ztlOP))...)
}

func (r *OsqueryPackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data osqueryPack

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.OsqueryPacks.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Osquery pack %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an Osquery pack")
}

func (r *OsqueryPackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Osquery pack", req, resp)
}
