package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMArtifactResource{}
var _ resource.ResourceWithImportState = &MDMArtifactResource{}

func NewMDMArtifactResource() resource.Resource {
	return &MDMArtifactResource{}
}

// MDMArtifactResource defines the resource implementation.
type MDMArtifactResource struct {
	client *goztl.Client
}

func (r *MDMArtifactResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_artifact"
}

func (r *MDMArtifactResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM artifacts.",
		MarkdownDescription: "The resource `zentral_mdm_artifact` manages MDM artifacts.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the artifact.",
				MarkdownDescription: "`ID` of the artifact.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the artifact.",
				MarkdownDescription: "Name of the artifact.",
				Required:            true,
			},
			"type": schema.StringAttribute{
				Description:         "Type of the artifact.",
				MarkdownDescription: "Type of the artifact.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"Activation",
						"Asset",
						"Configuration",
						"Certificate Asset",
						"Data Asset",
						"Enterprise App",
						"Configuration (manual)",
						"Profile",
						"Store App",
					}...),
				},
			},
			"channel": schema.StringAttribute{
				Description:         "Channel of the artifact.",
				MarkdownDescription: "Channel of the artifact.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"Device", "User"}...),
				},
			},
			"platforms": schema.SetAttribute{
				Description:         "Platforms of the artifact.",
				MarkdownDescription: "Platforms of the artifact.",
				ElementType:         types.StringType,
				Required:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"macOS", "iOS", "iPadOS", "tvOS"}...)),
					setvalidator.SizeAtLeast(1),
				},
			},
			"install_during_setup_assistant": schema.BoolAttribute{
				Description:         "If true, this artifact will be installed during the setup assistant. Defaults to false.",
				MarkdownDescription: "If `true`, this artifact will be installed during the setup assistant. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"auto_update": schema.BoolAttribute{
				Description:         "If true, new version of this artifact will be automatically installed. Defaults to true.",
				MarkdownDescription: "If `true`, new version of this artifact will be automatically installed. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"reinstall_interval": schema.Int64Attribute{
				Description:         "In days, the time interval after which the artifact will be reinstalled. If 0, the artifact will not be reinstalled. Defaults to 0.",
				MarkdownDescription: "In days, the time interval after which the artifact will be reinstalled. If `0`, the artifact will not be reinstalled. Defaults to `0`.",
				Optional:            true,
				Default:             int64default.StaticInt64(0),
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 366),
				},
			},
			"reinstall_on_os_update": schema.StringAttribute{
				Description:         "Possible values: No, Major, Minor, Patch. Defaults to No.",
				MarkdownDescription: "Possible values: `No`, `Major`, `Minor`, `Patch`. Defaults to `No`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("No"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"No", "Major", "Minor", "Patch"}...),
				},
			},
			"requires": schema.SetAttribute{
				Description:         "IDs of the artifacts required by this artifact.",
				MarkdownDescription: "`ID`s of the artifacts required by this artifact.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *MDMArtifactResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMArtifactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmArtifact

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMA, _, err := r.client.MDMArtifacts.Create(ctx, mdmArtifactRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM artifact, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a MDM artifact")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmArtifactForState(ztlMA))...)
}

func (r *MDMArtifactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmArtifact

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMA, _, err := r.client.MDMArtifacts.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM artifact %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a MDM artifact")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmArtifactForState(ztlMA))...)
}

func (r *MDMArtifactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmArtifact

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMA, _, err := r.client.MDMArtifacts.Update(ctx, data.ID.ValueString(), mdmArtifactRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM artifact %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a MDM artifact")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmArtifactForState(ztlMA))...)
}

func (r *MDMArtifactResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmArtifact

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMArtifacts.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM artifact %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a MDM artifact")
}

func (r *MDMArtifactResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "MDM artifact", req, resp)
}
