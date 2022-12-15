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
var _ resource.Resource = &JMESPathCheckResource{}
var _ resource.ResourceWithImportState = &JMESPathCheckResource{}

func NewJMESPathCheckResource() resource.Resource {
	return &JMESPathCheckResource{}
}

// JMESPathCheckResource defines the resource implementation.
type JMESPathCheckResource struct {
	client *goztl.Client
}

func (r *JMESPathCheckResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jmespath_check"
}

func (r *JMESPathCheckResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages JMESPath compliance checks.",
		MarkdownDescription: "The resource `zentral_jmespath_check` manages JMESPath compliance checks.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the JMESPath compliance check.",
				MarkdownDescription: "`ID` of the JMESPath compliance check.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the JMESPath compliance check.",
				MarkdownDescription: "Name of the JMESPath compliance check.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the JMESPath compliance check.",
				MarkdownDescription: "Description of the JMESPath compliance check.",
				Optional:            true,
				Computed:            true,
			},
			"source_name": schema.StringAttribute{
				Description:         "The name of the inventory source the JMESPath compliance check is restricted to.",
				MarkdownDescription: "The name of the inventory source the JMESPath compliance check is restricted to.",
				Required:            true,
			},
			"platforms": schema.SetAttribute{
				Description:         "The platforms the JMESPath compliance check is restricted to.",
				MarkdownDescription: "The platforms the JMESPath compliance check is restricted to.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags the JMESPath compliance check is restricted to.",
				MarkdownDescription: "The IDs of the tags the JMESPath compliance check is restricted to.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"jmespath_expression": schema.StringAttribute{
				Description:         "The JMESPath compliance check expression.",
				MarkdownDescription: "The JMESPath compliance check expression.",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				Description:         "The JMESPath compliance check version.",
				MarkdownDescription: "The JMESPath compliance check version.",
				Computed:            true,
			},
		},
	}
}

func (r *JMESPathCheckResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *JMESPathCheckResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data jmespathCheck

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	platforms := make([]string, 0)
	for _, pv := range data.Platforms.Elements() { // nil if null or unknown → no iterations
		platforms = append(platforms, pv.(types.String).ValueString())
	}

	tagIDs := make([]int, 0)
	for _, tv := range data.TagIDs.Elements() { // nil if null or unknown → no iterations
		tagIDs = append(tagIDs, int(tv.(types.Int64).ValueInt64()))
	}

	ztlReq := &goztl.JMESPathCheckCreateRequest{
		Name:               data.Name.ValueString(),
		Description:        data.Description.ValueString(), // default to "" if null or unknown
		SourceName:         data.SourceName.ValueString(),
		Platforms:          platforms,
		TagIDs:             tagIDs,
		JMESPathExpression: data.JMESPathExpression.ValueString(),
	}
	ztlJC, _, err := r.client.JMESPathChecks.Create(ctx, ztlReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create JMESPath check, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a JMESPath check")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, jmespathCheckForState(ztlJC))...)
}

func (r *JMESPathCheckResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data jmespathCheck

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlJC, _, err := r.client.JMESPathChecks.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read JMESPath check %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a JMESPath check")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, jmespathCheckForState(ztlJC))...)
}

func (r *JMESPathCheckResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data jmespathCheck

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	platforms := make([]string, 0)
	for _, pv := range data.Platforms.Elements() { // nil if null or unknown → no iterations
		platforms = append(platforms, pv.(types.String).ValueString())
	}

	tagIDs := make([]int, 0)
	for _, tv := range data.TagIDs.Elements() { // nil if null or unknown → no iterations
		tagIDs = append(tagIDs, int(tv.(types.Int64).ValueInt64()))
	}

	ztlReq := &goztl.JMESPathCheckUpdateRequest{
		Name:               data.Name.ValueString(),
		Description:        data.Description.ValueString(), // default to "" if null or unknown
		SourceName:         data.SourceName.ValueString(),
		Platforms:          platforms,
		TagIDs:             tagIDs,
		JMESPathExpression: data.JMESPathExpression.ValueString(),
	}
	ztlJC, _, err := r.client.JMESPathChecks.Update(ctx, int(data.ID.ValueInt64()), ztlReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update JMESPath check %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a JMESPath check")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, jmespathCheckForState(ztlJC))...)
}

func (r *JMESPathCheckResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data jmespathCheck

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.JMESPathChecks.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete JMESPath check %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a JMESPath check")
}

func (r *JMESPathCheckResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "JMESPath check", req, resp)
}
