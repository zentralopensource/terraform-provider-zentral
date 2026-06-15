package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMPackageResource{}
var _ resource.ResourceWithImportState = &MDMPackageResource{}

func NewMDMPackageResource() resource.Resource {
	return &MDMPackageResource{}
}

// MDMPackageResource defines the resource implementation.
type MDMPackageResource struct {
	client *goztl.Client
}

func (r *MDMPackageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_package"
}

func (r *MDMPackageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM packages.",
		MarkdownDescription: "The resource `zentral_mdm_package` manages MDM packages — `.pkg` or `.ipa` payloads served to devices via the MDM `ManifestURL` contract (InstallEnterpriseApplication, InstallApplication, DDM `com.apple.configuration.package`).",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the package.",
				MarkdownDescription: "`ID` of the package.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the package.",
				MarkdownDescription: "Name of the package.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the package.",
				MarkdownDescription: "Description of the package.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"source_uri": schema.StringAttribute{
				Description:         "URI the server downloads the package from at create time (s3://, https://, …). Immutable: changes force replacement.",
				MarkdownDescription: "URI the server downloads the package from at create time (`s3://`, `https://`, …). Immutable: changes force replacement.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sha256": schema.StringAttribute{
				Description:         "Hexadecimal digest of the SHA-256 hash of the package. Verified server-side against the downloaded file. Immutable: changes force replacement.",
				MarkdownDescription: "Hexadecimal digest of the SHA-256 hash of the package. Verified server-side against the downloaded file. Immutable: changes force replacement.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description:         "Type of the package, derived from the file extension at upload time (PKG or IPA).",
				MarkdownDescription: "Type of the package, derived from the file extension at upload time (`PKG` or `IPA`).",
				Computed:            true,
			},
			"size": schema.Int64Attribute{
				Description:         "Size in bytes of the uploaded package.",
				MarkdownDescription: "Size in bytes of the uploaded package.",
				Computed:            true,
			},
			"filename": schema.StringAttribute{
				Description:         "Original filename of the uploaded package.",
				MarkdownDescription: "Original filename of the uploaded package.",
				Computed:            true,
			},
			"product_id": schema.StringAttribute{
				Description:         "Product identifier extracted from the package metadata.",
				MarkdownDescription: "Product identifier extracted from the package metadata.",
				Computed:            true,
			},
			"product_version": schema.StringAttribute{
				Description:         "Product version extracted from the package metadata.",
				MarkdownDescription: "Product version extracted from the package metadata.",
				Computed:            true,
			},
			"bundles": schema.StringAttribute{
				Description:         "JSON-encoded list of bundles extracted from the package.",
				MarkdownDescription: "JSON-encoded list of bundles extracted from the package.",
				Computed:            true,
			},
			"manifest": schema.StringAttribute{
				Description:         "JSON-encoded Apple ManifestURL manifest generated for this package.",
				MarkdownDescription: "JSON-encoded Apple `ManifestURL` manifest generated for this package.",
				Computed:            true,
			},
		},
	}
}

func (r *MDMPackageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMPackageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmPackage

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMP, _, err := r.client.MDMPackages.Create(ctx, mdmPackageCreateRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM package, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a MDM package")

	mp, err := mdmPackageForState(ztlMP)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to serialize MDM package, got error: %s", err),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mp)...)
}

func (r *MDMPackageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmPackage

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMP, _, err := r.client.MDMPackages.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM package %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a MDM package")

	mp, err := mdmPackageForState(ztlMP)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to serialize MDM package, got error: %s", err),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mp)...)
}

func (r *MDMPackageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmPackage

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMP, _, err := r.client.MDMPackages.Update(ctx, data.ID.ValueString(), mdmPackageUpdateRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM package %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a MDM package")

	mp, err := mdmPackageForState(ztlMP)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to serialize MDM package, got error: %s", err),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mp)...)
}

func (r *MDMPackageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmPackage

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMPackages.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM package %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a MDM package")
}

func (r *MDMPackageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "MDM package", req, resp)
}
