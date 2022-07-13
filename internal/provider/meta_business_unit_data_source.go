package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.DataSourceType = metaBusinessUnitDataSourceType{}
var _ tfsdk.DataSource = metaBusinessUnitDataSource{}

type metaBusinessUnitDataSourceType struct{}

func (t metaBusinessUnitDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Meta business unit",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "Name of the meta business unit",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "ID of the meta business unit",
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

type metaBusinessUnitData struct {
	Name types.String `tfsdk:"name"`
	Id   types.Int64  `tfsdk:"id"`
}

type metaBusinessUnitDataSource struct {
	provider provider
}

func (d metaBusinessUnitDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data metaBusinessUnitData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.Value > 0 {
		mbu, _, err := d.provider.client.MetaBusinessUnits.GetByID(ctx, int(data.Id.Value))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get meta business unit '%d', got error: %s", data.Id.Value, err),
			)
		}
		data.Name = types.String{Value: mbu.Name}
	} else {
		mbu, _, err := d.provider.client.MetaBusinessUnits.GetByName(ctx, data.Name.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get meta business unit '%s', got error: %s", data.Name.Value, err),
			)
		}
		data.Id = types.Int64{Value: int64(mbu.ID)}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
