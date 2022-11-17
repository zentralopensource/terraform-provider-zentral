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
var _ provider.ResourceType = tagResourceType{}
var _ resource.Resource = tagResource{}
var _ resource.ResourceWithImportState = tagResource{}

type tagResourceType struct{}

func (t tagResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Manages tags.",
		MarkdownDescription: "The resource `zentral_tag` manages tags.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description:         "ID of the tag.",
				MarkdownDescription: "`ID` of the tag.",
				Type:                types.Int64Type,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"taxonomy_id": {
				Description:         "ID of the tag taxonomy.",
				MarkdownDescription: "`ID` of the tag taxonomy.",
				Type:                types.Int64Type,
				Optional:            true,
			},
			"name": {
				Description:         "Name of the tag.",
				MarkdownDescription: "Name of the tag.",
				Type:                types.StringType,
				Required:            true,
			},
			"color": {
				Description:         "Color of the tag.",
				MarkdownDescription: "Color of the tag.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
		},
	}, nil
}

func (t tagResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return tagResource{
		provider: provider,
	}, diags
}

type tagResource struct {
	provider zentralProvider
}

func (r tagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data tag

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tagCreateRequest := &goztl.TagCreateRequest{
		Name:  data.Name.Value,
		Color: data.Color.Value,
	}
	if !data.TaxonomyID.Null {
		tagCreateRequest.TaxonomyID = goztl.Int(int(data.TaxonomyID.Value))
	}
	tag, _, err := r.provider.client.Tags.Create(ctx, tagCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create tag, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a tag")

	diags = resp.State.Set(ctx, tagForState(tag))
	resp.Diagnostics.Append(diags...)
}

func (r tagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data tag

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tag, _, err := r.provider.client.Tags.GetByID(ctx, int(data.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read tag %d, got error: %s", data.ID.Value, err),
		)
		return
	}

	tflog.Trace(ctx, "read a tag")

	diags = resp.State.Set(ctx, tagForState(tag))
	resp.Diagnostics.Append(diags...)
}

func (r tagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data tag

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tagUpdateRequest := &goztl.TagUpdateRequest{
		Name:  data.Name.Value,
		Color: data.Color.Value,
	}
	if !data.TaxonomyID.Null {
		tagUpdateRequest.TaxonomyID = goztl.Int(int(data.TaxonomyID.Value))
	}
	tag, _, err := r.provider.client.Tags.Update(ctx, int(data.ID.Value), tagUpdateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update tag %d, got error: %s", data.ID.Value, err),
		)
		return
	}

	tflog.Trace(ctx, "updated a tag")

	diags = resp.State.Set(ctx, tagForState(tag))
	resp.Diagnostics.Append(diags...)
}

func (r tagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data tag

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.client.Tags.Delete(ctx, int(data.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete tag %d, got error: %s", data.ID.Value, err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a tag")
}

func (r tagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "tag", req, resp)
}
