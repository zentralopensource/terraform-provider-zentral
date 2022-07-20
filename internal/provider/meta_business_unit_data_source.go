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
var _ tfsdk.DataSourceType = metaBusinessUnitDataSourceType{}
var _ tfsdk.DataSource = metaBusinessUnitDataSource{}

type metaBusinessUnitDataSourceType struct{}

func (t metaBusinessUnitDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Allows details of a meta business unit to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `zentral_meta_business_unit` allows details of a meta business unit to be retrieved by its `ID` or name.",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Description:         "Name of the meta business unit.",
				MarkdownDescription: "Name of the meta business unit.",
				Type:                types.StringType,
				Optional:            true,
			},
			"id": {
				Description:         "ID of the meta business unit.",
				MarkdownDescription: "ID of the meta business unit.",
				Type:                types.Int64Type,
				Optional:            true,
			},
		},
	}, nil
}

func (t metaBusinessUnitDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return metaBusinessUnitDataSource{
		provider: provider,
	}, diags
}

type metaBusinessUnitDataSource struct {
	provider provider
}

func (d metaBusinessUnitDataSource) ValidateConfig(ctx context.Context, req tfsdk.ValidateDataSourceConfigRequest, resp *tfsdk.ValidateDataSourceConfigResponse) {
	var data metaBusinessUnit
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.Null && data.Name.Null {
		resp.Diagnostics.AddError(
			"Invalid `zentral_meta_business_unit` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.Null && !data.Name.Null {
		resp.Diagnostics.AddError(
			"Invalid `zentral_meta_business_unit` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d metaBusinessUnitDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data metaBusinessUnit

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var mbu *goztl.MetaBusinessUnit
	var err error
	if data.ID.Value > 0 {
		mbu, _, err = d.provider.client.MetaBusinessUnits.GetByID(ctx, int(data.ID.Value))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get meta business unit '%d', got error: %s", data.ID.Value, err),
			)
		}
	} else {
		mbu, _, err = d.provider.client.MetaBusinessUnits.GetByName(ctx, data.Name.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get meta business unit '%s', got error: %s", data.Name.Value, err),
			)
		}
	}

	if mbu != nil {
		diags = resp.State.Set(ctx, metaBusinessUnitForState(mbu))
		resp.Diagnostics.Append(diags...)
	}
}
