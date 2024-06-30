package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MDMPushCertificateDataSource{}

func NewMDMPushCertificateDataSource() datasource.DataSource {
	return &MDMPushCertificateDataSource{}
}

// MDMPushCertificateDataSource defines the data source implementation.
type MDMPushCertificateDataSource struct {
	client *goztl.Client
}

func (d *MDMPushCertificateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_push_certificate"
}

func (d *MDMPushCertificateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM push certificate to be retrieved by its ID or its name.",
		MarkdownDescription: "The data source `zentral_mdm_push_certificate` allows details of a MDM push certificate to be retrieved by its `ID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the MDM push certificate.",
				MarkdownDescription: "`ID` of the MDM push certificate.",
				Optional:            true,
			},
			"provisioning_uid": schema.StringAttribute{
				Description:         "Provisioning UID of the push certificate.",
				MarkdownDescription: "Provisioning `UID` of the push certificate.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the push certificate.",
				MarkdownDescription: "Name of the push certificate.",
				Optional:            true,
			},
			"topic": schema.StringAttribute{
				Description:         "APNS topic the push certificate.",
				MarkdownDescription: "APNS topic of the push certificate.",
				Computed:            true,
			},
			"certificate": schema.StringAttribute{
				Description:         "Push certificate in PEM form.",
				MarkdownDescription: "Push certificate in `PEM` form.",
				Computed:            true,
			},
		},
	}
}

func (d *MDMPushCertificateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MDMPushCertificateDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmPushCertificate
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_push_certificate` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_push_certificate` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MDMPushCertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmPushCertificate

	// Read Terraform push certificate data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMPC *goztl.MDMPushCertificate
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlMPC, _, err = d.client.MDMPushCertificates.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM push certificate '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlMPC, _, err = d.client.MDMPushCertificates.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM push certificate '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMPC != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmPushCertificateForState(ztlMPC))...)
	}
}
