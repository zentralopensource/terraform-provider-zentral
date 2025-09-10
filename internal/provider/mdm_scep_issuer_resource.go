package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMSCEPIssuerResource{}
var _ resource.ResourceWithImportState = &MDMSCEPIssuerResource{}

func NewMDMSCEPIssuerResource() resource.Resource {
	return &MDMSCEPIssuerResource{}
}

// MDMSCEPIssuerResource defines the resource implementation.
type MDMSCEPIssuerResource struct {
	client *goztl.Client
}

func (r *MDMSCEPIssuerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_scep_issuer"
}

func (r *MDMSCEPIssuerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM SCEP issuers.",
		MarkdownDescription: "The resource `zentral_mdm_scep_issuer` manages MDM SCEP issuers.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the SCEP issuer (UUID).",
				MarkdownDescription: "`ID` of the SCEP issuer (UUID).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the SCEP issuer.",
				MarkdownDescription: "Name of the SCEP issuer.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the SCEP issuer.",
				MarkdownDescription: "Description of the SCEP issuer.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"url": schema.StringAttribute{
				Description:         "URL of the SCEP issuer.",
				MarkdownDescription: "URL of the SCEP issuer.",
				Required:            true,
			},
			"key_size": schema.Int64Attribute{
				Description:         "The key size in bits, either '1024', '2048', or '4096'. Defaults to '2048'.",
				MarkdownDescription: "The key size in bits, either `1024`, `2048`, or `4096`. Defaults to `2048`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(tfSCEPIssuerDefaultKeySize),
				Validators: []validator.Int64{
					int64validator.OneOf([]int64{1024, tfSCEPIssuerDefaultKeySize, 4096}...),
				},
			},
			"key_usage": schema.Int64Attribute{
				Description:         "A bitmask that specifies the use of the key: '1' is signing, '4' is encryption, and '5' is both signing and encryption. Defaults to '0'.",
				MarkdownDescription: "A bitmask that specifies the use of the key: `1` is signing, `4` is encryption, and `5` is both signing and encryption. Defaults to `0`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(tfSCEPIssuerDefaultKeyUsage),
				Validators: []validator.Int64{
					int64validator.OneOf([]int64{0, 1, 4, 5}...),
				},
			},
			"backend": schema.StringAttribute{
				Description:         "SCEP issuer backend.",
				MarkdownDescription: "SCEP issuer backend.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						tfCertIssuerIDentBackend,
						tfCertIssuerMicrosoftCABackend,
						tfCertIssuerOktaCABackend,
						tfCertIssuerStaticChallengeBackend,
					}...),
				},
			},
			"ident":            makeIDentBackendResourceAttribute(),
			"microsoft_ca":     makeMicrosoftCABackendResourceAttribute("Microsoft CA"),
			"okta_ca":          makeMicrosoftCABackendResourceAttribute("Okta CA"),
			"static_challenge": makeStaticChallengeBackendResourceAttribute(),
		},
	}
}

func (r *MDMSCEPIssuerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMSCEPIssuerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmSCEPIssuer

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSI, _, err := r.client.MDMSCEPIssuers.Create(ctx, mdmSCEPIssuerRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM SCEP issuer, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a MDM SCEP issuer")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmSCEPIssuerForState(ztlMSI))...)
}

func (r *MDMSCEPIssuerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmSCEPIssuer

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSI, _, err := r.client.MDMSCEPIssuers.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read SCEP issuer %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a MDM SCEP issuer")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmSCEPIssuerForState(ztlMSI))...)
}

func (r *MDMSCEPIssuerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmSCEPIssuer

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMSI, _, err := r.client.MDMSCEPIssuers.Update(ctx, data.ID.ValueString(), mdmSCEPIssuerRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update SCEP issuer %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a MDM SCEP issuer")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmSCEPIssuerForState(ztlMSI))...)
}

func (r *MDMSCEPIssuerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmSCEPIssuer

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMSCEPIssuers.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete SCEP issuer %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a MDM SCEP issuer")
}

func (r *MDMSCEPIssuerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "MDM SCEP issuer", req, resp)
}
