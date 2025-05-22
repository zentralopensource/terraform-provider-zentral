package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMDataAssetResource{}
var _ resource.ResourceWithImportState = &MDMDataAssetResource{}

func NewMDMDataAssetResource() resource.Resource {
	return &MDMDataAssetResource{}
}

// MDMDataAssetResource defines the resource implementation.
type MDMDataAssetResource struct {
	client *goztl.Client
}

func (r *MDMDataAssetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_data_asset"
}

func (r *MDMDataAssetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM data assets.",
		MarkdownDescription: "The resource `zentral_mdm_data_asset` manages MDM data assets.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the data asset.",
				MarkdownDescription: "`ID` of the data asset.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"artifact_id": schema.StringAttribute{
				Description:         "ID of the data asset artifact.",
				MarkdownDescription: "`ID` of the data asset artifact.",
				Required:            true,
			},
			"type": schema.StringAttribute{
				Description:         "The type of the data asset file.",
				MarkdownDescription: "The type of the data asset file.",
				Required:            true,
			},
			"file_uri": schema.StringAttribute{
				Description:         "The URI of the data asset file.",
				MarkdownDescription: "The URI of the data asset file.",
				Required:            true,
			},
			"file_sha256": schema.StringAttribute{
				Description:         "The hexadecimal digest of the sha256 hash of the file.",
				MarkdownDescription: "The hexadecimal digest of the sha256 hash of the file.",
				Required:            true,
			},
			"file_size": schema.Int64Attribute{
				Description:         "The size of the data asset file.",
				MarkdownDescription: "The size of the data asset file.",
				Computed:            true,
			},
			"filename": schema.StringAttribute{
				Description:         "The name of the data asset file.",
				MarkdownDescription: "The name of the data asset file.",
				Computed:            true,
			},
			"ios": schema.BoolAttribute{
				Description:         "Toggles the installation of the data asset on iOS devices.",
				MarkdownDescription: "Toggles the installation of the data asset on iOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ios_max_version": schema.StringAttribute{
				Description:         "Devices with this iOS version or higher will not receive this data asset.",
				MarkdownDescription: "Devices with this iOS version or higher will **not** receive this data asset.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ios_min_version": schema.StringAttribute{
				Description:         "Devices with this iOS version or higher will receive this data asset.",
				MarkdownDescription: "Devices with this iOS version or higher will receive this data asset.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ipados": schema.BoolAttribute{
				Description:         "Toggles the installation of the data asset on iPadOS devices.",
				MarkdownDescription: "Toggles the installation of the data asset on iPadOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ipados_max_version": schema.StringAttribute{
				Description:         "Devices with this iPadOS version or higher will not receive this data asset.",
				MarkdownDescription: "Devices with this iPadOS version or higher will **not** receive this data asset.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ipados_min_version": schema.StringAttribute{
				Description:         "Devices with this iPadOS version or higher will receive this data asset.",
				MarkdownDescription: "Devices with this iPadOS version or higher will receive this data asset.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"macos": schema.BoolAttribute{
				Description:         "Toggles the installation of the data asset on macOS devices.",
				MarkdownDescription: "Toggles the installation of the data asset on macOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"macos_max_version": schema.StringAttribute{
				Description:         "Devices with this macOS version or higher will not receive this data asset.",
				MarkdownDescription: "Devices with this macOS version or higher will **not** receive this data asset.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"macos_min_version": schema.StringAttribute{
				Description:         "Devices with this macOS version or higher will receive this data asset.",
				MarkdownDescription: "Devices with this macOS version or higher will receive this data asset.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tvos": schema.BoolAttribute{
				Description:         "Toggles the installation of the data asset on tvOS devices.",
				MarkdownDescription: "Toggles the installation of the data asset on tvOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"tvos_max_version": schema.StringAttribute{
				Description:         "Devices with this tvOS version or higher will not receive this data asset.",
				MarkdownDescription: "Devices with this tvOS version or higher will **not** receive this data asset.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tvos_min_version": schema.StringAttribute{
				Description:         "Devices with this tvOS version or higher will receive this data asset.",
				MarkdownDescription: "Devices with this tvOS version or higher will receive this data asset.",
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
				Description:         "Machines tagged with one of these tags will not receive the data asset.",
				MarkdownDescription: "Machines tagged with one of these tags will not receive the data asset.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.Int64Type, []attr.Value{})),
			},
			"tag_shards": schema.SetNestedAttribute{
				Description:         "A set of tag shard values different from the default shard, to determine if the tagged machines will receive the data asset.",
				MarkdownDescription: "A set of tag shard values different from the default shard, to determine if the tagged machines will receive the data asset.",
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
			"version": schema.Int64Attribute{
				Description:         "Version of the data asset.",
				MarkdownDescription: "Version of the data asset.",
				Required:            true,
			},
		},
	}
}

func (r *MDMDataAssetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMDataAssetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmDataAsset

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMDA, _, err := r.client.MDMDataAssets.Create(ctx, mdmDataAssetRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM data asset, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an MDM data asset")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmDataAssetForState(ztlMDA, data.FileURI))...)
}

func (r *MDMDataAssetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmDataAsset

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMDA, _, err := r.client.MDMDataAssets.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM data asset %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read an MDM data asset")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmDataAssetForState(ztlMDA, data.FileURI))...)
}

func (r *MDMDataAssetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmDataAsset

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMDA, _, err := r.client.MDMDataAssets.Update(ctx, data.ID.ValueString(), mdmDataAssetRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM data asset %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an MDM data asset")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmDataAssetForState(ztlMDA, data.FileURI))...)
}

func (r *MDMDataAssetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmDataAsset

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMDataAssets.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM data asset %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an MDM data asset")
}

func (r *MDMDataAssetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "MDM data asset", req, resp)
}
