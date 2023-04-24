package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MonolithSubManifestPkgInfoResource{}
var _ resource.ResourceWithImportState = &MonolithSubManifestPkgInfoResource{}

func NewMonolithSubManifestPkgInfoResource() resource.Resource {
	return &MonolithSubManifestPkgInfoResource{}
}

// MonolithSubManifestPkgInfoResource defines the resource implementation.
type MonolithSubManifestPkgInfoResource struct {
	client *goztl.Client
}

func (r *MonolithSubManifestPkgInfoResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monolith_sub_manifest_pkg_info"
}

func (r *MonolithSubManifestPkgInfoResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Monolith sub manifest pkg infos.",
		MarkdownDescription: "The resource `zentral_monolith_sub_manifest_pkg_info` manages Monolith sub manifest pkg infos.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the sub manifest pkg info.",
				MarkdownDescription: "`ID` of the sub manifest pkg info.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"sub_manifest_id": schema.Int64Attribute{
				Description:         "The ID of the sub manifest this pkg info is attached to.",
				MarkdownDescription: "The `ID` of the sub manifest this pkg info is attached to.",
				Required:            true,
			},
			"key": schema.StringAttribute{
				Description:         "Sub manifest key under which this pkg info will be included.",
				MarkdownDescription: "Sub manifest key under which this pkg info will be included.",
				Required:            true,
			},
			"pkg_info_name": schema.StringAttribute{
				Description:         "The name of the pkg info to include.",
				MarkdownDescription: "The name of the pkg info to include.",
				Required:            true,
			},
			"featured_item": schema.BoolAttribute{
				Description:         "If true, this pkg info will be displayed in the featured items section in Managed Software Center. Defaults to false.",
				MarkdownDescription: "If `true`, this pkg info will be displayed in the featured items section in Managed Software Center. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"condition_id": schema.Int64Attribute{
				Description:         "The ID of the condition that is evaluated to decide if this pkg info is included.",
				MarkdownDescription: "The `ID` of the condition that is evaluated to decide if this pkg info is included.",
				Optional:            true,
			},
			"shard_modulo": schema.Int64Attribute{
				Description:         "The modulo used to calculate the shards. Defaults to 100.",
				MarkdownDescription: "The modulo used to calculate the shards. Defaults to `100`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(100),
			},
			"default_shard": schema.Int64Attribute{
				Description:         "The default shard value. Defaults to 100.",
				MarkdownDescription: "The default shard value. Defaults to `100`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(100),
			},
			"excluded_tag_ids": schema.SetAttribute{
				Description:         "Machines tagged with one of these tags will not receive the pkg info.",
				MarkdownDescription: "Machines tagged with one of these tags will not receive the pkg info.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.Int64Type, []attr.Value{})),
			},
			"tag_shards": schema.SetNestedAttribute{
				Description:         "A set of tag shard values different from the default shard, to determine if the tagged machines will receive the pkg info.",
				MarkdownDescription: "A set of tag shard values different from the default shard, to determine if the tagged machines will receive the pkg info.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag_id": schema.Int64Attribute{
							Description:         "The ID of the tag.",
							MarkdownDescription: "The `ID` of the tag.",
							Required:            true,
						},
						"shard": schema.Int64Attribute{
							Description:         "The shard for the tag.",
							MarkdownDescription: "The shard for the tag.",
							Required:            true,
						},
					},
				},
				Optional: true,
				Computed: true,
				Default:  setdefault.StaticValue(types.SetValueMust(types.ObjectType{AttrTypes: tagShardAttrTypes}, []attr.Value{})),
			},
		},
	}
}

func (r *MonolithSubManifestPkgInfoResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MonolithSubManifestPkgInfoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data monolithSubManifestPkgInfo

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSMPI, _, err := r.client.MonolithSubManifestPkgInfos.Create(ctx, monolithSubManifestPkgInfoRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Monolith sub manifest pkg info, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Monolith sub manifest pkg info")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithSubManifestPkgInfoForState(ztlMSMPI))...)
}

func (r *MonolithSubManifestPkgInfoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data monolithSubManifestPkgInfo

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSMPI, _, err := r.client.MonolithSubManifestPkgInfos.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Monolith sub manifest pkg info %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Monolith sub manifest pkg info")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithSubManifestPkgInfoForState(ztlMSMPI))...)
}

func (r *MonolithSubManifestPkgInfoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data monolithSubManifestPkgInfo

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSMPI, _, err := r.client.MonolithSubManifestPkgInfos.Update(ctx, int(data.ID.ValueInt64()), monolithSubManifestPkgInfoRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Monolith sub manifest pkg info %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Monolith sub manifest pkg info")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithSubManifestPkgInfoForState(ztlMSMPI))...)
}

func (r *MonolithSubManifestPkgInfoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data monolithSubManifestPkgInfo

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MonolithSubManifestPkgInfos.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Monolith sub manifest pkg info %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Monolith sub manifest pkg info")
}

func (r *MonolithSubManifestPkgInfoResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Monolith sub manifest pkg info", req, resp)
}
