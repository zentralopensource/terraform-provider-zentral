package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.DataSourceType = tagDataSourceType{}
var _ tfsdk.DataSource = tagDataSource{}

type tagDataSourceType struct{}

func (t tagDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Tag",
		MarkdownDescription: "Tag",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description:         "ID of the tag",
				MarkdownDescription: "ID of the tag",
				Type:                types.Int64Type,
				Optional:            true,
			},
			"name": {
				Description:         "Name of the tag",
				MarkdownDescription: "Name of the tag",
				Type:                types.StringType,
				Optional:            true,
			},
			"color": {
				Description:         "Color of the tag",
				MarkdownDescription: "Color of the tag",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (t tagDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return tagDataSource{
		provider: provider,
	}, diags
}

type tagDataSource struct {
	provider provider
}

func (d tagDataSource) ValidateConfig(ctx context.Context, req tfsdk.ValidateDataSourceConfigRequest, resp *tfsdk.ValidateDataSourceConfigResponse) {
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

func (d tagDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
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
