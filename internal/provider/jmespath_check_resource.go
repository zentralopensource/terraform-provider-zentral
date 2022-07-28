package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
	"github.com/zentralopensource/terraform-provider-zentral/internal/planmodifiers"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = jmespathCheckResourceType{}
var _ tfsdk.Resource = jmespathCheckResource{}
var _ tfsdk.ResourceWithImportState = jmespathCheckResource{}

type jmespathCheckResourceType struct{}

func (t jmespathCheckResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Manages JMESPath compliance checks.",
		MarkdownDescription: "The resource `zentral_jmespath_check` manages JMESPath compliance checks.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description:         "ID of the JMESPath compliance check.",
				MarkdownDescription: "`ID` of the JMESPath compliance check.",
				Type:                types.Int64Type,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				Description:         "Name of the JMESPath compliance check.",
				MarkdownDescription: "Name of the JMESPath compliance check.",
				Type:                types.StringType,
				Required:            true,
			},
			"description": {
				Description:         "Description of the JMESPath compliance check.",
				MarkdownDescription: "Description of the JMESPath compliance check.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					planmodifiers.DefaultValue(types.String{Value: ""}),
				},
			},
			"source_name": {
				Description:         "The name of the inventory source the JMESPath compliance check is restricted to.",
				MarkdownDescription: "The name of the inventory source the JMESPath compliance check is restricted to.",
				Type:                types.StringType,
				Required:            true,
			},
			"platforms": {
				Description:         "The platforms the JMESPath compliance check is restricted to.",
				MarkdownDescription: "The platforms the JMESPath compliance check is restricted to.",
				Type:                types.SetType{ElemType: types.StringType},
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					planmodifiers.DefaultValue(types.Set{ElemType: types.StringType, Elems: []attr.Value{}}),
				},
			},
			"tag_ids": {
				Description:         "The IDs of the tags the JMESPath compliance check is restricted to.",
				MarkdownDescription: "The IDs of the tags the JMESPath compliance check is restricted to.",
				Type:                types.SetType{ElemType: types.Int64Type},
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					planmodifiers.DefaultValue(types.Set{ElemType: types.Int64Type, Elems: []attr.Value{}}),
				},
			},
			"jmespath_expression": {
				Description:         "The JMESPath compliance check expression.",
				MarkdownDescription: "The JMESPath compliance check expression.",
				Type:                types.StringType,
				Required:            true,
			},
			"version": {
				Description:         "The JMESPath compliance check version.",
				MarkdownDescription: "The JMESPath compliance check version.",
				Type:                types.Int64Type,
				Computed:            true,
			},
		},
	}, nil
}

func (t jmespathCheckResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return jmespathCheckResource{
		provider: provider,
	}, diags
}

type jmespathCheckResource struct {
	provider provider
}

func (r jmespathCheckResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data jmespathCheck

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	platforms := make([]string, 0)
	for _, pv := range data.Platforms.Elems {
		platforms = append(platforms, pv.(types.String).Value)
	}

	tagIDs := make([]int, 0)
	for _, tv := range data.TagIDs.Elems {
		tagIDs = append(tagIDs, int(tv.(types.Int64).Value))
	}

	ztlReq := &goztl.JMESPathCheckCreateRequest{
		Name:               data.Name.Value,
		Description:        data.Description.Value,
		SourceName:         data.SourceName.Value,
		Platforms:          platforms,
		TagIDs:             tagIDs,
		JMESPathExpression: data.JMESPathExpression.Value,
	}
	ztlJC, _, err := r.provider.client.JMESPathChecks.Create(ctx, ztlReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create JMESPath check, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a JMESPath check")
	newC := jmespathCheckForState(ztlJC)
	tflog.Error(ctx, goztl.Stringify(newC))

	diags = resp.State.Set(ctx, newC)
	resp.Diagnostics.Append(diags...)
}

func (r jmespathCheckResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data jmespathCheck

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlJC, _, err := r.provider.client.JMESPathChecks.GetByID(ctx, int(data.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read JMESPath check %d, got error: %s", data.ID.Value, err),
		)
		return
	}

	tflog.Trace(ctx, "read a JMESPath check")

	diags = resp.State.Set(ctx, jmespathCheckForState(ztlJC))
	resp.Diagnostics.Append(diags...)
}

func (r jmespathCheckResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data jmespathCheck

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	platforms := make([]string, 0)
	for _, pv := range data.Platforms.Elems {
		platforms = append(platforms, pv.(types.String).Value)
	}

	tagIDs := make([]int, 0)
	for _, tv := range data.TagIDs.Elems {
		tagIDs = append(tagIDs, int(tv.(types.Int64).Value))
	}

	ztlReq := &goztl.JMESPathCheckUpdateRequest{
		Name:               data.Name.Value,
		Description:        data.Description.Value,
		SourceName:         data.SourceName.Value,
		Platforms:          platforms,
		TagIDs:             tagIDs,
		JMESPathExpression: data.JMESPathExpression.Value,
	}
	ztlJC, _, err := r.provider.client.JMESPathChecks.Update(ctx, int(data.ID.Value), ztlReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update JMESPath check %d, got error: %s", data.ID.Value, err),
		)
		return
	}

	tflog.Trace(ctx, "updated a JMESPath check")

	diags = resp.State.Set(ctx, jmespathCheckForState(ztlJC))
	resp.Diagnostics.Append(diags...)
}

func (r jmespathCheckResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data jmespathCheck

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.client.JMESPathChecks.Delete(ctx, int(data.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete JMESPath check %d, got error: %s", data.ID.Value, err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a JMESPath check")
}

func (r jmespathCheckResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "JMESPath check", req, resp)
}
