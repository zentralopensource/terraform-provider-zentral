package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
	"github.com/zentralopensource/terraform-provider-zentral/internal/planmodifiers"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OsqueryPackQueryResource{}
var _ resource.ResourceWithImportState = &OsqueryPackQueryResource{}

func NewOsqueryPackQueryResource() resource.Resource {
	return &OsqueryPackQueryResource{}
}

// OsqueryPackQueryResource defines the resource implementation.
type OsqueryPackQueryResource struct {
	client *goztl.Client
}

func (r *OsqueryPackQueryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_osquery_pack_query"
}

func (r *OsqueryPackQueryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Osquery pack queries.",
		MarkdownDescription: "The resource `zentral_osquery_query` manages Osquery pack queries.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the pack query.",
				MarkdownDescription: "`ID` of the pack query.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"pack_id": schema.Int64Attribute{
				Description:         "ID of the pack.",
				MarkdownDescription: "ID of the pack.",
				Required:            true,
			},
			"query_id": schema.Int64Attribute{
				Description:         "ID of the query.",
				MarkdownDescription: "ID of the query.",
				Required:            true,
			},
			"slug": schema.StringAttribute{
				Description:         "Slug of the pack query.",
				MarkdownDescription: "Slug of the pack query.",
				Computed:            true,
			},
			"interval": schema.Int64Attribute{
				Description:         "Query frequency, in seconds.",
				MarkdownDescription: "Query frequency, in seconds.",
				Required:            true,
			},
			"log_removed_actions": schema.BoolAttribute{
				Description:         "If true, 'removed' actions should be logged. Defaults to true.",
				MarkdownDescription: "If `true`, `removed` actions should be logged. Default to `true`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefault(true),
				},
			},
			"snapshot_mode": schema.BoolAttribute{
				Description:         "If true, differentials will not be stored and this query will not emulate an event stream. Defaults to false.",
				MarkdownDescription: "If `true`, differentials will not be stored and this query will not emulate an event stream. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefault(false),
				},
			},
			"shard": schema.Int64Attribute{
				Description:         "Restrict this query to a percentage (1-100) of target hosts.",
				MarkdownDescription: "Restrict this query to a percentage (1-100) of target hosts.",
				Optional:            true,
			},
			"can_be_denylisted": schema.BoolAttribute{
				Description:         "If true, this query can be denylisted when stopped by the watchdog for excessive resource consumption. Defaults to true.",
				MarkdownDescription: "If `true`, this query can be denylisted when stopped by the watchdog for excessive resource consumption. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefault(true),
				},
			},
		},
	}
}

func (r *OsqueryPackQueryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OsqueryPackQueryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data osqueryPackQuery

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOPQ, _, err := r.client.OsqueryPackQueries.Create(ctx, osqueryPackQueryRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Osquery pack query, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an Osquery pack query")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryPackQueryForState(ztlOPQ))...)
}

func (r *OsqueryPackQueryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data osqueryPackQuery

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOPQ, _, err := r.client.OsqueryPackQueries.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Osquery pack query %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read an Osquery pack query")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryPackQueryForState(ztlOPQ))...)
}

func (r *OsqueryPackQueryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data osqueryPackQuery

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOPQ, _, err := r.client.OsqueryPackQueries.Update(ctx, int(data.ID.ValueInt64()), osqueryPackQueryRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Osquery pack query %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an Osquery pack query")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryPackQueryForState(ztlOPQ))...)
}

func (r *OsqueryPackQueryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data osqueryPackQuery

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.OsqueryPackQueries.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Osquery pack query %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an Osquery pack query")
}

func (r *OsqueryPackQueryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Osquery pack query", req, resp)
}
