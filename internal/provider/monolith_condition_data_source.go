package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MonolithConditionDataSource{}

func NewMonolithConditionDataSource() datasource.DataSource {
	return &MonolithConditionDataSource{}
}

// MonolithConditionDataSource defines the data source implementation.
type MonolithConditionDataSource struct {
	client *goztl.Client
}

func (d *MonolithConditionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monolith_condition"
}

func (d *MonolithConditionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Monolith condition to be retrieved by its ID or its name.",
		MarkdownDescription: "The data source `zentral_monolith_condition` allows details of a Monolith condition to be retrieved by its `ID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Monolith condition.",
				MarkdownDescription: "`ID` of the Monolith condition.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the condition.",
				MarkdownDescription: "Name of the condition.",
				Optional:            true,
			},
			"predicate": schema.StringAttribute{
				Description:         "Predicate of the condition.",
				MarkdownDescription: "Predicate of the condition.",
				Computed:            true,
			},
		},
	}
}

func (d *MonolithConditionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonolithConditionDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data monolithCondition
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_monolith_condition` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_monolith_condition` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MonolithConditionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data monolithCondition

	// Read Terraform condition data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMM *goztl.MonolithCondition
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlMM, _, err = d.client.MonolithConditions.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Monolith condition '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlMM, _, err = d.client.MonolithConditions.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Monolith condition '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMM != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, monolithConditionForState(ztlMM))...)
	}
}
