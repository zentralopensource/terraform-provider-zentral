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
var _ datasource.DataSource = &JMESPathCheckDataSource{}

func NewJMESPathCheckDataSource() datasource.DataSource {
	return &JMESPathCheckDataSource{}
}

// JMESPathCheckDataSource defines the data source implementation.
type JMESPathCheckDataSource struct {
	client *goztl.Client
}

func (d *JMESPathCheckDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jmespath_check"
}

func (d *JMESPathCheckDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a JMESPath compliance check to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `zentral_jmespath_check` allows details of a JMESPath compliance check to be retrieved by its `ID` or name.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the JMESPath compliance check.",
				MarkdownDescription: "`ID` of the JMESPath compliance check.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the JMESPath compliance check.",
				MarkdownDescription: "Name of the JMESPath compliance check.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the JMESPath compliance check.",
				MarkdownDescription: "Description of the JMESPath compliance check.",
				Computed:            true,
			},
			"source_name": schema.StringAttribute{
				Description:         "The name of the inventory source the JMESPath compliance check is restricted to.",
				MarkdownDescription: "The name of the inventory source the JMESPath compliance check is restricted to.",
				Computed:            true,
			},
			"platforms": schema.SetAttribute{
				Description:         "The platforms the JMESPath compliance check is restricted to.",
				MarkdownDescription: "The platforms the JMESPath compliance check is restricted to.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags the JMESPath compliance check is restricted to.",
				MarkdownDescription: "The IDs of the tags the JMESPath compliance check is restricted to.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"jmespath_expression": schema.StringAttribute{
				Description:         "The JMESPath compliance check expression.",
				MarkdownDescription: "The JMESPath compliance check expression.",
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				Description:         "The JMESPath compliance check version.",
				MarkdownDescription: "The JMESPath compliance check version.",
				Computed:            true,
			},
		},
	}
}

func (d *JMESPathCheckDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *JMESPathCheckDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data jmespathCheck
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_jmespath_check` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_jmespath_check` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *JMESPathCheckDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data jmespathCheck

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlJC *goztl.JMESPathCheck
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlJC, _, err = d.client.JMESPathChecks.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get JMESPath check '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlJC, _, err = d.client.JMESPathChecks.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get JMESPath check '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlJC != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, jmespathCheckForState(ztlJC))...)
	}
}
