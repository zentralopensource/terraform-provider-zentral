package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MonolithCatalogResource{}
var _ resource.ResourceWithImportState = &MonolithCatalogResource{}

func NewMonolithCatalogResource() resource.Resource {
	return &MonolithCatalogResource{}
}

// MonolithCatalogResource defines the resource implementation.
type MonolithCatalogResource struct {
	client *goztl.Client
}

func (r *MonolithCatalogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monolith_catalog"
}

func (r *MonolithCatalogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Monolith catalogs.",
		MarkdownDescription: "The resource `zentral_monolith_catalog` manages Monolith catalogs.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the catalog.",
				MarkdownDescription: "`ID` of the catalog.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the catalog.",
				MarkdownDescription: "Name of the catalog.",
				Required:            true,
			},
			"priority": schema.Int64Attribute{
				Description:         "Priority of the catalog.",
				MarkdownDescription: "Priority of the catalog.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
		},
	}
}

func (r *MonolithCatalogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MonolithCatalogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data monolithCatalog

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMC, _, err := r.client.MonolithCatalogs.Create(ctx, monolithCatalogRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Monolith catalog, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Monolith catalog")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithCatalogForState(ztlMC))...)
}

func (r *MonolithCatalogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data monolithCatalog

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMC, _, err := r.client.MonolithCatalogs.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Monolith catalog %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Monolith catalog")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithCatalogForState(ztlMC))...)
}

func (r *MonolithCatalogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data monolithCatalog

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMC, _, err := r.client.MonolithCatalogs.Update(ctx, int(data.ID.ValueInt64()), monolithCatalogRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Monolith catalog %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Monolith catalog")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithCatalogForState(ztlMC))...)
}

func (r *MonolithCatalogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data monolithCatalog

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MonolithCatalogs.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Monolith catalog %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Monolith catalog")
}

func (r *MonolithCatalogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Monolith catalog", req, resp)
}
