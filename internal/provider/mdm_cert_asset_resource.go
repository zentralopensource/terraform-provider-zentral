package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMCertAssetResource{}
var _ resource.ResourceWithImportState = &MDMCertAssetResource{}

func NewMDMCertAssetResource() resource.Resource {
	return &MDMCertAssetResource{}
}

// MDMCertAssetResource defines the resource implementation.
type MDMCertAssetResource struct {
	client *goztl.Client
}

func (r *MDMCertAssetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_cert_asset"
}

func (r *MDMCertAssetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM cert assets.",
		MarkdownDescription: "The resource `zentral_mdm_cert_asset` manages MDM cert assets.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the cert asset.",
				MarkdownDescription: "`ID` of the cert asset.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"acme_issuer_id": schema.StringAttribute{
				Description:         "ID of the optional MDM ACME issuer that could be used to issue the certificates.",
				MarkdownDescription: "`ID` of the optional MDM ACME issuer that could be used to issue the certificates.",
				Optional:            true,
			},
			"scep_issuer_id": schema.StringAttribute{
				Description:         "ID of the MDM SCEP issuer that could be used to issue the certificates.",
				MarkdownDescription: "`ID` of the MDM SCEP issuer that could be used to issue the certificates.",
				Optional:            true,
			},
			"subject": schema.SetNestedAttribute{
				Description:         "A set of relative distinguished names.",
				MarkdownDescription: "A set of relative distinguished names.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description:         "Type of the RDN. You can represent OIDs as dotted numbers or use shortcuts for country ('C'), locality ('L'), state ('ST'), organization ('O'), organizational unit ('OU'), and common name ('CN').",
							MarkdownDescription: "Type of the RDN. You can represent OIDs as dotted numbers or use shortcuts for country (`C`), locality (`L`), state (`ST`), organization (`O`), organizational unit (`OU`), and common name (`CN`).",
							Required:            true,
						},
						"value": schema.StringAttribute{
							Description:         "The value of the RDN.",
							MarkdownDescription: "The value of the RDN.",
							Required:            true,
						},
					},
				},
				Optional: true,
				Computed: true,
				Default:  setdefault.StaticValue(types.SetValueMust(types.ObjectType{AttrTypes: rdnAttrTypes}, []attr.Value{})),
			},
			"subject_alt_name": schema.SingleNestedAttribute{
				Description:         "The subject's alternative name for the certificate.",
				MarkdownDescription: "The subject's alternative name for the certificate.",
				Attributes: map[string]schema.Attribute{
					"rfc822_name": schema.StringAttribute{
						Description:         "The RFC 822 email address.",
						MarkdownDescription: "The RFC 822 email address.",
						Optional:            true,
					},
					"dns_name": schema.StringAttribute{
						Description:         "The DNS name.",
						MarkdownDescription: "The DNS name.",
						Optional:            true,
					},
					"uri": schema.StringAttribute{
						Description:         "The uniform resource identifier (URI).",
						MarkdownDescription: "The uniform resource identifier (URI).",
						Optional:            true,
					},
					"nt_principal_name": schema.StringAttribute{
						Description:         "The NT principal name. Use an OID set to '1.3.6.1.4.1.311.20.2.3'.",
						MarkdownDescription: "The NT principal name. Use an OID set to `1.3.6.1.4.1.311.20.2.3`.",
						Optional:            true,
					},
				},
				Optional: true,
				Computed: true,
				Default:  objectdefault.StaticValue(defaultSAN()),
			},
			"accessible": schema.StringAttribute{
				Description:         "The keychain accessibility that determines when the keychain item is available for use. Defaults to 'Default'.",
				MarkdownDescription: "The keychain accessibility that determines when the keychain item is available for use. Defaults to `Default`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(tfCertAssetAccessibilityDefault),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{tfCertAssetAccessibilityDefault, tfCertAssetAccessibilityAfterFirstUnlock}...),
				},
			},
			"artifact_id": schema.StringAttribute{
				Description:         "ID of the cert asset artifact.",
				MarkdownDescription: "`ID` of the cert asset artifact.",
				Required:            true,
			},
			"ios": schema.BoolAttribute{
				Description:         "Toggles the installation of the cert asset on iOS devices.",
				MarkdownDescription: "Toggles the installation of the cert asset on iOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ios_max_version": schema.StringAttribute{
				Description:         "Devices with this iOS version or higher will not receive this cert asset.",
				MarkdownDescription: "Devices with this iOS version or higher will **not** receive this cert asset.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ios_min_version": schema.StringAttribute{
				Description:         "Devices with this iOS version or higher will receive this cert asset.",
				MarkdownDescription: "Devices with this iOS version or higher will receive this cert asset.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ipados": schema.BoolAttribute{
				Description:         "Toggles the installation of the cert asset on iPadOS devices.",
				MarkdownDescription: "Toggles the installation of the cert asset on iPadOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ipados_max_version": schema.StringAttribute{
				Description:         "Devices with this iPadOS version or higher will not receive this cert asset.",
				MarkdownDescription: "Devices with this iPadOS version or higher will **not** receive this cert asset.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ipados_min_version": schema.StringAttribute{
				Description:         "Devices with this iPadOS version or higher will receive this cert asset.",
				MarkdownDescription: "Devices with this iPadOS version or higher will receive this cert asset.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"macos": schema.BoolAttribute{
				Description:         "Toggles the installation of the cert asset on macOS devices.",
				MarkdownDescription: "Toggles the installation of the cert asset on macOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"macos_max_version": schema.StringAttribute{
				Description:         "Devices with this macOS version or higher will not receive this cert asset.",
				MarkdownDescription: "Devices with this macOS version or higher will **not** receive this cert asset.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"macos_min_version": schema.StringAttribute{
				Description:         "Devices with this macOS version or higher will receive this cert asset.",
				MarkdownDescription: "Devices with this macOS version or higher will receive this cert asset.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tvos": schema.BoolAttribute{
				Description:         "Toggles the installation of the cert asset on tvOS devices.",
				MarkdownDescription: "Toggles the installation of the cert asset on tvOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"tvos_max_version": schema.StringAttribute{
				Description:         "Devices with this tvOS version or higher will not receive this cert asset.",
				MarkdownDescription: "Devices with this tvOS version or higher will **not** receive this cert asset.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tvos_min_version": schema.StringAttribute{
				Description:         "Devices with this tvOS version or higher will receive this cert asset.",
				MarkdownDescription: "Devices with this tvOS version or higher will receive this cert asset.",
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
				Description:         "Machines tagged with one of these tags will not receive the cert asset.",
				MarkdownDescription: "Machines tagged with one of these tags will not receive the cert asset.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.Int64Type, []attr.Value{})),
			},
			"tag_shards": schema.SetNestedAttribute{
				Description:         "A set of tag shard values different from the default shard, to determine if the tagged machines will receive the cert asset.",
				MarkdownDescription: "A set of tag shard values different from the default shard, to determine if the tagged machines will receive the cert asset.",
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
				Description:         "Version of the cert asset.",
				MarkdownDescription: "Version of the cert asset.",
				Required:            true,
			},
		},
	}
}

func (r *MDMCertAssetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMCertAssetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmCertAsset

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMCA, _, err := r.client.MDMCertAssets.Create(ctx, mdmCertAssetRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM cert asset, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an MDM cert asset")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmCertAssetForState(ztlMCA))...)
}

func (r *MDMCertAssetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmCertAsset

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMCA, _, err := r.client.MDMCertAssets.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM cert asset %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read an MDM cert asset")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmCertAssetForState(ztlMCA))...)
}

func (r *MDMCertAssetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmCertAsset

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMCA, _, err := r.client.MDMCertAssets.Update(ctx, data.ID.ValueString(), mdmCertAssetRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM cert asset %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an MDM cert asset")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmCertAssetForState(ztlMCA))...)
}

func (r *MDMCertAssetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmCertAsset

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMCertAssets.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM cert asset %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an MDM cert asset")
}

func (r *MDMCertAssetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "MDM cert asset", req, resp)
}
