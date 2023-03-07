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
var _ resource.Resource = &MonolithManifestSubManifestResource{}
var _ resource.ResourceWithImportState = &MonolithManifestSubManifestResource{}

func NewMonolithManifestSubManifestResource() resource.Resource {
	return &MonolithManifestSubManifestResource{}
}

// MonolithManifestSubManifestResource defines the resource implementation.
type MonolithManifestSubManifestResource struct {
	client *goztl.Client
}

func (r *MonolithManifestSubManifestResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monolith_manifest_sub_manifest"
}

func (r *MonolithManifestSubManifestResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Monolith manifest sub manifests.",
		MarkdownDescription: "The resource `zentral_monolith_manifest_sub_manifest` manages Monolith manifest sub manifests.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the manifest sub manifest.",
				MarkdownDescription: "`ID` of the manifest sub manifest.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"manifest_id": schema.Int64Attribute{
				Description:         "ID of the manifest.",
				MarkdownDescription: "ID of the manifest.",
				Required:            true,
			},
			"sub_manifest_id": schema.Int64Attribute{
				Description:         "ID of the sub manifest.",
				MarkdownDescription: "ID of the sub manifest.",
				Required:            true,
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags used to scope the sub manifest.",
				MarkdownDescription: "The `ID`s of the tags used to scope the sub manifest.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *MonolithManifestSubManifestResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MonolithManifestSubManifestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data monolithManifestSubManifest

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMMSM, _, err := r.client.MonolithManifestSubManifests.Create(ctx, monolithManifestSubManifestRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Monolith manifest sub manifest, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an Monolith manifest sub manifest")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithManifestSubManifestForState(ztlMMSM))...)
}

func (r *MonolithManifestSubManifestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data monolithManifestSubManifest

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMMSM, _, err := r.client.MonolithManifestSubManifests.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Monolith manifest sub manifest %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read an Monolith manifest sub manifest")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithManifestSubManifestForState(ztlMMSM))...)
}

func (r *MonolithManifestSubManifestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data monolithManifestSubManifest

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMMSM, _, err := r.client.MonolithManifestSubManifests.Update(ctx, int(data.ID.ValueInt64()), monolithManifestSubManifestRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Monolith manifest sub manifest %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an Monolith manifest sub manifest")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithManifestSubManifestForState(ztlMMSM))...)
}

func (r *MonolithManifestSubManifestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data monolithManifestSubManifest

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MonolithManifestSubManifests.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Monolith manifest sub manifest %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an Monolith manifest sub manifest")
}

func (r *MonolithManifestSubManifestResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Monolith manifest sub manifest", req, resp)
}
