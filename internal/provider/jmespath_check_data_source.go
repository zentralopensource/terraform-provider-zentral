package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = jmespathCheckDataSource{}
var _ provider.DataSourceType = jmespathCheckDataSourceType{}

type jmespathCheckDataSourceType struct{}

func (t jmespathCheckDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Allows details of a JMESPath compliance check to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `zentral_jmespath_check` allows details of a JMESPath compliance check to be retrieved by its `ID` or name.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description:         "ID of the JMESPath compliance check.",
				MarkdownDescription: "`ID` of the JMESPath compliance check.",
				Type:                types.Int64Type,
				Optional:            true,
			},
			"name": {
				Description:         "Name of the JMESPath compliance check.",
				MarkdownDescription: "Name of the JMESPath compliance check.",
				Type:                types.StringType,
				Optional:            true,
			},
			"description": {
				Description:         "Description of the JMESPath compliance check.",
				MarkdownDescription: "Description of the JMESPath compliance check.",
				Type:                types.StringType,
				Computed:            true,
			},
			"source_name": {
				Description:         "The name of the inventory source the JMESPath compliance check is restricted to.",
				MarkdownDescription: "The name of the inventory source the JMESPath compliance check is restricted to.",
				Type:                types.StringType,
				Computed:            true,
			},
			"platforms": {
				Description:         "The platforms the JMESPath compliance check is restricted to.",
				MarkdownDescription: "The platforms the JMESPath compliance check is restricted to.",
				Type:                types.SetType{ElemType: types.StringType},
				Computed:            true,
			},
			"tag_ids": {
				Description:         "The IDs of the tags the JMESPath compliance check is restricted to.",
				MarkdownDescription: "The IDs of the tags the JMESPath compliance check is restricted to.",
				Type:                types.SetType{ElemType: types.Int64Type},
				Computed:            true,
			},
			"jmespath_expression": {
				Description:         "The JMESPath compliance check expression.",
				MarkdownDescription: "The JMESPath compliance check expression.",
				Type:                types.StringType,
				Computed:            true,
			},
			"version": {
				Description:         "The JMESPath compliance check version.",
				MarkdownDescription: "The JMESPath compliance check version.",
				Type:                types.Int64Type,
				Computed:            true,
			},
		},
	}, nil
}

func (t jmespathCheckDataSourceType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return jmespathCheckDataSource{
		provider: provider,
	}, diags
}

type jmespathCheckDataSource struct {
	provider zentralProvider
}

func (d jmespathCheckDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data jmespathCheck
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.Null && data.Name.Null {
		resp.Diagnostics.AddError(
			"Invalid `zentral_jmespath_check` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.Null && !data.Name.Null {
		resp.Diagnostics.AddError(
			"Invalid `zentral_jmespath_check` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d jmespathCheckDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data jmespathCheck

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ztlJC *goztl.JMESPathCheck
	var err error
	if data.ID.Value > 0 {
		ztlJC, _, err = d.provider.client.JMESPathChecks.GetByID(ctx, int(data.ID.Value))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get JMESPath check '%d' by ID, got error: %s", data.ID.Value, err),
			)
		}
	} else {
		ztlJC, _, err = d.provider.client.JMESPathChecks.GetByName(ctx, data.Name.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get JMESPath check '%s' by name, got error: %s", data.Name.Value, err),
			)
		}
	}

	if ztlJC != nil {
		diags = resp.State.Set(ctx, jmespathCheckForState(ztlJC))
		resp.Diagnostics.Append(diags...)
	}
}
