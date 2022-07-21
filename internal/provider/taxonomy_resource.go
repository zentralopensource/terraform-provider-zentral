package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = taxonomyResourceType{}
var _ tfsdk.Resource = taxonomyResource{}
var _ tfsdk.ResourceWithImportState = taxonomyResource{}

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
					tfsdk.UseStateForUnknown(),
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

func (t taxonomyResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return taxonomyResource{
		provider: provider,
	}, diags
}

type taxonomyResource struct {
	provider provider
}

func (r taxonomyResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
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

func (r taxonomyResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
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

func (r taxonomyResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
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

func (r taxonomyResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
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

func (r taxonomyResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "taxonomy", req, resp)
}
