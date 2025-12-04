package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &GWSGroupTagMappingResource{}
var _ resource.ResourceWithImportState = &GWSGroupTagMappingResource{}

func NewGWSGroupTagMappingResource() resource.Resource {
	return &GWSGroupTagMappingResource{}
}

// GWSGroupTagMappingResource defines the resource implementation.
type GWSGroupTagMappingResource struct {
	client *goztl.Client
}

func (r *GWSGroupTagMappingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gws_group_tag_mapping"
}

func (r *GWSGroupTagMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Google Workspace group tag mappings.",
		MarkdownDescription: "The resource `zentral_gws_group_tag_mapping` manages Google Workspace group tag mappings.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the group tag mapping (UUID).",
				MarkdownDescription: "`ID` of the group tag mapping (UUID).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_email": schema.StringAttribute{
				Description:         "Group email of the Google Workspace group for this group tag mapping.",
				MarkdownDescription: "Group email of the Google Workspace group for this group tag mapping.",
				Required:            true,
			},
			"connection_id": schema.StringAttribute{
				Description:         "ID of the Google Workspace connection for this group tag mapping (UUID).",
				MarkdownDescription: "`ID` of the Google Workspace connection for this group tag mapping (UUID).",
				Required:            true,
			},
			"tag_ids": schema.SetAttribute{
				Description:         "IDs of the tags mapped with this group tag mapping.",
				MarkdownDescription: "`ID`s of the tags mapped with this group tag mapping.",
				ElementType:         types.Int64Type,
				Required:            true,
			},
		},
	}
}

func (r *GWSGroupTagMappingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GWSGroupTagMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data gwsGroupTagMapping

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlGTM, _, err := r.client.GWSGroupTagMappings.Create(ctx, gwsGroupTagMappingRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Google Workspace group tag mapping, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Google Workspace group tag mapping")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, gwsGroupTagMappingForState(ztlGTM))...)
}

func (r *GWSGroupTagMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data gwsGroupTagMapping

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlGTM, _, err := r.client.GWSGroupTagMappings.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Google Workspace group tag mapping %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Google Workspace group tag mapping")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, gwsGroupTagMappingForState(ztlGTM))...)
}

func (r *GWSGroupTagMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data gwsGroupTagMapping

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlGTM, _, err := r.client.GWSGroupTagMappings.Update(ctx, data.ID.ValueString(), gwsGroupTagMappingRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Google Workspace group tag mapping %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Google Workspace group tag mapping")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, gwsGroupTagMappingForState(ztlGTM))...)
}

func (r *GWSGroupTagMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data gwsGroupTagMapping

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.GWSGroupTagMappings.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Google Workspace group tag mapping %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Google Workspace group tag mapping")
}

func (r *GWSGroupTagMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "Google Workspace group tag mapping", req, resp)
}
