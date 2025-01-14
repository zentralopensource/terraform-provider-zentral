package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &MonolithRepositoryDataSource{}

func NewMonolithRepositoryDataSource() datasource.DataSource {
	return &MonolithRepositoryDataSource{}
}

// MonolithRepositoryDataSource defines the data source implementation.
type MonolithRepositoryDataSource struct {
	client *goztl.Client
}

func (d *MonolithRepositoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monolith_repository"
}

func (d *MonolithRepositoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Monolith repository to be retrieved by its ID or its name.",
		MarkdownDescription: "The data source `zentral_monolith_repository` allows details of a Monolith repository to be retrieved by its `ID` or its `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the repository.",
				MarkdownDescription: "`ID` of the repository.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the repository.",
				MarkdownDescription: "Name of the repository.",
				Optional:            true,
			},
			"meta_business_unit_id": schema.Int64Attribute{
				Description:         "The ID of the meta business unit this repository is restricted to.",
				MarkdownDescription: "The `ID` of the meta business unit this repository is restricted to.",
				Computed:            true,
			},
			"backend": schema.StringAttribute{
				Description:         "Repository backend.",
				MarkdownDescription: "Repository backend.",
				Computed:            true,
			},
			"azure": schema.SingleNestedAttribute{
				Description:         "Azure Blob Storage backend parameters.",
				MarkdownDescription: "Azure Blob Storage backend parameters.",
				Attributes: map[string]schema.Attribute{
					"storage_account": schema.StringAttribute{
						Description:         "Name of the storage account.",
						MarkdownDescription: "Name of the storage account.",
						Computed:            true,
					},
					"container": schema.StringAttribute{
						Description:         "Name of the blob container.",
						MarkdownDescription: "Name of the blob container.",
						Computed:            true,
					},
					"prefix": schema.StringAttribute{
						Description:         "Prefix of the Munki repository in the container.",
						MarkdownDescription: "Prefix of the Munki repository in the container.",
						Computed:            true,
					},
					"client_id": schema.StringAttribute{
						Description:         "Client ID of the Azure app registration.",
						MarkdownDescription: "Client ID of the Azure app registration.",
						Computed:            true,
					},
					"tenant_id": schema.StringAttribute{
						Description:         "Azure tenant ID.",
						MarkdownDescription: "Azure tenant ID.",
						Computed:            true,
					},
					"client_secret": schema.StringAttribute{
						Description:         "Client secret of the Azure app registration.",
						MarkdownDescription: "Client secret of the Azure app registration.",
						Computed:            true,
					},
				},
				Optional: true,
			},
			"s3": schema.SingleNestedAttribute{
				Description:         "S3 backend parameters.",
				MarkdownDescription: "S3 backend parameters.",
				Attributes: map[string]schema.Attribute{
					"bucket": schema.StringAttribute{
						Description:         "Name of the S3 bucket.",
						MarkdownDescription: "Name of the S3 bucket.",
						Computed:            true,
					},
					"region_name": schema.StringAttribute{
						Description:         "Name of the S3 bucket region.",
						MarkdownDescription: "Name of the S3 bucket region.",
						Computed:            true,
					},
					"prefix": schema.StringAttribute{
						Description:         "Prefix of the Munki repository in the S3 bucket.",
						MarkdownDescription: "Prefix of the Munki repository in the S3 bucket.",
						Computed:            true,
					},
					"access_key_id": schema.StringAttribute{
						Description:         "AWS access key ID.",
						MarkdownDescription: "AWS access key ID.",
						Computed:            true,
					},
					"secret_access_key": schema.StringAttribute{
						Description:         "AWS secret access key.",
						MarkdownDescription: "AWS secret access key.",
						Sensitive:           true,
						Computed:            true,
					},
					"assume_role_arn": schema.StringAttribute{
						Description:         "ARN of the IAM role to assume.",
						MarkdownDescription: "ARN of the IAM role to assume.",
						Computed:            true,
					},
					"signature_version": schema.StringAttribute{
						Description:         "Version of the AWS request signature to use.",
						MarkdownDescription: "Version of the AWS request signature to use.",
						Computed:            true,
					},
					"endpoint_url": schema.StringAttribute{
						Description:         "S3 endpoint URL.",
						MarkdownDescription: "S3 endpoint URL.",
						Computed:            true,
					},
					"cloudfront_domain": schema.StringAttribute{
						Description:         "Cloudfront domain.",
						MarkdownDescription: "Cloudfront domain.",
						Computed:            true,
					},
					"cloudfront_key_id": schema.StringAttribute{
						Description:         "Cloudfront key ID.",
						MarkdownDescription: "Cloudfront key ID.",
						Computed:            true,
					},
					"cloudfront_privkey_pem": schema.StringAttribute{
						Description:         "Cloudfront private key in PEM form.",
						MarkdownDescription: "Cloudfront private key in PEM form.",
						Sensitive:           true,
						Computed:            true,
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *MonolithRepositoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonolithRepositoryDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data monolithRepository
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_monolith_repository` data source",
			"`id` or `name` missing",
		)
	} else if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_monolith_repository` data source",
			"`id` and `name` cannot be both set",
		)
	}
}

func (d *MonolithRepositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data monolithRepository

	// Read Terraform repository data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlMR *goztl.MonolithRepository
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlMR, _, err = d.client.MonolithRepositories.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Monolith repository '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	} else {
		ztlMR, _, err = d.client.MonolithRepositories.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Monolith repository '%s' by name, got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if ztlMR != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, monolithRepositoryForState(ztlMR))...)
	}
}
