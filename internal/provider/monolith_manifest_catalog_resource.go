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
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MonolithManifestCatalogResource{}
var _ resource.ResourceWithImportState = &MonolithManifestCatalogResource{}

func NewMonolithManifestCatalogResource() resource.Resource {
	return &MonolithManifestCatalogResource{}
}

// MonolithManifestCatalogResource defines the resource implementation.
type MonolithManifestCatalogResource struct {
	client *goztl.Client
}

func (r *MonolithManifestCatalogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monolith_manifest_catalog"
}

func (r *MonolithManifestCatalogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Monolith manifest catalogs.",
		MarkdownDescription: "The resource `zentral_monolith_manifest_catalog` manages Monolith manifest catalogs.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the manifest catalog.",
				MarkdownDescription: "`ID` of the manifest catalog.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"manifest_id": schema.Int64Attribute{
				Description:         "ID of the manifest.",
				MarkdownDescription: "ID of the manifest.",
				Required:            true,
			},
			"catalog_id": schema.Int64Attribute{
				Description:         "ID of the catalog.",
				MarkdownDescription: "ID of the catalog.",
				Required:            true,
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags used to scope the catalog.",
				MarkdownDescription: "The `ID`s of the tags used to scope the catalog.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *MonolithManifestCatalogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MonolithManifestCatalogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data monolithManifestCatalog

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMMC, _, err := r.client.MonolithManifestCatalogs.Create(ctx, monolithManifestCatalogRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Monolith manifest catalog, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an Monolith manifest catalog")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithManifestCatalogForState(ztlMMC))...)
}

func (r *MonolithManifestCatalogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data monolithManifestCatalog

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMMC, _, err := r.client.MonolithManifestCatalogs.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Monolith manifest catalog %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read an Monolith manifest catalog")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithManifestCatalogForState(ztlMMC))...)
}

func (r *MonolithManifestCatalogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data monolithManifestCatalog

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMMC, _, err := r.client.MonolithManifestCatalogs.Update(ctx, int(data.ID.ValueInt64()), monolithManifestCatalogRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Monolith manifest catalog %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an Monolith manifest catalog")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithManifestCatalogForState(ztlMMC))...)
}

func (r *MonolithManifestCatalogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data monolithManifestCatalog

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MonolithManifestCatalogs.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Monolith manifest catalog %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an Monolith manifest catalog")
}

func (r *MonolithManifestCatalogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Monolith manifest catalog", req, resp)
}
