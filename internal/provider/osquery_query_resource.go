package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OsqueryQueryResource{}
var _ resource.ResourceWithImportState = &OsqueryQueryResource{}

func NewOsqueryQueryResource() resource.Resource {
	return &OsqueryQueryResource{}
}

// OsqueryQueryResource defines the resource implementation.
type OsqueryQueryResource struct {
	client *goztl.Client
}

func (r *OsqueryQueryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_osquery_query"
}

func (r *OsqueryQueryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Osquery queries.",
		MarkdownDescription: "The resource `zentral_osquery_query` manages Osquery queries.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the query.",
				MarkdownDescription: "`ID` of the query.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the query.",
				MarkdownDescription: "Name of the query.",
				Required:            true,
			},
			"sql": schema.StringAttribute{
				Description:         "The SQL query to run.",
				MarkdownDescription: "The SQL query to run.",
				Required:            true,
			},
			"platforms": schema.SetAttribute{
				Description:         "Restrict the query to some platforms, default is 'all' platforms",
				MarkdownDescription: "Restrict the query to some platforms, default is 'all' platforms",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"minimum_osquery_version": schema.StringAttribute{
				Description:         "Only run on Osquery versions greater than or equal-to this version string",
				MarkdownDescription: "Only run on Osquery versions greater than or equal-to this version string",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the query.",
				MarkdownDescription: "Description of the query.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"value": schema.StringAttribute{
				Description:         "Description of the results returned by the query.",
				MarkdownDescription: "Description of the results returned by the query.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"version": schema.Int64Attribute{
				Description:         "Version of the query.",
				MarkdownDescription: "Version of the query.",
				Computed:            true,
			},
			"compliance_check_enabled": schema.BoolAttribute{
				Description:         "If true, the query will be used as compliance check. Defaults to false.",
				MarkdownDescription: "If `true`, the query will be used as compliance check. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"scheduling": schema.SingleNestedAttribute{
				Description:         "Attributes to link a query to a pack for scheduling.",
				MarkdownDescription: "Attributes to link a query to a pack for scheduling.",
				Attributes: map[string]schema.Attribute{
					"pack_id": schema.Int64Attribute{
						Description:         "The ID of the pack.",
						MarkdownDescription: "The `ID` of the pack.",
						Required:            true,
					},
					"interval": schema.Int64Attribute{
						Description:         "the query frequency, in seconds. It has a maximum value of 604,800 (1 week).",
						MarkdownDescription: "the query frequency, in seconds. It has a maximum value of 604,800 (1 week).",
						Required:            true,
						Validators: []validator.Int64{
							int64validator.Between(10, 604800),
						},
					},
					"log_removed_actions": schema.BoolAttribute{
						Description:         "If true, removed actions should be logged. Default to true.",
						MarkdownDescription: "If `true`, remove actions should be logged. Default to `true`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
					},
					"snapshot_mode": schema.BoolAttribute{
						Description:         "If true, differentials will not be stored and this query will not emulate an event stream. Defaults to false.",
						MarkdownDescription: "If `true`, differentials will not be stored and this query will not emulate an event stream. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"shard": schema.Int64Attribute{
						Description:         "Restrict this query to a percentage (1-100) of target hosts.",
						MarkdownDescription: "Restrict this query to a percentage (1-100) of target hosts.",
						Optional:            true,
						Validators: []validator.Int64{
							int64validator.Between(1, 100),
						},
					},
					"can_be_denylisted": schema.BoolAttribute{
						Description:         "If true, this query can be denylisted when stopped by the watchdog for excessive resource consumption. Defaults to true.",
						MarkdownDescription: "If `true`, this query can be denylisted when stopped by the watchdog for excessive resource consumption. Defaults to `true`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
					},
				},
				Optional: true,
			},
		},
	}
}

func (r *OsqueryQueryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OsqueryQueryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data osqueryQuery

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOQ, _, err := r.client.OsqueryQueries.Create(ctx, osqueryQueryRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Osquery query, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an Osquery query")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryQueryForState(ztlOQ))...)
}

func (r *OsqueryQueryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data osqueryQuery

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOQ, _, err := r.client.OsqueryQueries.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Osquery query %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read an Osquery query")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryQueryForState(ztlOQ))...)
}

func (r *OsqueryQueryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data osqueryQuery

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOQ, _, err := r.client.OsqueryQueries.Update(ctx, int(data.ID.ValueInt64()), osqueryQueryRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Osquery query %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an Osquery query")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryQueryForState(ztlOQ))...)
}

func (r *OsqueryQueryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data osqueryQuery

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.OsqueryQueries.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Osquery query %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an Osquery query")
}

func (r *OsqueryQueryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Osquery query", req, resp)
}
