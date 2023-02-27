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
var _ resource.Resource = &OsqueryATCResource{}
var _ resource.ResourceWithImportState = &OsqueryATCResource{}

func NewOsqueryATCResource() resource.Resource {
	return &OsqueryATCResource{}
}

// OsqueryATCResource defines the resource implementation.
type OsqueryATCResource struct {
	client *goztl.Client
}

func (r *OsqueryATCResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_osquery_atc"
}

func (r *OsqueryATCResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Osquery automatic table constructions.",
		MarkdownDescription: "The resource `zentral_osquery_atc` manages Osquery automatic table constructions.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Osquery ATC.",
				MarkdownDescription: "`ID` of the Osquery ATC.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the Osquery ATC.",
				MarkdownDescription: "Name of the Osquery ATC.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the Osquery ATC.",
				MarkdownDescription: "Description of the Osquery ATC.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.StringDefault(""),
				},
			},
			"table_name": schema.StringAttribute{
				Description:         "Name of the Osquery ATC table.",
				MarkdownDescription: "Name of the Osquery ATC table.",
				Required:            true,
			},
			"query": schema.StringAttribute{
				Description:         "Query used to fetch the ATC data.",
				MarkdownDescription: "Query used to fetch the ATC data.",
				Required:            true,
			},
			"path": schema.StringAttribute{
				Description:         "Path of the SQLite table on the device.",
				MarkdownDescription: "Path of the SQLite table on the device.",
				Required:            true,
			},
			"columns": schema.ListAttribute{
				Description:         "List of the column names corresponding the the query.",
				MarkdownDescription: "List of the column names corresponding the the query.",
				ElementType:         types.StringType,
				Required:            true,
			},
			"platforms": schema.SetAttribute{
				Description:         "Platform on which this ATC can be activated",
				MarkdownDescription: "Platform on which this ATC can be activated",
				ElementType:         types.StringType,
				Required:            true,
			},
		},
	}
}

func (r *OsqueryATCResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OsqueryATCResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data osqueryATC

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOA, _, err := r.client.OsqueryATC.Create(ctx, osqueryATCRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Osquery ATC, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an Osquery ATC")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryATCForState(ztlOA))...)
}

func (r *OsqueryATCResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data osqueryATC

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOA, _, err := r.client.OsqueryATC.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Osquery ATC %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read an Osquery ATC")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryATCForState(ztlOA))...)
}

func (r *OsqueryATCResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data osqueryATC

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOA, _, err := r.client.OsqueryATC.Update(ctx, int(data.ID.ValueInt64()), osqueryATCRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Osquery ATC %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an Osquery ATC")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryATCForState(ztlOA))...)
}

func (r *OsqueryATCResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data osqueryATC

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.OsqueryATC.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Osquery ATC %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an Osquery ATC")
}

func (r *OsqueryATCResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Osquery ATC", req, resp)
}
