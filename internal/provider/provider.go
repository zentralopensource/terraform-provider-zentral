package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.Provider = &provider{}

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type provider struct {
	configured bool
	version    string
	client     *goztl.Client
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	BaseURL types.String `tfsdk:"base_url"`
	Token   types.String `tfsdk:"token"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// base URL
	var baseURL string
	if data.BaseURL.Unknown {
		resp.Diagnostics.AddWarning(
			"Zentral provider configuration error",
			"Cannot use unknown value as base URL",
		)
		return
	}

	if data.BaseURL.Null {
		baseURL = os.Getenv("ZTL_API_BASE_URL")
	} else {
		baseURL = data.BaseURL.Value
	}

	if baseURL == "" {
		resp.Diagnostics.AddError(
			"Zentral provider configuration error",
			"Base URL cannot be an empty string",
		)
		return
	}

	// base URL
	var token string
	if data.Token.Unknown {
		resp.Diagnostics.AddWarning(
			"Zentral provider configuration error",
			"Cannot use unknown value as token",
		)
		return
	}

	if data.Token.Null {
		token = os.Getenv("ZTL_API_TOKEN")
	} else {
		token = data.Token.Value
	}

	if token == "" {
		resp.Diagnostics.AddError(
			"Zentral provider configuration error",
			"Token cannot be an empty string",
		)
		return
	}

	userAgent := fmt.Sprintf("terraform-provider-zentral/%s", p.version)
	c, err := goztl.NewClient(nil, baseURL, token, goztl.SetUserAgent(userAgent))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create Zentral client:\n\n"+err.Error(),
		)
	}

	p.client = c
	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"zentral_jmespath_check":     jmespathCheckResourceType{},
		"zentral_meta_business_unit": metaBusinessUnitResourceType{},
		"zentral_tag":                tagResourceType{},
		"zentral_taxonomy":           taxonomyResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"zentral_jmespath_check":     jmespathCheckDataSourceType{},
		"zentral_meta_business_unit": metaBusinessUnitDataSourceType{},
		"zentral_tag":                tagDataSourceType{},
		"zentral_taxonomy":           TaxonomyDataSourceType{},
	}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"base_url": {
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "The API base URL.",
				MarkdownDescription: "The base URL where the Zentral API is mounted, including the path. " +
					"Can also be set using the `ZTL_API_BASE_URL` environment variable.",
			},
			"token": {
				Type:      types.StringType,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Description: "The Zentral service account or user token. " +
					"Can also be set using the `ZTL_API_TOKEN` environment variable.",
			},
		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
