package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = metaBusinessUnitResourceType{}
var _ tfsdk.Resource = metaBusinessUnitResource{}
var _ tfsdk.ResourceWithImportState = metaBusinessUnitResource{}

type metaBusinessUnitResourceType struct{}

func (t metaBusinessUnitResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Meta business unit",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "Name of the meta business unit",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "ID of the meta business unit",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.Int64Type,
			},
		},
	}, nil
}

func (t metaBusinessUnitResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return metaBusinessUnitResource{
		provider: provider,
	}, diags
}

type metaBusinessUnitResourceData struct {
	Name types.String `tfsdk:"name"`
	Id   types.Int64  `tfsdk:"id"`
}

type metaBusinessUnitResource struct {
	provider provider
}

func (r metaBusinessUnitResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data metaBusinessUnitResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	mbuCreateRequest := &goztl.MetaBusinessUnitCreateRequest{
		Name: data.Name.Value,
	}
	mbu, _, err := r.provider.client.MetaBusinessUnits.Create(ctx, mbuCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create meta business unit, got error: %s", err))
		return
	}

	data.Id = types.Int64{Value: int64(mbu.ID)}

	tflog.Trace(ctx, "created a meta business unit")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r metaBusinessUnitResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data metaBusinessUnitResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	mbu, _, err := r.provider.client.MetaBusinessUnits.GetByID(ctx, int(data.Id.Value))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read meta business unit %d, got error: %s", data.Id.Value, err),
		)
		return
	}

	data.Name = types.String{Value: mbu.Name}

	tflog.Trace(ctx, "read a meta business unit")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r metaBusinessUnitResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data metaBusinessUnitResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	mbuUpdateRequest := &goztl.MetaBusinessUnitUpdateRequest{
		Name: data.Name.Value,
	}
	_, _, err := r.provider.client.MetaBusinessUnits.Update(ctx, int(data.Id.Value), mbuUpdateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update meta business unit %d, got error: %s", data.Id.Value, err),
		)
		return
	}

	tflog.Trace(ctx, "updated a meta business unit")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r metaBusinessUnitResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data metaBusinessUnitResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.client.MetaBusinessUnits.Delete(ctx, int(data.Id.Value))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete meta business unit %d, got error: %s", data.Id.Value, err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a meta business unit")
}

func (r metaBusinessUnitResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
