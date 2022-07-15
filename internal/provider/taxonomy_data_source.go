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
var _ tfsdk.DataSourceType = TaxonomyDataSourceType{}
var _ tfsdk.DataSource = TaxonomyDataSource{}

type TaxonomyDataSourceType struct{}

func (t TaxonomyDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Taxonomy",
		MarkdownDescription: "Taxonomy",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description:         "ID of the Taxonomy",
				MarkdownDescription: "ID of the Taxonomy",
				Type:                types.Int64Type,
				Optional:            true,
			},
			"name": {
				Description:         "Name of the Taxonomy",
				MarkdownDescription: "Name of the Taxonomy",
				Type:                types.StringType,
				Optional:            true,
			},
		},
	}, nil
}

func (t TaxonomyDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return TaxonomyDataSource{
		provider: provider,
	}, diags
}

type TaxonomyDataSource struct {
	provider provider
}

func (d TaxonomyDataSource) ValidateConfig(ctx context.Context, req tfsdk.ValidateDataSourceConfigRequest, resp *tfsdk.ValidateDataSourceConfigResponse) {
	var data taxonomy
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.Null && data.Name.Null {
		resp.Diagnostics.AddError(
			"Invalid `zentral_taxonomy` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.Null && !data.Name.Null {
		resp.Diagnostics.AddError(
			"Invalid `zentral_taxonomy` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d TaxonomyDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data taxonomy

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ztlTaxonomy *goztl.Taxonomy
	var err error
	if data.ID.Value > 0 {
		ztlTaxonomy, _, err = d.provider.client.Taxonomies.GetByID(ctx, int(data.ID.Value))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get taxonomy '%d', got error: %s", data.ID.Value, err),
			)
		}
	} else {
		ztlTaxonomy, _, err = d.provider.client.Taxonomies.GetByName(ctx, data.Name.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get taxonomy '%s', got error: %s", data.Name.Value, err),
			)
		}
	}

	if ztlTaxonomy != nil {
		diags = resp.State.Set(ctx, taxonomyForState(ztlTaxonomy))
		resp.Diagnostics.Append(diags...)
	}
}
