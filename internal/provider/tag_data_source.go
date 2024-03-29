package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &TagDataSource{}

func NewTagDataSource() datasource.DataSource {
	return &TagDataSource{}
}

// TagDataSource defines the data source implementation.
type TagDataSource struct {
	client *goztl.Client
}

func (d *TagDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (d *TagDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a tag to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `zentral_tag` allows details of a tag to be retrieved by its `ID` or name.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the tag.",
				MarkdownDescription: "`ID` of the tag.",
				Optional:            true,
			},
			"taxonomy_id": schema.Int64Attribute{
				Description:         "ID of the tag taxonomy.",
				MarkdownDescription: "`ID` of the tag taxonomy.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the tag.",
				MarkdownDescription: "Name of the tag.",
				Optional:            true,
			},
			"color": schema.StringAttribute{
				Description:         "Color of the tag.",
				MarkdownDescription: "Color of the tag.",
				Computed:            true,
			},
		},
	}
}

func (d *TagDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TagDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data tag
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_tag` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_tag` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *TagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data tag

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlTag *goztl.Tag
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlTag, _, err = d.client.Tags.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get tag '%d', got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlTag, _, err = d.client.Tags.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get tag '%s', got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlTag != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, tagForState(ztlTag))...)
	}
}
