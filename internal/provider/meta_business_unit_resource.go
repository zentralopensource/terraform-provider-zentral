package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
	"github.com/zentralopensource/terraform-provider-zentral/internal/planmodifiers"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = metaBusinessUnitResourceType{}
var _ tfsdk.Resource = metaBusinessUnitResource{}
var _ tfsdk.ResourceWithImportState = metaBusinessUnitResource{}

type metaBusinessUnitResourceType struct{}

func (t metaBusinessUnitResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Manages meta business units.",
		MarkdownDescription: "The resource `zentral_meta_business_unit` manages meta business units.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description:         "ID of the meta business unit.",
				MarkdownDescription: "`ID` of the meta business unit.",
				Type:                types.Int64Type,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				Description:         "Name of the meta business unit.",
				MarkdownDescription: "Name of the meta business unit.",
				Type:                types.StringType,
				Required:            true,
			},
			"api_enrollment_enabled": {
				Description: "Enables API enrollments.",
				MarkdownDescription: "Enables API enrollments. Once enabled, it **CANNOT** be disabled. " +
					"Defaults to `true`.",
				Type:     types.BoolType,
				Optional: true,
				Computed: true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					planmodifiers.DefaultValue(types.Bool{Value: true}),
				},
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

type metaBusinessUnitResource struct {
	provider provider
}

func (r metaBusinessUnitResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data metaBusinessUnit

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	mbuCreateRequest := &goztl.MetaBusinessUnitCreateRequest{
		Name: data.Name.Value,
	}
	if data.APIEnrollmentEnabled.Value {
		mbuCreateRequest.APIEnrollmentEnabled = true
	}
	mbu, _, err := r.provider.client.MetaBusinessUnits.Create(ctx, mbuCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create meta business unit, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a meta business unit")

	diags = resp.State.Set(ctx, metaBusinessUnitForState(mbu))
	resp.Diagnostics.Append(diags...)
}

func (r metaBusinessUnitResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data metaBusinessUnit

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	mbu, _, err := r.provider.client.MetaBusinessUnits.GetByID(ctx, int(data.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read meta business unit %d, got error: %s", data.ID.Value, err),
		)
		return
	}

	tflog.Trace(ctx, "read a meta business unit")

	diags = resp.State.Set(ctx, metaBusinessUnitForState(mbu))
	resp.Diagnostics.Append(diags...)
}

func (r metaBusinessUnitResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data metaBusinessUnit

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	mbuUpdateRequest := &goztl.MetaBusinessUnitUpdateRequest{
		Name: data.Name.Value,
	}
	if data.APIEnrollmentEnabled.Value {
		mbuUpdateRequest.APIEnrollmentEnabled = true
	}
	mbu, _, err := r.provider.client.MetaBusinessUnits.Update(ctx, int(data.ID.Value), mbuUpdateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update meta business unit %d, got error: %s", data.ID.Value, err),
		)
		return
	}

	tflog.Trace(ctx, "updated a meta business unit")

	diags = resp.State.Set(ctx, metaBusinessUnitForState(mbu))
	resp.Diagnostics.Append(diags...)
}

func (r metaBusinessUnitResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data metaBusinessUnit

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.client.MetaBusinessUnits.Delete(ctx, int(data.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete meta business unit %d, got error: %s", data.ID.Value, err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a meta business unit")
}

func (r metaBusinessUnitResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "meta business unit", req, resp)
}
