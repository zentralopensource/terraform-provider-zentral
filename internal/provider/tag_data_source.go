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
var _ datasource.DataSource = tagDataSource{}
var _ provider.DataSourceType = tagDataSourceType{}

type tagDataSourceType struct{}

func (t tagDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Allows details of a tag to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `zentral_tag` allows details of a tag to be retrieved by its `ID` or name.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description:         "ID of the tag.",
				MarkdownDescription: "`ID` of the tag.",
				Type:                types.Int64Type,
				Optional:            true,
			},
			"taxonomy_id": {
				Description:         "ID of the tag taxonomy.",
				MarkdownDescription: "`ID` of the tag taxonomy.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"name": {
				Description:         "Name of the tag.",
				MarkdownDescription: "Name of the tag.",
				Type:                types.StringType,
				Optional:            true,
			},
			"color": {
				Description:         "Color of the tag.",
				MarkdownDescription: "Color of the tag.",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (t tagDataSourceType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return tagDataSource{
		provider: provider,
	}, diags
}

type tagDataSource struct {
	provider zentralProvider
}

func (d tagDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data tag
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.Null && data.Name.Null {
		resp.Diagnostics.AddError(
			"Invalid `zentral_tag` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.Null && !data.Name.Null {
		resp.Diagnostics.AddError(
			"Invalid `zentral_tag` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d tagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data tag

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ztlTag *goztl.Tag
	var err error
	if data.ID.Value > 0 {
		ztlTag, _, err = d.provider.client.Tags.GetByID(ctx, int(data.ID.Value))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get tag '%d', got error: %s", data.ID.Value, err),
			)
		}
	} else {
		ztlTag, _, err = d.provider.client.Tags.GetByName(ctx, data.Name.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get tag '%s', got error: %s", data.Name.Value, err),
			)
		}
	}

	if ztlTag != nil {
		diags = resp.State.Set(ctx, tagForState(ztlTag))
		resp.Diagnostics.Append(diags...)
	}
}
