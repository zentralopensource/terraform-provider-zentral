package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
var _ resource.Resource = &MDMACMEIssuerResource{}
var _ resource.ResourceWithImportState = &MDMACMEIssuerResource{}

func NewMDMACMEIssuerResource() resource.Resource {
	return &MDMACMEIssuerResource{}
}

// MDMACMEIssuerResource defines the resource implementation.
type MDMACMEIssuerResource struct {
	client *goztl.Client
}

func (r *MDMACMEIssuerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_acme_issuer"
}

func (r *MDMACMEIssuerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM ACME issuers.",
		MarkdownDescription: "The resource `zentral_mdm_acme_issuer` manages MDM ACME issuers.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the ACME issuer (UUID).",
				MarkdownDescription: "`ID` of the ACME issuer (UUID).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the ACME issuer.",
				MarkdownDescription: "Name of the ACME issuer.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the ACME issuer.",
				MarkdownDescription: "Description of the ACME issuer.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"directory_url": schema.StringAttribute{
				Description:         "Directory URL of the ACME issuer.",
				MarkdownDescription: "Directory URL of the ACME issuer.",
				Required:            true,
			},
			"key_type": schema.StringAttribute{
				Description:         "Private key type. One of 'ECSECPrimeRandom' or 'RSA'.",
				MarkdownDescription: "Private key type. One of `ECSECPrimeRandom` or `RSA`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{tfACMEIssuerECSECPrimeRandomKeyType, tfACMEIssuerRSAKeyType}...),
				},
			},
			"key_size": schema.Int64Attribute{
				Description:         "The key size in bits.",
				MarkdownDescription: "The key size in bits.",
				Required:            true,
			},
			"usage_flags": schema.Int64Attribute{
				Description:         "A bitmask that specifies the use of the key: '1' is signing, '4' is encryption, and '5' is both signing and encryption. Defaults to '0'.",
				MarkdownDescription: "A bitmask that specifies the use of the key: `1` is signing, `4` is encryption, and `5` is both signing and encryption. Defaults to `0`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(tfACMEIssuerDefaultKeyUsage),
				Validators: []validator.Int64{
					int64validator.OneOf([]int64{0, 1, 4, 5}...),
				},
			},
			"extended_key_usage": schema.SetAttribute{
				Description:         "The device requests this extended key usage for the certificate that the ACME server issues.",
				MarkdownDescription: "The device requests this extended key usage for the certificate that the ACME server issues.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"hardware_bound": schema.BoolAttribute{
				Description:         "Indicates if the private key must be bound to the device.",
				MarkdownDescription: "Indicates if the private key must be bound to the device.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"attest": schema.BoolAttribute{
				Description:         "Indicates if the device must provide an attestation that describe the device and the generated key to the ACME issuer.",
				MarkdownDescription: "Indicates if the device must provide an attestation that describe the device and the generated key to the ACME issuer.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"backend": schema.StringAttribute{
				Description:         "ACME issuer backend.",
				MarkdownDescription: "ACME issuer backend.",
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

func (r *MDMACMEIssuerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMACMEIssuerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmACMEIssuer

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMAI, _, err := r.client.MDMACMEIssuers.Create(ctx, mdmACMEIssuerRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM ACME issuer, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a MDM ACME issuer")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmACMEIssuerForState(ztlMAI))...)
}

func (r *MDMACMEIssuerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmACMEIssuer

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMAI, _, err := r.client.MDMACMEIssuers.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read ACME issuer %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a MDM ACME issuer")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmACMEIssuerForState(ztlMAI))...)
}

func (r *MDMACMEIssuerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmACMEIssuer

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMAI, _, err := r.client.MDMACMEIssuers.Update(ctx, data.ID.ValueString(), mdmACMEIssuerRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update ACME issuer %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a MDM ACME issuer")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmACMEIssuerForState(ztlMAI))...)
}

func (r *MDMACMEIssuerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmACMEIssuer

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMACMEIssuers.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete ACME issuer %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a MDM ACME issuer")
}

func (r *MDMACMEIssuerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "MDM ACME issuer", req, resp)
}
