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
var _ resource.Resource = &MDMEnterpriseAppResource{}
var _ resource.ResourceWithImportState = &MDMEnterpriseAppResource{}

func NewMDMEnterpriseAppResource() resource.Resource {
	return &MDMEnterpriseAppResource{}
}

// MDMEnterpriseAppResource defines the resource implementation.
type MDMEnterpriseAppResource struct {
	client *goztl.Client
}

func (r *MDMEnterpriseAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_enterprise_app"
}

func (r *MDMEnterpriseAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM enterprise apps.",
		MarkdownDescription: "The resource `zentral_mdm_enterprise_app` manages MDM enterprise apps.",

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
			"package_uri": schema.StringAttribute{
				Description:         "The URI of the app package.",
				MarkdownDescription: "The URI of the app package.",
				Required:            true,
			},
			"package_sha256": schema.StringAttribute{
				Description:         "The hexadecimal digest of the sha256 hash of the package.",
				MarkdownDescription: "The hexadecimal digest of the sha256 hash of the package.",
				Required:            true,
			},
			"ios_app": schema.BoolAttribute{
				Description:         "Indicates if the app is an iOS app that can run on an Apple silicon in macOS 11 and later. Defaults to false.",
				MarkdownDescription: "Indicates if the app is an iOS app that can run on an Apple silicon in macOS 11 and later. Defaults to `false`",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"configuration": schema.StringAttribute{
				Description:         "A dictionary serialized as plist that contains the initial configuration of the app.",
				MarkdownDescription: "A dictionary serialized as plist that contains the initial configuration of the app.",
				Optional:            true,
			},
			"install_as_managed": schema.BoolAttribute{
				Description:         "If true, install the app as a managed app. Defaults to false.",
				MarkdownDescription: "If `true`, install the app as a managed app. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"remove_on_unenroll": schema.BoolAttribute{
				Description:         "If true, the app will be removed when the device is unenrolled. Defaults to false.",
				MarkdownDescription: "If `true`, the app will be removed when the device is unenrolled. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
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

func (r *MDMEnterpriseAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMEnterpriseAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmEnterpriseApp

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMEA, _, err := r.client.MDMEnterpriseApps.Create(ctx, mdmEnterpriseAppRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM enterprise app, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an MDM enterprise app")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmEnterpriseAppForState(ztlMEA))...)
}

func (r *MDMEnterpriseAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmEnterpriseApp

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMEA, _, err := r.client.MDMEnterpriseApps.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM enterprise app %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read an MDM enterprise app")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmEnterpriseAppForState(ztlMEA))...)
}

func (r *MDMEnterpriseAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmEnterpriseApp

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMEA, _, err := r.client.MDMEnterpriseApps.Update(ctx, data.ID.ValueString(), mdmEnterpriseAppRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM enterprise app %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an MDM enterprise app")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmEnterpriseAppForState(ztlMEA))...)
}

func (r *MDMEnterpriseAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmEnterpriseApp

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMEnterpriseApps.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM enterprise app %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an MDM enterprise app")
}

func (r *MDMEnterpriseAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "MDM enterprise app", req, resp)
}
