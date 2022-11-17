package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = taxonomyResourceType{}
var _ resource.Resource = taxonomyResource{}
var _ resource.ResourceWithImportState = taxonomyResource{}

type taxonomyResourceType struct{}

func (t taxonomyResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Manages taxonomies.",
		MarkdownDescription: "The resource `zentral_taxonomy` manages taxonomies.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description:         "ID of the taxonomy.",
				MarkdownDescription: "`ID` of the taxonomy.",
				Type:                types.Int64Type,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"name": {
				Description:         "Name of the taxonomy.",
				MarkdownDescription: "Name of the taxonomy.",
				Type:                types.StringType,
				Required:            true,
			},
		},
	}, nil
}

func (t taxonomyResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return taxonomyResource{
		provider: provider,
	}, diags
}

type taxonomyResource struct {
	provider zentralProvider
}

func (r taxonomyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data taxonomy

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	taxonomyCreateRequest := &goztl.TaxonomyCreateRequest{
		Name: data.Name.Value,
	}
	taxonomy, _, err := r.provider.client.Taxonomies.Create(ctx, taxonomyCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create taxonomy, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a taxonomy")

	diags = resp.State.Set(ctx, taxonomyForState(taxonomy))
	resp.Diagnostics.Append(diags...)
}

func (r taxonomyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data taxonomy

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	taxonomy, _, err := r.provider.client.Taxonomies.GetByID(ctx, int(data.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read taxonomy %d, got error: %s", data.ID.Value, err),
		)
		return
	}

	tflog.Trace(ctx, "read a taxonomy")

	diags = resp.State.Set(ctx, taxonomyForState(taxonomy))
	resp.Diagnostics.Append(diags...)
}

func (r taxonomyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data taxonomy

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	taxonomyUpdateRequest := &goztl.TaxonomyUpdateRequest{
		Name: data.Name.Value,
	}
	taxonomy, _, err := r.provider.client.Taxonomies.Update(ctx, int(data.ID.Value), taxonomyUpdateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update taxonomy %d, got error: %s", data.ID.Value, err),
		)
		return
	}

	tflog.Trace(ctx, "updated a taxonomy")

	diags = resp.State.Set(ctx, taxonomyForState(taxonomy))
	resp.Diagnostics.Append(diags...)
}

func (r taxonomyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data taxonomy

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.client.Taxonomies.Delete(ctx, int(data.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete taxonomy %d, got error: %s", data.ID.Value, err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a taxonomy")
}

func (r taxonomyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "taxonomy", req, resp)
}
