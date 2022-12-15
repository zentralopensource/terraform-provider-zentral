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
var _ resource.Resource = &TaxonomyResource{}
var _ resource.ResourceWithImportState = &TaxonomyResource{}

func NewTaxonomyResource() resource.Resource {
	return &TaxonomyResource{}
}

// TaxonomyResource defines the resource implementation.
type TaxonomyResource struct {
	client *goztl.Client
}

func (r *TaxonomyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_taxonomy"
}

func (t *TaxonomyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages taxonomies.",
		MarkdownDescription: "The resource `zentral_taxonomy` manages taxonomies.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the taxonomy.",
				MarkdownDescription: "`ID` of the taxonomy.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the taxonomy.",
				MarkdownDescription: "Name of the taxonomy.",
				Required:            true,
			},
		},
	}
}

func (r *TaxonomyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TaxonomyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data taxonomy

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	taxonomyCreateRequest := &goztl.TaxonomyCreateRequest{
		Name: data.Name.ValueString(),
	}
	taxonomy, _, err := r.client.Taxonomies.Create(ctx, taxonomyCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create taxonomy, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a taxonomy")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, taxonomyForState(taxonomy))...)
}

func (r *TaxonomyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data taxonomy

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	taxonomy, _, err := r.client.Taxonomies.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read taxonomy %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a taxonomy")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, taxonomyForState(taxonomy))...)
}

func (r *TaxonomyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data taxonomy

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	taxonomyUpdateRequest := &goztl.TaxonomyUpdateRequest{
		Name: data.Name.ValueString(),
	}
	taxonomy, _, err := r.client.Taxonomies.Update(ctx, int(data.ID.ValueInt64()), taxonomyUpdateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update taxonomy %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a taxonomy")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, taxonomyForState(taxonomy))...)
}

func (r *TaxonomyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data taxonomy

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Taxonomies.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete taxonomy %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a taxonomy")
}

func (r *TaxonomyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "taxonomy", req, resp)
}
