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
var _ resource.Resource = &OsqueryFileCategoryResource{}
var _ resource.ResourceWithImportState = &OsqueryFileCategoryResource{}

func NewOsqueryFileCategoryResource() resource.Resource {
	return &OsqueryFileCategoryResource{}
}

// OsqueryFileCategoryResource defines the resource implementation.
type OsqueryFileCategoryResource struct {
	client *goztl.Client
}

func (r *OsqueryFileCategoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_osquery_file_category"
}

func (r *OsqueryFileCategoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Osquery file categories.",
		MarkdownDescription: "The resource `zentral_osquery_file_category` manages Osquery file categories.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Osquery file category.",
				MarkdownDescription: "`ID` of the Osquery file category.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the Osquery file category.",
				MarkdownDescription: "Name of the Osquery file category.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the Osquery file category.",
				MarkdownDescription: "Description of the Osquery file category.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.StringDefault(""),
				},
			},
			"file_paths": schema.SetAttribute{
				Description:         "Set of paths to include in the Osquery file category.",
				MarkdownDescription: "Set of paths to include in the Osquery file category.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"exclude_paths": schema.SetAttribute{
				Description:         "Set of paths to exclude from the Osquery file category.",
				MarkdownDescription: "Set of paths to exclude from the Osquery file category.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"file_paths_queries": schema.SetAttribute{
				Description:         "Set of queries returning paths to monitor as path columns in the results.",
				MarkdownDescription: "Set of queries returning paths to monitor as path columns in the results.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"access_monitoring": schema.BoolAttribute{
				Description:         "If true, FIM will include file access for this file category. Defaults to false.",
				MarkdownDescription: "If `true`, FIM will include file access for this file category. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefault(false),
				},
			},
		},
	}
}

func (r *OsqueryFileCategoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OsqueryFileCategoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data osqueryFileCategory

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOFC, _, err := r.client.OsqueryFileCategories.Create(ctx, osqueryFileCategoryRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Osquery file category, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an Osquery file category")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryFileCategoryForState(ztlOFC))...)
}

func (r *OsqueryFileCategoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data osqueryFileCategory

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOFC, _, err := r.client.OsqueryFileCategories.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Osquery file category %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read an Osquery file category")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryFileCategoryForState(ztlOFC))...)
}

func (r *OsqueryFileCategoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data osqueryFileCategory

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOFC, _, err := r.client.OsqueryFileCategories.Update(ctx, int(data.ID.ValueInt64()), osqueryFileCategoryRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Osquery file category %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an Osquery file category")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryFileCategoryForState(ztlOFC))...)
}

func (r *OsqueryFileCategoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data osqueryFileCategory

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.OsqueryFileCategories.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Osquery file category %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an Osquery file category")
}

func (r *OsqueryFileCategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Osquery file category", req, resp)
}
