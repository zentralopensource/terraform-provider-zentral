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
var _ resource.Resource = &OsqueryConfigurationPackResource{}
var _ resource.ResourceWithImportState = &OsqueryConfigurationPackResource{}

func NewOsqueryConfigurationPackResource() resource.Resource {
	return &OsqueryConfigurationPackResource{}
}

// OsqueryConfigurationPackResource defines the resource implementation.
type OsqueryConfigurationPackResource struct {
	client *goztl.Client
}

func (r *OsqueryConfigurationPackResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_osquery_configuration_pack"
}

func (r *OsqueryConfigurationPackResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Osquery configuration packs.",
		MarkdownDescription: "The resource `zentral_osquery_configuration_pack` manages Osquery configuration packs.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the configuration pack.",
				MarkdownDescription: "`ID` of the configuration pack.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"configuration_id": schema.Int64Attribute{
				Description:         "ID of the configuration.",
				MarkdownDescription: "ID of the configuration.",
				Required:            true,
			},
			"pack_id": schema.Int64Attribute{
				Description:         "ID of the pack.",
				MarkdownDescription: "ID of the pack.",
				Required:            true,
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags used to scope the pack.",
				MarkdownDescription: "The `ID`s of the tags used to scope the pack.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *OsqueryConfigurationPackResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OsqueryConfigurationPackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data osqueryConfigurationPack

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOCP, _, err := r.client.OsqueryConfigurationPacks.Create(ctx, osqueryConfigurationPackRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Osquery configuration pack, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an Osquery configuration pack")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryConfigurationPackForState(ztlOCP))...)
}

func (r *OsqueryConfigurationPackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data osqueryConfigurationPack

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOCP, _, err := r.client.OsqueryConfigurationPacks.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Osquery configuration pack %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read an Osquery configuration pack")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryConfigurationPackForState(ztlOCP))...)
}

func (r *OsqueryConfigurationPackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data osqueryConfigurationPack

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlOCP, _, err := r.client.OsqueryConfigurationPacks.Update(ctx, int(data.ID.ValueInt64()), osqueryConfigurationPackRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Osquery configuration pack %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an Osquery configuration pack")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, osqueryConfigurationPackForState(ztlOCP))...)
}

func (r *OsqueryConfigurationPackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data osqueryConfigurationPack

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.OsqueryConfigurationPacks.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Osquery configuration pack %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an Osquery configuration pack")
}

func (r *OsqueryConfigurationPackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Osquery configuration pack", req, resp)
}
