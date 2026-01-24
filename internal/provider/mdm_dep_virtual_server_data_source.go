package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MDMDEPVirtualServerDataSource{}

func NewMDMDEPVirtualServerDataSource() datasource.DataSource {
	return &MDMDEPVirtualServerDataSource{}
}

// MDMDEPVirtualServerDataSource defines the data source implementation.
type MDMDEPVirtualServerDataSource struct {
	client *goztl.Client
}

func (d *MDMDEPVirtualServerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_dep_virtual_server"
}

func (d *MDMDEPVirtualServerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a MDM DEP virtual server to be retrieved by its ID and name.",
		MarkdownDescription: "The data source `zentral_mdm_dep_virtual_server` allows details of a MDM DEP virtual server to be retrieved by its `ID` or `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the MDM DEP virtual server.",
				MarkdownDescription: "ID of the MDM DEP virtual server.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the MDM DEP virtual server.",
				MarkdownDescription: "Name of the MDM DEP virtual server.",
				Optional:            true,
			},
		},
	}
}

func (d *MDMDEPVirtualServerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MDMDEPVirtualServerDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data mdmDEPVirtualServer
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_dep_virtual_server` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_mdm_dep_virtual_server` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MDMDEPVirtualServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mdmDEPVirtualServer

	// Read Terraform enrollment data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlVirtualServer *goztl.MDMDEPVirtualServer
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlVirtualServer, _, err = d.client.MDMDEPVirtualServers.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM DEP virtual server '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		var result []goztl.MDMDEPVirtualServer
		result, _, err = d.client.MDMDEPVirtualServers.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get MDM DEP virtual server '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
		if len(result) == 1 {
			ztlVirtualServer = &result[0]
		} else {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("No unique result for MDM DEP virtual server'%s' by name, found %d results.", data.Name.ValueString(), len(result)),
			)
		}
	}

	if ztlVirtualServer != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, mdmDEPVirtualServerForState(ztlVirtualServer))...)
	}
}
