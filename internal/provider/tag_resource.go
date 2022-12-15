package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &TagResource{}
var _ resource.ResourceWithImportState = &TagResource{}

func NewTagResource() resource.Resource {
	return &TagResource{}
}

// TagResource defines the resource implementation.
type TagResource struct {
	client *goztl.Client
}

func (r *TagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (r *TagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages tags.",
		MarkdownDescription: "The resource `zentral_tag` manages tags.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the tag.",
				MarkdownDescription: "`ID` of the tag.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"taxonomy_id": schema.Int64Attribute{
				Description:         "ID of the tag taxonomy.",
				MarkdownDescription: "`ID` of the tag taxonomy.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the tag.",
				MarkdownDescription: "Name of the tag.",
				Required:            true,
			},
			"color": schema.StringAttribute{
				Description:         "Color of the tag.",
				MarkdownDescription: "Color of the tag.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *TagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data tag

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tagCreateRequest := &goztl.TagCreateRequest{
		Name:  data.Name.ValueString(),
		Color: data.Color.ValueString(),
	}
	if !data.TaxonomyID.IsNull() {
		tagCreateRequest.TaxonomyID = goztl.Int(int(data.TaxonomyID.ValueInt64()))
	}
	tag, _, err := r.client.Tags.Create(ctx, tagCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create tag, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a tag")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, tagForState(tag))...)
}

func (r *TagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data tag

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tag, _, err := r.client.Tags.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read tag %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a tag")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, tagForState(tag))...)
}

func (r *TagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data tag

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tagUpdateRequest := &goztl.TagUpdateRequest{
		Name:  data.Name.ValueString(),
		Color: data.Color.ValueString(),
	}
	if !data.TaxonomyID.IsNull() {
		tagUpdateRequest.TaxonomyID = goztl.Int(int(data.TaxonomyID.ValueInt64()))
	}
	tag, _, err := r.client.Tags.Update(ctx, int(data.ID.ValueInt64()), tagUpdateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update tag %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a tag")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, tagForState(tag))...)
}

func (r *TagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data tag

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Tags.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete tag %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a tag")
}

func (r *TagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "tag", req, resp)
}
