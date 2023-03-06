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
	"github.com/zentralopensource/terraform-provider-zentral/internal/planmodifiers"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MonolithSubManifestResource{}
var _ resource.ResourceWithImportState = &MonolithSubManifestResource{}

func NewMonolithSubManifestResource() resource.Resource {
	return &MonolithSubManifestResource{}
}

// MonolithSubManifestResource defines the resource implementation.
type MonolithSubManifestResource struct {
	client *goztl.Client
}

func (r *MonolithSubManifestResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monolith_sub_manifest"
}

func (r *MonolithSubManifestResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Monolith sub manifests.",
		MarkdownDescription: "The resource `zentral_monolith_sub_manifest` manages Monolith sub manifests.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the sub manifest.",
				MarkdownDescription: "`ID` of the sub manifest.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the sub manifest.",
				MarkdownDescription: "Name of the sub manifest.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the sub manifest.",
				MarkdownDescription: "Description of the sub manifest.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.StringDefault(""),
				},
			},
			"meta_business_unit_id": schema.Int64Attribute{
				Description:         "The ID of the meta business unit this sub manifest is restricted to.",
				MarkdownDescription: "The `ID` of the meta business unit this sub manifest is restricted to.",
				Optional:            true,
			},
		},
	}
}

func (r *MonolithSubManifestResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MonolithSubManifestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data monolithSubManifest

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSM, _, err := r.client.MonolithSubManifests.Create(ctx, monolithSubManifestRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Monolith sub manifest, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Monolith sub manifest")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithSubManifestForState(ztlMSM))...)
}

func (r *MonolithSubManifestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data monolithSubManifest

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSM, _, err := r.client.MonolithSubManifests.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Monolith sub manifest %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Monolith sub manifest")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithSubManifestForState(ztlMSM))...)
}

func (r *MonolithSubManifestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data monolithSubManifest

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSM, _, err := r.client.MonolithSubManifests.Update(ctx, int(data.ID.ValueInt64()), monolithSubManifestRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Monolith sub manifest %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Monolith sub manifest")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithSubManifestForState(ztlMSM))...)
}

func (r *MonolithSubManifestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data monolithSubManifest

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MonolithSubManifests.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Monolith sub manifest %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Monolith sub manifest")
}

func (r *MonolithSubManifestResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Monolith sub manifest", req, resp)
}
