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
var _ resource.Resource = &MonolithManifestEnrollmentPackageResource{}
var _ resource.ResourceWithImportState = &MonolithManifestEnrollmentPackageResource{}

func NewMonolithManifestEnrollmentPackageResource() resource.Resource {
	return &MonolithManifestEnrollmentPackageResource{}
}

// MonolithManifestEnrollmentPackageResource defines the resource implementation.
type MonolithManifestEnrollmentPackageResource struct {
	client *goztl.Client
}

func (r *MonolithManifestEnrollmentPackageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monolith_manifest_enrollment_package"
}

func (r *MonolithManifestEnrollmentPackageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Monolith manifest enrollment packages.",
		MarkdownDescription: "The resource `zentral_monolith_manifest_enrollment_package` manages Monolith manifest enrollment packages.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the manifest enrollment package.",
				MarkdownDescription: "`ID` of the manifest enrollment package.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"manifest_id": schema.Int64Attribute{
				Description:         "ID of the manifest.",
				MarkdownDescription: "ID of the manifest.",
				Required:            true,
			},
			"builder": schema.StringAttribute{
				Description:         "Enrollment package builder module.",
				MarkdownDescription: "Enrollment package builder module.",
				Required:            true,
			},
			"enrollment_id": schema.Int64Attribute{
				Description:         "ID of the enrollment.",
				MarkdownDescription: "ID of the enrollment.",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				Description:         "Version of the enrollment package.",
				MarkdownDescription: "Version of the enrollment package.",
				Computed:            true,
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags used to scope the enrollment package.",
				MarkdownDescription: "The `ID`s of the tags used to scope the enrollment package.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *MonolithManifestEnrollmentPackageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MonolithManifestEnrollmentPackageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data monolithManifestEnrollmentPackage

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMMEP, _, err := r.client.MonolithManifestEnrollmentPackages.Create(ctx, monolithManifestEnrollmentPackageRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Monolith manifest enrollment package, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an Monolith manifest enrollment package")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithManifestEnrollmentPackageForState(ztlMMEP))...)
}

func (r *MonolithManifestEnrollmentPackageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data monolithManifestEnrollmentPackage

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMMEP, _, err := r.client.MonolithManifestEnrollmentPackages.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Monolith manifest enrollment package %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read an Monolith manifest enrollment package")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithManifestEnrollmentPackageForState(ztlMMEP))...)
}

func (r *MonolithManifestEnrollmentPackageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data monolithManifestEnrollmentPackage

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMMEP, _, err := r.client.MonolithManifestEnrollmentPackages.Update(ctx, int(data.ID.ValueInt64()), monolithManifestEnrollmentPackageRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Monolith manifest enrollment package %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an Monolith manifest enrollment package")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithManifestEnrollmentPackageForState(ztlMMEP))...)
}

func (r *MonolithManifestEnrollmentPackageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data monolithManifestEnrollmentPackage

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MonolithManifestEnrollmentPackages.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Monolith manifest enrollment package %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an Monolith manifest enrollment package")
}

func (r *MonolithManifestEnrollmentPackageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Monolith manifest enrollment package", req, resp)
}
