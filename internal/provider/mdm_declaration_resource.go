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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMDeclarationResource{}
var _ resource.ResourceWithImportState = &MDMDeclarationResource{}

func NewMDMDeclarationResource() resource.Resource {
	return &MDMDeclarationResource{}
}

// MDMDeclarationResource defines the resource implementation.
type MDMDeclarationResource struct {
	client *goztl.Client
}

func (r *MDMDeclarationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_declaration"
}

func (r *MDMDeclarationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM declarations.",
		MarkdownDescription: "The resource `zentral_mdm_declaration` manages MDM declarations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the declaration.",
				MarkdownDescription: "`ID` of the declaration.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"artifact_id": schema.StringAttribute{
				Description:         "ID of the declaration artifact.",
				MarkdownDescription: "`ID` of the declaration artifact.",
				Required:            true,
			},
			"source": schema.StringAttribute{
				Description:         "The actual DDM declaration (JSON).",
				MarkdownDescription: "The actual DDM declaration (JSON).",
				Required:            true,
				Validators:          []validator.String{mdmDeclarationSourceValidator{}},
			},
			"ios": schema.BoolAttribute{
				Description:         "Toggles the installation of the declaration on iOS devices.",
				MarkdownDescription: "Toggles the installation of the declaration on iOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ios_max_version": schema.StringAttribute{
				Description:         "Devices with this iOS version or higher will not receive this declaration.",
				MarkdownDescription: "Devices with this iOS version or higher will **not** receive this declaration.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ios_min_version": schema.StringAttribute{
				Description:         "Devices with this iOS version or higher will receive this declaration.",
				MarkdownDescription: "Devices with this iOS version or higher will receive this declaration.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ipados": schema.BoolAttribute{
				Description:         "Toggles the installation of the declaration on iPadOS devices.",
				MarkdownDescription: "Toggles the installation of the declaration on iPadOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ipados_max_version": schema.StringAttribute{
				Description:         "Devices with this iPadOS version or higher will not receive this declaration.",
				MarkdownDescription: "Devices with this iPadOS version or higher will **not** receive this declaration.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ipados_min_version": schema.StringAttribute{
				Description:         "Devices with this iPadOS version or higher will receive this declaration.",
				MarkdownDescription: "Devices with this iPadOS version or higher will receive this declaration.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"macos": schema.BoolAttribute{
				Description:         "Toggles the installation of the declaration on macOS devices.",
				MarkdownDescription: "Toggles the installation of the declaration on macOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"macos_max_version": schema.StringAttribute{
				Description:         "Devices with this macOS version or higher will not receive this declaration.",
				MarkdownDescription: "Devices with this macOS version or higher will **not** receive this declaration.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"macos_min_version": schema.StringAttribute{
				Description:         "Devices with this macOS version or higher will receive this declaration.",
				MarkdownDescription: "Devices with this macOS version or higher will receive this declaration.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tvos": schema.BoolAttribute{
				Description:         "Toggles the installation of the declaration on tvOS devices.",
				MarkdownDescription: "Toggles the installation of the declaration on tvOS devices.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"tvos_max_version": schema.StringAttribute{
				Description:         "Devices with this tvOS version or higher will not receive this declaration.",
				MarkdownDescription: "Devices with this tvOS version or higher will **not** receive this declaration.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tvos_min_version": schema.StringAttribute{
				Description:         "Devices with this tvOS version or higher will receive this declaration.",
				MarkdownDescription: "Devices with this tvOS version or higher will receive this declaration.",
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
				Description:         "Machines tagged with one of these tags will not receive the declaration.",
				MarkdownDescription: "Machines tagged with one of these tags will not receive the declaration.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.Int64Type, []attr.Value{})),
			},
			"tag_shards": schema.SetNestedAttribute{
				Description:         "A set of tag shard values different from the default shard, to determine if the tagged machines will receive the declaration.",
				MarkdownDescription: "A set of tag shard values different from the default shard, to determine if the tagged machines will receive the declaration.",
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
				Description:         "Version of the declaration.",
				MarkdownDescription: "Version of the declaration.",
				Required:            true,
			},
		},
	}
}

func (r *MDMDeclarationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMDeclarationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmDeclaration

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMDReq, err := mdmDeclarationRequestWithState(data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to prepare MDM declaration creation request, got error: %s", err),
		)
		return
	}

	ztlMD, _, err := r.client.MDMDeclarations.Create(ctx, ztlMDReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM declaration, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an MDM declaration")

	// Save data into Terraform state
	ztlMDFS, err := mdmDeclarationForState(ztlMD)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to prepare MDM declaration for state, got error: %s", err),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, ztlMDFS)...)
}

func (r *MDMDeclarationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmDeclaration

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMD, _, err := r.client.MDMDeclarations.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM declaration %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read an MDM declaration")

	// Save updated data into Terraform state
	ztlMDFS, err := mdmDeclarationForState(ztlMD)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to prepare MDM declaration for state, got error: %s", err),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, ztlMDFS)...)
}

func (r *MDMDeclarationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmDeclaration

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMDReq, err := mdmDeclarationRequestWithState(data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to prepare MDM declaration update request, got error: %s", err),
		)
		return
	}

	ztlMD, _, err := r.client.MDMDeclarations.Update(ctx, data.ID.ValueString(), ztlMDReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM declaration %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an MDM declaration")

	// Save updated data into Terraform state
	ztlMDFS, err := mdmDeclarationForState(ztlMD)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to prepare MDM declaration for state, got error: %s", err),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, ztlMDFS)...)
}

func (r *MDMDeclarationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmDeclaration

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMDeclarations.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM declaration %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an MDM declaration")
}

func (r *MDMDeclarationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "MDM declaration", req, resp)
}
