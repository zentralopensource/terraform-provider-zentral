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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMBlueprintArtifactResource{}
var _ resource.ResourceWithImportState = &MDMBlueprintArtifactResource{}

func NewMDMBlueprintArtifactResource() resource.Resource {
	return &MDMBlueprintArtifactResource{}
}

// MDMBlueprintArtifactResource defines the resource implementation.
type MDMBlueprintArtifactResource struct {
	client *goztl.Client
}

func (r *MDMBlueprintArtifactResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_blueprint_artifact"
}

func (r *MDMBlueprintArtifactResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM blueprint artifacts.",
		MarkdownDescription: "The resource `zentral_mdm_blueprint_artifact` manages MDM blueprint artifacts.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the blueprint artifact.",
				MarkdownDescription: "`ID` of the blueprint artifact.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"blueprint_id": schema.Int64Attribute{
				Description:         "ID of the blueprint artifact blueprint.",
				MarkdownDescription: "`ID` of the blueprint artifact blueprint.",
				Required:            true,
			},
			"artifact_id": schema.StringAttribute{
				Description:         "ID of the blueprint artifact artifact.",
				MarkdownDescription: "`ID` of the blueprint artifact artifact.",
				Required:            true,
			},
			"ios": schema.BoolAttribute{
				Description:         "Toggles the installation of the blueprint artifact on iOS devices.",
				MarkdownDescription: "Toggles the installation of the blueprint artifact on iOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ios_max_version": schema.StringAttribute{
				Description:         "Devices with this iOS version or higher will not receive this blueprint artifact.",
				MarkdownDescription: "Devices with this iOS version or higher will **not** receive this blueprint artifact.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ios_min_version": schema.StringAttribute{
				Description:         "Devices with this iOS version or higher will receive this blueprint artifact.",
				MarkdownDescription: "Devices with this iOS version or higher will receive this blueprint artifact.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ipados": schema.BoolAttribute{
				Description:         "Toggles the installation of the blueprint artifact on iPadOS devices.",
				MarkdownDescription: "Toggles the installation of the blueprint artifact on iPadOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ipados_max_version": schema.StringAttribute{
				Description:         "Devices with this iPadOS version or higher will not receive this blueprint artifact.",
				MarkdownDescription: "Devices with this iPadOS version or higher will **not** receive this blueprint artifact.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ipados_min_version": schema.StringAttribute{
				Description:         "Devices with this iPadOS version or higher will receive this blueprint artifact.",
				MarkdownDescription: "Devices with this iPadOS version or higher will receive this blueprint artifact.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"macos": schema.BoolAttribute{
				Description:         "Toggles the installation of the blueprint artifact on macOS devices.",
				MarkdownDescription: "Toggles the installation of the blueprint artifact on macOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"macos_max_version": schema.StringAttribute{
				Description:         "Devices with this macOS version or higher will not receive this blueprint artifact.",
				MarkdownDescription: "Devices with this macOS version or higher will **not** receive this blueprint artifact.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"macos_min_version": schema.StringAttribute{
				Description:         "Devices with this macOS version or higher will receive this blueprint artifact.",
				MarkdownDescription: "Devices with this macOS version or higher will receive this blueprint artifact.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tvos": schema.BoolAttribute{
				Description:         "Toggles the installation of the blueprint artifact on tvOS devices.",
				MarkdownDescription: "Toggles the installation of the blueprint artifact on tvOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"tvos_max_version": schema.StringAttribute{
				Description:         "Devices with this tvOS version or higher will not receive this blueprint artifact.",
				MarkdownDescription: "Devices with this tvOS version or higher will **not** receive this blueprint artifact.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tvos_min_version": schema.StringAttribute{
				Description:         "Devices with this tvOS version or higher will receive this blueprint artifact.",
				MarkdownDescription: "Devices with this tvOS version or higher will receive this blueprint artifact.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
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
				Description:         "Machines tagged with one of these tags will not receive the blueprint artifact.",
				MarkdownDescription: "Machines tagged with one of these tags will not receive the blueprint artifact.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.Int64Type, []attr.Value{})),
			},
			"tag_shards": schema.SetNestedAttribute{
				Description:         "A set of tag shard values different from the default shard, to determine if the tagged machines will receive the blueprint artifact.",
				MarkdownDescription: "A set of tag shard values different from the default shard, to determine if the tagged machines will receive the blueprint artifact.",
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

func (r *MDMBlueprintArtifactResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMBlueprintArtifactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmBlueprintArtifact

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMBA, _, err := r.client.MDMBlueprintArtifacts.Create(ctx, mdmBlueprintArtifactRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM blueprint artifact, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a MDM blueprint artifact")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmBlueprintArtifactForState(ztlMBA))...)
}

func (r *MDMBlueprintArtifactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmBlueprintArtifact

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMBA, _, err := r.client.MDMBlueprintArtifacts.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM blueprint artifact %d, got error: %s", int(data.ID.ValueInt64()), err),
		)
		return
	}

	tflog.Trace(ctx, "read a MDM blueprint artifact")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmBlueprintArtifactForState(ztlMBA))...)
}

func (r *MDMBlueprintArtifactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmBlueprintArtifact

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMBA, _, err := r.client.MDMBlueprintArtifacts.Update(ctx, int(data.ID.ValueInt64()), mdmBlueprintArtifactRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM blueprint artifact %d, got error: %s", int(data.ID.ValueInt64()), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a MDM blueprint artifact")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmBlueprintArtifactForState(ztlMBA))...)
}

func (r *MDMBlueprintArtifactResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmBlueprintArtifact

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMBlueprintArtifacts.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM blueprint artifact %d, got error: %s", int(data.ID.ValueInt64()), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a MDM blueprint artifact")
}

func (r *MDMBlueprintArtifactResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "MDM blueprint artifact", req, resp)
}
