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
var _ datasource.DataSource = &OsqueryFileCategoryDataSource{}

func NewOsqueryFileCategoryDataSource() datasource.DataSource {
	return &OsqueryFileCategoryDataSource{}
}

// OsqueryFileCategoryDataSource defines the data source implementation.
type OsqueryFileCategoryDataSource struct {
	client *goztl.Client
}

func (d *OsqueryFileCategoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_osquery_file_category"
}

func (d *OsqueryFileCategoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Osquery file category to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `zentral_osquery_file_category` allows details of a Osquery file category to be retrieved by its `ID` or name.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Osquery file category.",
				MarkdownDescription: "`ID` of the Osquery file category.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the Osquery file category.",
				MarkdownDescription: "Name of the Osquery file category.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the Osquery file category.",
				MarkdownDescription: "Description of the Osquery file category.",
				Computed:            true,
			},
			"file_paths": schema.SetAttribute{
				Description:         "Set of paths to include in the Osquery file category.",
				MarkdownDescription: "Set of paths to include in the Osquery file category.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"exclude_paths": schema.SetAttribute{
				Description:         "Set of paths to exclude from the Osquery file category.",
				MarkdownDescription: "Set of paths to exclude from the Osquery file category.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"file_paths_queries": schema.SetAttribute{
				Description:         "Set of queries returning paths to monitor as path columns in the results.",
				MarkdownDescription: "Set of queries returning paths to monitor as path columns in the results.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"access_monitoring": schema.BoolAttribute{
				Description:         "If true, FIM will include file access for this file category. Defaults to false.",
				MarkdownDescription: "If `true`, FIM will include file access for this file category. Defaults to `false`.",
				Computed:            true,
			},
		},
	}
}

func (d *OsqueryFileCategoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OsqueryFileCategoryDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data osqueryFileCategory
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_osquery_file_category` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_osquery_file_category` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *OsqueryFileCategoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data osqueryFileCategory

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlOFC *goztl.OsqueryFileCategory
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlOFC, _, err = d.client.OsqueryFileCategories.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Osquery file category '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlOFC, _, err = d.client.OsqueryFileCategories.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Osquery file category '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlOFC != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, osqueryFileCategoryForState(ztlOFC))...)
	}
}
