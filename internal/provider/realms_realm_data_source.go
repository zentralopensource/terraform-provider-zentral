package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &RealmsRealmDataSource{}

func NewRealmsRealmDataSource() datasource.DataSource {
	return &RealmsRealmDataSource{}
}

// RealmsRealmDataSource defines the data source implementation.
type RealmsRealmDataSource struct {
	client *goztl.Client
}

func (d *RealmsRealmDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_realm"
}

func (d *RealmsRealmDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a realm to be retrieved by its UUID or its name.",
		MarkdownDescription: "The data source `zentral_realm` allows details of a realm to be retrieved by its `UUID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "alias for the UUID of the realm.",
				MarkdownDescription: "alias for the UUID of the realm.",
				Computed:            true,
			},
			"uuid": schema.StringAttribute{
				Description:         "UUID of the realm.",
				MarkdownDescription: "`UUID` of the realm.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the realm.",
				MarkdownDescription: "Name of the realm.",
				Optional:            true,
			},
		},
	}
}

func (d *RealmsRealmDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RealmsRealmDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data realmsRealm
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.UUID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_realm` data source",
			"`uuid` or `name` missing",
		)
	} else if !data.UUID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_realm` data source",
			"`uuid` and `name` cannot be both set",
		)
	}
}

func (d *RealmsRealmDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data realmsRealm

	// Read Terraform SCEP config data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlR *goztl.RealmsRealm
	var err error
	if !data.UUID.IsNull() {
		ztlR, _, err = d.client.RealmsRealms.GetByUUID(ctx, data.UUID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get realm '%s' by UUID, got error: %s", data.UUID.ValueString(), err),
			)
		}
	} else {
		ztlR, _, err = d.client.RealmsRealms.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get realm '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlR != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, realmsRealmForState(ztlR))...)
	}
}
