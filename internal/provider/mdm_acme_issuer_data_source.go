package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MDMACMEIssuerDataSource{}

func NewMDMACMEIssuerDataSource() datasource.DataSource {
	return &MDMACMEIssuerDataSource{}
}

// MDMACMEIssuerDataSource defines the data source implementation.
type MDMACMEIssuerDataSource struct {
	client *goztl.Client
}

func (d *MDMACMEIssuerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_acme_issuer"
}

func (d *MDMACMEIssuerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM ACME issuer to be retrieved by its ID or its name.",
		MarkdownDescription: "The data source `zentral_mdm_acme_issuer` allows details of a MDM ACME issuer to be retrieved by its `ID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the MDM ACME issuer (UUID).",
				MarkdownDescription: "`ID` of the MDM ACME issuer (UUID).",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the ACME issuer.",
				MarkdownDescription: "Name of the ACME issuer.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the ACME issuer.",
				MarkdownDescription: "Description of the ACME issuer.",
				Computed:            true,
			},
			"directory_url": schema.StringAttribute{
				Description:         "Directory URL of the ACME server.",
				MarkdownDescription: "Directory URL of the ACME server.",
				Computed:            true,
			},
			"key_type": schema.StringAttribute{
				Description:         "Private key type for the ACME requests.",
				MarkdownDescription: "Private key type for the ACME requests.",
				Computed:            true,
			},
			"key_size": schema.Int64Attribute{
				Description:         "Private key size for the ACME requests.",
				MarkdownDescription: "Private key size for the ACME requests.",
				Computed:            true,
			},
			"usage_flags": schema.Int64Attribute{
				Description:         "Private key usage flags for the ACME requests.",
				MarkdownDescription: "Private key usage flags for the ACME requests.",
				Computed:            true,
			},
			"extended_key_usage": schema.SetAttribute{
				Description:         "The device requests this extended key usage for the certificate that the ACME server issues.",
				MarkdownDescription: "The device requests this extended key usage for the certificate that the ACME server issues.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"hardware_bound": schema.BoolAttribute{
				Description:         "Indicates if the private key must be bound to the device.",
				MarkdownDescription: "Indicates if the private key must be bound to the device.",
				Computed:            true,
			},
			"attest": schema.BoolAttribute{
				Description:         "Indicates if the device must provide an attestation that describe the device and the generated key to the ACME issuer.",
				MarkdownDescription: "Indicates if the device must provide an attestation that describe the device and the generated key to the ACME issuer.",
				Computed:            true,
			},
			"backend": schema.StringAttribute{
				Description:         "ACME issuer backend.",
				MarkdownDescription: "ACME issuer backend.",
				Computed:            true,
			},
			"ident":            makeIDentBackendDataSourceAttribute(),
			"microsoft_ca":     makeMicrosoftCABackendDataSourceAttribute("Microsoft CA"),
			"okta_ca":          makeMicrosoftCABackendDataSourceAttribute("Okta CA"),
			"static_challenge": makeStaticChallengeBackendDataSourceAttribute(),
		},
	}
}

func (d *MDMACMEIssuerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*goztl.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *goztl.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *MDMACMEIssuerDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmACMEIssuer
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_acme_issuer` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_acme_issuer` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MDMACMEIssuerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmACMEIssuer

	// Read Terraform ACME issuer data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMAI *goztl.MDMACMEIssuer
	var err error
	if !data.ID.IsNull() {
		ztlMAI, _, err = d.client.MDMACMEIssuers.GetByID(ctx, data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM ACME issuer '%s' by ID, got error: %s", data.ID.ValueString(), err),
			)
		}
	} else {
		ztlMAI, _, err = d.client.MDMACMEIssuers.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM ACME issuer '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMAI != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmACMEIssuerForState(ztlMAI))...)
	}
}
