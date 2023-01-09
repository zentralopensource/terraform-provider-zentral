package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
	"github.com/zentralopensource/terraform-provider-zentral/internal/planmodifiers"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MetaBusinessUnitResource{}
var _ resource.ResourceWithImportState = &MetaBusinessUnitResource{}

func NewMetaBusinessUnitResource() resource.Resource {
	return &MetaBusinessUnitResource{}
}

// MetaBusinessUnitResource defines the resource implementation.
type MetaBusinessUnitResource struct {
	client *goztl.Client
}

func (r *MetaBusinessUnitResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_meta_business_unit"
}

func (r *MetaBusinessUnitResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages meta business units.",
		MarkdownDescription: "The resource `zentral_meta_business_unit` manages meta business units.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the meta business unit.",
				MarkdownDescription: "`ID` of the meta business unit.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the meta business unit.",
				MarkdownDescription: "Name of the meta business unit.",
				Required:            true,
			},
			"api_enrollment_enabled": schema.BoolAttribute{
				Description: "Enables API enrollments.",
				MarkdownDescription: "Enables API enrollments. Once enabled, it **CANNOT** be disabled. " +
					"Defaults to `true`.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefault(true),
				},
			},
		},
	}
}

func (r *MetaBusinessUnitResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MetaBusinessUnitResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data metaBusinessUnit

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	mbuCreateRequest := &goztl.MetaBusinessUnitCreateRequest{
		Name: data.Name.ValueString(),
	}
	if data.APIEnrollmentEnabled.ValueBool() {
		mbuCreateRequest.APIEnrollmentEnabled = true
	}
	mbu, _, err := r.client.MetaBusinessUnits.Create(ctx, mbuCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create meta business unit, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a meta business unit")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, metaBusinessUnitForState(mbu))...)
}

func (r *MetaBusinessUnitResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data metaBusinessUnit

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	mbu, _, err := r.client.MetaBusinessUnits.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read meta business unit %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a meta business unit")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, metaBusinessUnitForState(mbu))...)
}

func (r *MetaBusinessUnitResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data metaBusinessUnit

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	mbuUpdateRequest := &goztl.MetaBusinessUnitUpdateRequest{
		Name: data.Name.ValueString(),
	}
	if data.APIEnrollmentEnabled.ValueBool() {
		mbuUpdateRequest.APIEnrollmentEnabled = true
	}
	mbu, _, err := r.client.MetaBusinessUnits.Update(ctx, int(data.ID.ValueInt64()), mbuUpdateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update meta business unit %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a meta business unit")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, metaBusinessUnitForState(mbu))...)
}

func (r *MetaBusinessUnitResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data metaBusinessUnit

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MetaBusinessUnits.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete meta business unit %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a meta business unit")
}

func (r *MetaBusinessUnitResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "meta business unit", req, resp)
}
