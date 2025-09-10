package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MDMSCEPIssuerDataSource{}

func NewMDMSCEPIssuerDataSource() datasource.DataSource {
	return &MDMSCEPIssuerDataSource{}
}

// MDMSCEPIssuerDataSource defines the data source implementation.
type MDMSCEPIssuerDataSource struct {
	client *goztl.Client
}

func (d *MDMSCEPIssuerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_scep_issuer"
}

func (d *MDMSCEPIssuerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM SCEP issuer to be retrieved by its ID or its name.",
		MarkdownDescription: "The data source `zentral_mdm_scep_issuer` allows details of a MDM SCEP issuer to be retrieved by its `ID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the MDM SCEP issuer (UUID).",
				MarkdownDescription: "`ID` of the MDM SCEP issuer (UUID).",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the SCEP issuer.",
				MarkdownDescription: "Name of the SCEP issuer.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the SCEP issuer.",
				MarkdownDescription: "Description of the SCEP issuer.",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				Description:         "URL of the SCEP server.",
				MarkdownDescription: "URL of the SCEP server.",
				Computed:            true,
			},
			"key_size": schema.Int64Attribute{
				Description:         "The size of the private key in bits.",
				MarkdownDescription: "The size of the private key in bits.",
				Computed:            true,
			},
			"key_usage": schema.Int64Attribute{
				Description:         "Key usage for the SCEP requests.",
				MarkdownDescription: "Key usage for the SCEP requests.",
				Computed:            true,
			},
			"backend": schema.StringAttribute{
				Description:         "SCEP issuer backend.",
				MarkdownDescription: "SCEP issuer backend.",
				Computed:            true,
			},
			"ident":            makeIDentBackendDataSourceAttribute(),
			"microsoft_ca":     makeMicrosoftCABackendDataSourceAttribute("Microsoft CA"),
			"okta_ca":          makeMicrosoftCABackendDataSourceAttribute("Okta CA"),
			"static_challenge": makeStaticChallengeBackendDataSourceAttribute(),
		},
	}
}

func (d *MDMSCEPIssuerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MDMSCEPIssuerDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmSCEPIssuer
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_scep_issuer` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_scep_issuer` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MDMSCEPIssuerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmSCEPIssuer

	// Read Terraform SCEP issuer data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMSI *goztl.MDMSCEPIssuer
	var err error
	if !data.ID.IsNull() {
		ztlMSI, _, err = d.client.MDMSCEPIssuers.GetByID(ctx, data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM SCEP issuer '%s' by ID, got error: %s", data.ID.ValueString(), err),
			)
		}
	} else {
		ztlMSI, _, err = d.client.MDMSCEPIssuers.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM SCEP issuer '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMSI != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmSCEPIssuerForState(ztlMSI))...)
	}
}
