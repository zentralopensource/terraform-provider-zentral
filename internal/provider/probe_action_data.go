package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &ProbeActionDataSource{}

func NewProbeActionDataSource() datasource.DataSource {
	return &ProbeActionDataSource{}
}

// ProbeActionDataSource defines the data source implementation.
type ProbeActionDataSource struct {
	client *goztl.Client
}

func (r *ProbeActionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_probe_action"
}

func (r *ProbeActionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a probe action to be retrieved by its ID or its name.",
		MarkdownDescription: "The data source `zentral_probe_action` allows details of a probe action to be retrieved by its `ID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the action.",
				MarkdownDescription: "`ID` of the action.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the action.",
				MarkdownDescription: "Name of the action.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the action.",
				MarkdownDescription: "Description of the action.",
				Computed:            true,
			},
			"backend": schema.StringAttribute{
				Description:         "Action backend.",
				MarkdownDescription: "Action backend.",
				Computed:            true,
			},
			"http_post": schema.SingleNestedAttribute{
				Description:         "HTTP Post backend parameters.",
				MarkdownDescription: "HTTP Post backend parameters.",
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Description:         "URL.",
						MarkdownDescription: "`URL`.",
						Computed:            true,
					},
					"username": schema.StringAttribute{
						Description:         "Username for basic authentication.",
						MarkdownDescription: "Username for basic authentication.",
						Computed:            true,
					},
					"password": schema.StringAttribute{
						Description:         "Password for basic authentication.",
						MarkdownDescription: "Password for basic authentication.",
						Sensitive:           true,
						Computed:            true,
					},
					"headers": schema.SetNestedAttribute{
						Description:         "A set of additional HTTP headers to add to the requests.",
						MarkdownDescription: "A set of additional HTTP headers to add to the requests.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description:         "Name of the HTTP header.",
									MarkdownDescription: "Name of the HTTP header.",
									Computed:            true,
								},
								"value": schema.StringAttribute{
									Description:         "Value of the HTTP header.",
									MarkdownDescription: "Value of the HTTP header.",
									Computed:            true,
									Sensitive:           true,
								},
							},
						},
						Computed: true,
					},
					"cel_transformation": schema.StringAttribute{
						Description:         "CEL expression that is used to transform the event data. The input to the expression is a Map with two keys: metadata for the event metadata and payload for the event payload.",
						MarkdownDescription: "CEL expression that is used to transform the event data. The input to the expression is a `Map` with two keys: `metadata` for the event metadata and `payload` for the event payload.",
						Computed:            true,
					},
				},
				Computed: true,
			},
			"slack_incoming_webhook": schema.SingleNestedAttribute{
				Description:         "Slack incoming webhook backend parameters.",
				MarkdownDescription: "Slack incoming webhook backend parameters.",
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Description:         "URL.",
						MarkdownDescription: "`URL`.",
						Computed:            true,
						Sensitive:           true,
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *ProbeActionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*goztl.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *goztl.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ProbeActionDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data probeAction
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_probe_action` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_probe_action` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *ProbeActionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data probeAction

	// Read Terraform probe action data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlPA *goztl.ProbeAction
	var err error
	if len(data.ID.ValueString()) > 0 {
		ztlPA, _, err = d.client.ProbesActions.GetByID(ctx, data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get probe action '%s' by ID, got error: %s", data.ID.ValueString(), err),
			)
		}
	} else {
		ztlPA, _, err = d.client.ProbesActions.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get probe action '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlPA != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, probeActionForState(ztlPA))...)
	}
}
