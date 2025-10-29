package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMStoreAppResource{}
var _ resource.ResourceWithImportState = &MDMStoreAppResource{}

func NewMDMStoreAppResource() resource.Resource {
	return &MDMStoreAppResource{}
}

// MDMStoreAppResource defines the resource implementation.
type MDMStoreAppResource struct {
	client *goztl.Client
}

func (r *MDMStoreAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_store_app"
}

func (r *MDMStoreAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM store apps.",
		MarkdownDescription: "The resource `zentral_mdm_store_app` manages MDM store apps.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the app.",
				MarkdownDescription: "`ID` of the app.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"artifact_id": schema.StringAttribute{
				Description:         "ID of the app artifact.",
				MarkdownDescription: "`ID` of the app artifact.",
				Required:            true,
			},
			"location_asset_id": schema.Int64Attribute{
				Description:         "The location asset ID.",
				MarkdownDescription: "The location asset `ID`.",
				Required:            true,
			},
			"associated_domains": schema.ListAttribute{
				Description:         "An array that contains the associated domains to add to this app. Available in iOS 13 and later.",
				MarkdownDescription: "An array that contains the associated domains to add to this app. Available in iOS 13 and later.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"associated_domains_enable_direct_downloads": schema.BoolAttribute{
				Description:         "If true, perform claimed site association verification directly at the domain instead of on Apple's servers. Only set this to true for domains that can't access the internet. Available in iOS 14 and later. Defaults to false.",
				MarkdownDescription: "If `true`, perform claimed site association verification directly at the domain instead of on Apple's servers. Only set this to `true` for domains that can't access the internet. Available in iOS 14 and later. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"configuration": schema.StringAttribute{
				Description:         "A dictionary serialized as plist that contains the initial configuration of the app.",
				MarkdownDescription: "A dictionary serialized as plist that contains the initial configuration of the app.",
				Optional:            true,
			},
			"content_filter_uuid": schema.StringAttribute{
				Description:         "The content filter UUID for this app. Available in iOS 16 and later.",
				MarkdownDescription: "The content filter UUID for this app. Available in iOS 16 and later.",
				Optional:            true,
			},
			"dns_proxy_uuid": schema.StringAttribute{
				Description:         "The DNS proxy UUID for this app. Available in iOS 16 and later.",
				MarkdownDescription: "The DNS proxy UUID for this app. Available in iOS 16 and later..",
				Optional:            true,
			},
			"vpn_uuid": schema.StringAttribute{
				Description:         "A per-app VPN unique identifier for this app. Available in iOS 7 and later.",
				MarkdownDescription: "A per-app VPN unique identifier for this app. Available in iOS 7 and later.",
				Optional:            true,
			},
			"prevent_backup": schema.BoolAttribute{
				Description:         "If true, prevent backup of app data. Defaults to false.",
				MarkdownDescription: "If `true`, prevent backup of app data. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"removable": schema.BoolAttribute{
				Description:         "If false, this app isn't removable while it's a managed app. Available in iOS 14 and later, and tvOS 14 and later. Defaults to false.",
				MarkdownDescription: "If `false`, this app isn't removable while it's a managed app. Available in iOS 14 and later, and tvOS 14 and later. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"remove_on_unenroll": schema.BoolAttribute{
				Description:         "If true, the app will be removed when the device is unenrolled. Defaults to true.",
				MarkdownDescription: "If `true`, the app will be removed when the device is unenrolled. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"ios": schema.BoolAttribute{
				Description:         "Toggles the installation of the app on iOS devices.",
				MarkdownDescription: "Toggles the installation of the app on iOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ios_max_version": schema.StringAttribute{
				Description:         "Devices with this iOS version or higher will not receive this app.",
				MarkdownDescription: "Devices with this iOS version or higher will **not** receive this app.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ios_min_version": schema.StringAttribute{
				Description:         "Devices with this iOS version or higher will receive this app.",
				MarkdownDescription: "Devices with this iOS version or higher will receive this app.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ipados": schema.BoolAttribute{
				Description:         "Toggles the installation of the app on iPadOS devices.",
				MarkdownDescription: "Toggles the installation of the app on iPadOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ipados_max_version": schema.StringAttribute{
				Description:         "Devices with this iPadOS version or higher will not receive this app.",
				MarkdownDescription: "Devices with this iPadOS version or higher will **not** receive this app.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ipados_min_version": schema.StringAttribute{
				Description:         "Devices with this iPadOS version or higher will receive this app.",
				MarkdownDescription: "Devices with this iPadOS version or higher will receive this app.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"macos": schema.BoolAttribute{
				Description:         "Toggles the installation of the app on macOS devices.",
				MarkdownDescription: "Toggles the installation of the app on macOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"macos_max_version": schema.StringAttribute{
				Description:         "Devices with this macOS version or higher will not receive this app.",
				MarkdownDescription: "Devices with this macOS version or higher will **not** receive this app.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"macos_min_version": schema.StringAttribute{
				Description:         "Devices with this macOS version or higher will receive this app.",
				MarkdownDescription: "Devices with this macOS version or higher will receive this app.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tvos": schema.BoolAttribute{
				Description:         "Toggles the installation of the app on tvOS devices.",
				MarkdownDescription: "Toggles the installation of the app on tvOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"tvos_max_version": schema.StringAttribute{
				Description:         "Devices with this tvOS version or higher will not receive this app.",
				MarkdownDescription: "Devices with this tvOS version or higher will **not** receive this app.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tvos_min_version": schema.StringAttribute{
				Description:         "Devices with this tvOS version or higher will receive this app.",
				MarkdownDescription: "Devices with this tvOS version or higher will receive this app.",
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
				Description:         "Machines tagged with one of these tags will not receive the app.",
				MarkdownDescription: "Machines tagged with one of these tags will not receive the app.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.Int64Type, []attr.Value{})),
			},
			"tag_shards": schema.SetNestedAttribute{
				Description:         "A set of tag shard values different from the default shard, to determine if the tagged machines will receive the app.",
				MarkdownDescription: "A set of tag shard values different from the default shard, to determine if the tagged machines will receive the app.",
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
				Description:         "Version of the app.",
				MarkdownDescription: "Version of the app.",
				Required:            true,
			},
		},
	}
}

func (r *MDMStoreAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMStoreAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmStoreApp

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSA, _, err := r.client.MDMStoreApps.Create(ctx, mdmStoreAppRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM store app, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an MDM store app")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmStoreAppForState(ztlMSA))...)
}

func (r *MDMStoreAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmStoreApp

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSA, _, err := r.client.MDMStoreApps.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM store app %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read an MDM store app")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmStoreAppForState(ztlMSA))...)
}

func (r *MDMStoreAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmStoreApp

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSA, _, err := r.client.MDMStoreApps.Update(ctx, data.ID.ValueString(), mdmStoreAppRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM store app %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an MDM store app")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmStoreAppForState(ztlMSA))...)
}

func (r *MDMStoreAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmStoreApp

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMStoreApps.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM store app %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an MDM store app")
}

func (r *MDMStoreAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "MDM store app", req, resp)
}
