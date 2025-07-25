package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.Provider = &ZentralProvider{}

// ZentralProvider defines the provider implementation.
type ZentralProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ZentralProviderModel describes the provider data model.
type ZentralProviderModel struct {
	BaseURL types.String `tfsdk:"base_url"`
	Token   types.String `tfsdk:"token"`
}

func (p *ZentralProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "zentral"
	resp.Version = p.version
}

func (p *ZentralProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Optional:    true,
				Description: "The API base URL.",
				MarkdownDescription: "The base URL where the Zentral API is mounted, including the path. " +
					"Can also be set using the `ZTL_API_BASE_URL` environment variable.",
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Description: "The Zentral service account or user token. " +
					"Can also be set using the `ZTL_API_TOKEN` environment variable.",
			},
		},
	}
}

func (p *ZentralProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ZentralProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// base URL
	var baseURL string
	if data.BaseURL.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Zentral provider configuration error",
			"Cannot use unknown value as base URL",
		)
		return
	}

	if data.BaseURL.IsNull() {
		baseURL = os.Getenv("ZTL_API_BASE_URL")
	} else {
		baseURL = data.BaseURL.ValueString()
	}

	if baseURL == "" {
		resp.Diagnostics.AddError(
			"Zentral provider configuration error",
			"Base URL cannot be an empty string",
		)
		return
	}

	// API token
	var token string
	if data.Token.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Zentral provider configuration error",
			"Cannot use unknown value as token",
		)
		return
	}

	if data.Token.IsNull() {
		token = os.Getenv("ZTL_API_TOKEN")
	} else {
		token = data.Token.ValueString()
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

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *ZentralProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewJMESPathCheckResource,
		NewMetaBusinessUnitResource,
		NewMDMArtifactResource,
		NewMDMBlueprintResource,
		NewMDMBlueprintArtifactResource,
		NewMDMDataAssetResource,
		NewMDMDeclarationResource,
		NewMDMEnterpriseAppResource,
		NewMDMFileVaultConfigResource,
		NewMDMOTAEnrollmentResource,
		NewMDMProfileResource,
		NewMDMRecoveryPasswordConfigResource,
		NewMDMSoftwareUpdateEnforcementResource,
		NewMonolithCatalogResource,
		NewMonolithConditionResource,
		NewMonolithEnrollmentResource,
		NewMonolithManifestResource,
		NewMonolithManifestCatalogResource,
		NewMonolithManifestEnrollmentPackageResource,
		NewMonolithManifestSubManifestResource,
		NewMonolithRepositoryResource,
		NewMonolithSubManifestResource,
		NewMonolithSubManifestPkgInfoResource,
		NewMunkiConfigurationResource,
		NewMunkiEnrollmentResource,
		NewMunkiScriptCheckResource,
		NewOsqueryATCResource,
		NewOsqueryConfigurationResource,
		NewOsqueryConfigurationPackResource,
		NewOsqueryEnrollmentResource,
		NewOsqueryFileCategoryResource,
		NewOsqueryPackResource,
		NewOsqueryQueryResource,
		NewProbeResource,
		NewProbeActionResource,
		NewSantaConfigurationResource,
		NewSantaEnrollmentResource,
		NewSantaRuleResource,
		NewStoreResource,
		NewTagResource,
		NewTaxonomyResource,
	}
}

func (p *ZentralProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewJMESPathCheckDataSource,
		NewMetaBusinessUnitDataSource,
		NewMDMArtifactDataSource,
		NewMDMBlueprintDataSource,
		NewMDMFileVaultConfigDataSource,
		NewMDMOTAEnrollmentDataSource,
		NewMDMPushCertificateDataSource,
		NewMDMRecoveryPasswordConfigDataSource,
		NewMDMSCEPConfigDataSource,
		NewMDMSoftwareUpdateEnforcementDataSource,
		NewMonolithCatalogDataSource,
		NewMonolithConditionDataSource,
		NewMonolithEnrollmentDataSource,
		NewMonolithManifestDataSource,
		NewMonolithRepositoryDataSource,
		NewMonolithSubManifestDataSource,
		NewMunkiConfigurationDataSource,
		NewMunkiEnrollmentDataSource,
		NewMunkiScriptCheckDataSource,
		NewOsqueryATCDataSource,
		NewOsqueryConfigurationDataSource,
		NewOsqueryEnrollmentDataSource,
		NewOsqueryFileCategoryDataSource,
		NewOsqueryPackDataSource,
		NewOsqueryQueryDataSource,
		NewProbeActionDataSource,
		NewRealmsRealmDataSource,
		NewSantaConfigurationDataSource,
		NewSantaEnrollmentDataSource,
		NewSantaRuleDataSource,
		NewTagDataSource,
		NewTaxonomyDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ZentralProvider{
			version: version,
		}
	}
}
