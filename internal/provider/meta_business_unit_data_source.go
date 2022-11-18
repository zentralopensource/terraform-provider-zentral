package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MetaBusinessUnitDataSource{}

func NewMetaBusinessUnitDataSource() datasource.DataSource {
	return &MetaBusinessUnitDataSource{}
}

// MetaBusinessUnitDataSource defines the data source implementation.
type MetaBusinessUnitDataSource struct {
	client *goztl.Client
}

func (d *MetaBusinessUnitDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_meta_business_unit"
}

func (d *MetaBusinessUnitDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Allows details of a meta business unit to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `zentral_meta_business_unit` allows details of a meta business unit to be retrieved by its `ID` or name.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description:         "ID of the meta business unit.",
				MarkdownDescription: "ID of the meta business unit.",
				Type:                types.Int64Type,
				Optional:            true,
			},
			"name": {
				Description:         "Name of the meta business unit.",
				MarkdownDescription: "Name of the meta business unit.",
				Type:                types.StringType,
				Optional:            true,
			},
			"api_enrollment_enabled": {
				Description:         "If API enrollments are enabled.",
				MarkdownDescription: "If API enrollments are enabled.",
				Type:                types.BoolType,
				Computed:            true,
			},
		},
	}, nil
}

func (d *MetaBusinessUnitDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MetaBusinessUnitDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
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

func (d *MetaBusinessUnitDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data metaBusinessUnit

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var mbu *goztl.MetaBusinessUnit
	var err error
	if data.ID.Value > 0 {
		mbu, _, err = d.client.MetaBusinessUnits.GetByID(ctx, int(data.ID.Value))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get meta business unit '%d', got error: %s", data.ID.Value, err),
			)
		}
	} else {
		mbu, _, err = d.client.MetaBusinessUnits.GetByName(ctx, data.Name.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get meta business unit '%s', got error: %s", data.Name.Value, err),
			)
		}
	}

	if mbu != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, metaBusinessUnitForState(mbu))...)
	}
}
