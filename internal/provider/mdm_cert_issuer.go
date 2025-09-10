package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dataschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const (
	tfCertIssuerIDentBackend           string = "IDENT"
	tfCertIssuerMicrosoftCABackend            = "MICROSOFT_CA"
	tfCertIssuerOktaCABackend                 = "OKTA_CA"
	tfCertIssuerStaticChallengeBackend        = "STATIC_CHALLENGE"
	tfIDentDefaultRequestTimeout       int64  = 30
	tfIDentDefaultMaxRetries                  = 3
)

// IDente

var identAttrTypes = map[string]attr.Type{
	"url":             types.StringType,
	"bearer_token":    types.StringType,
	"request_timeout": types.Int64Type,
	"max_retries":     types.Int64Type,
}

func identBackendForState(ident *goztl.IDent) types.Object {
	var b types.Object
	if ident != nil {
		b = types.ObjectValueMust(
			identAttrTypes,
			map[string]attr.Value{
				"url":             types.StringValue(ident.URL),
				"bearer_token":    types.StringValue(ident.BearerToken),
				"request_timeout": types.Int64Value(int64(ident.RequestTimeout)),
				"max_retries":     types.Int64Value(int64(ident.MaxRetries)),
			},
		)
	} else {
		b = types.ObjectNull(identAttrTypes)
	}
	return b
}

func identBackendWithState(backend types.Object) *goztl.IDent {
	var b *goztl.IDent
	if !backend.IsNull() {
		bMap := backend.Attributes()
		b = &goztl.IDent{
			URL:            bMap["url"].(types.String).ValueString(),
			BearerToken:    bMap["bearer_token"].(types.String).ValueString(),
			RequestTimeout: int(bMap["request_timeout"].(types.Int64).ValueInt64()),
			MaxRetries:     int(bMap["max_retries"].(types.Int64).ValueInt64()),
		}
	}
	return b
}

func makeIDentBackendDataSourceAttribute() dataschema.SingleNestedAttribute {
	desc := "IDent backend parameters."
	return dataschema.SingleNestedAttribute{
		Description:         desc,
		MarkdownDescription: desc,
		Attributes: map[string]dataschema.Attribute{
			"url": dataschema.StringAttribute{
				Description:         "HTTP endpoint URL.",
				MarkdownDescription: "HTTP endpoint URL.",
				Computed:            true,
			},
			"bearer_token": dataschema.StringAttribute{
				Description:         "Bearer token.",
				MarkdownDescription: "Bearer token.",
				Sensitive:           true,
				Computed:            true,
			},
			"request_timeout": dataschema.Int64Attribute{
				Description:         "Request timeout.",
				MarkdownDescription: "Request timeout.",
				Computed:            true,
			},
			"max_retries": dataschema.Int64Attribute{
				Description:         "Max number of retries.",
				MarkdownDescription: "Max number of retries.",
				Computed:            true,
			},
		},
		Computed: true,
	}
}

func makeIDentBackendResourceAttribute() schema.SingleNestedAttribute {
	desc := "IDent backend parameters."
	return schema.SingleNestedAttribute{
		Description:         desc,
		MarkdownDescription: desc,
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description:         "HTTP endpoint URL.",
				MarkdownDescription: "HTTP endpoint URL.",
				Required:            true,
			},
			"bearer_token": schema.StringAttribute{
				Description:         "Bearer token.",
				MarkdownDescription: "Bearer token.",
				Sensitive:           true,
				Required:            true,
			},
			"request_timeout": schema.Int64Attribute{
				Description:         "Request timeout.",
				MarkdownDescription: "Request timeout.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(tfIDentDefaultRequestTimeout),
			},
			"max_retries": schema.Int64Attribute{
				Description:         "Max number of retries.",
				MarkdownDescription: "Max number of retries.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(tfIDentDefaultMaxRetries),
			},
		},
		Optional: true,
	}
}

// Microsoft CA

var microsoftCAAttrTypes = map[string]attr.Type{
	"url":      types.StringType,
	"username": types.StringType,
	"password": types.StringType,
}

func microsoftCABackendForState(mc *goztl.MicrosoftCA) types.Object {
	var b types.Object
	if mc != nil {
		b = types.ObjectValueMust(
			microsoftCAAttrTypes,
			map[string]attr.Value{
				"url":      types.StringValue(mc.URL),
				"username": types.StringValue(mc.Username),
				"password": types.StringValue(mc.Password),
			},
		)
	} else {
		b = types.ObjectNull(microsoftCAAttrTypes)
	}
	return b
}

func microsoftCABackendWithState(backend types.Object) *goztl.MicrosoftCA {
	var b *goztl.MicrosoftCA
	if !backend.IsNull() {
		bMap := backend.Attributes()
		b = &goztl.MicrosoftCA{
			URL:      bMap["url"].(types.String).ValueString(),
			Username: bMap["username"].(types.String).ValueString(),
			Password: bMap["password"].(types.String).ValueString(),
		}
	}
	return b
}

func makeMicrosoftCABackendDataSourceAttribute(desc string) dataschema.SingleNestedAttribute {
	desc = fmt.Sprintf("%s backend parameters.", desc)
	return dataschema.SingleNestedAttribute{
		Description:         desc,
		MarkdownDescription: desc,
		Attributes: map[string]dataschema.Attribute{
			"url": schema.StringAttribute{
				Description:         "HTTP endpoint URL.",
				MarkdownDescription: "HTTP endpoint URL.",
				Computed:            true,
			},
			"username": dataschema.StringAttribute{
				Description:         "Username.",
				MarkdownDescription: "Username.",
				Computed:            true,
			},
			"password": dataschema.StringAttribute{
				Description:         "Password.",
				MarkdownDescription: "Password.",
				Sensitive:           true,
				Computed:            true,
			},
		},
		Computed: true,
	}
}

func makeMicrosoftCABackendResourceAttribute(desc string) schema.SingleNestedAttribute {
	desc = fmt.Sprintf("%s backend parameters.", desc)
	return schema.SingleNestedAttribute{
		Description:         desc,
		MarkdownDescription: desc,
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description:         "HTTP endpoint URL.",
				MarkdownDescription: "HTTP endpoint URL.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				Description:         "Username.",
				MarkdownDescription: "Username.",
				Required:            true,
			},
			"password": schema.StringAttribute{
				Description:         "Password.",
				MarkdownDescription: "Password.",
				Sensitive:           true,
				Required:            true,
			},
		},
		Optional: true,
	}
}

// Static Challenge

var staticChallengeAttrTypes = map[string]attr.Type{
	"challenge": types.StringType,
}

func staticChallengeBackendForState(sc *goztl.StaticChallenge) types.Object {
	var b types.Object
	if sc != nil {
		b = types.ObjectValueMust(
			staticChallengeAttrTypes,
			map[string]attr.Value{
				"challenge": types.StringValue(sc.Challenge),
			},
		)
	} else {
		b = types.ObjectNull(staticChallengeAttrTypes)
	}
	return b
}

func staticChallengeBackendWithState(backend types.Object) *goztl.StaticChallenge {
	var b *goztl.StaticChallenge
	if !backend.IsNull() {
		bMap := backend.Attributes()
		b = &goztl.StaticChallenge{
			Challenge: bMap["challenge"].(types.String).ValueString(),
		}
	}
	return b
}

func makeStaticChallengeBackendDataSourceAttribute() dataschema.SingleNestedAttribute {
	desc := "Static Challenge backend parameters."
	return dataschema.SingleNestedAttribute{
		Description:         desc,
		MarkdownDescription: desc,
		Attributes: map[string]dataschema.Attribute{
			"challenge": dataschema.StringAttribute{
				Description:         "Challenge.",
				MarkdownDescription: "Challenge.",
				Sensitive:           true,
				Computed:            true,
			},
		},
		Computed: true,
	}
}

func makeStaticChallengeBackendResourceAttribute() schema.SingleNestedAttribute {
	desc := "Static Challenge backend parameters."
	return schema.SingleNestedAttribute{
		Description:         desc,
		MarkdownDescription: desc,
		Attributes: map[string]schema.Attribute{
			"challenge": schema.StringAttribute{
				Description:         "Challenge.",
				MarkdownDescription: "Challenge.",
				Sensitive:           true,
				Required:            true,
			},
		},
		Optional: true,
	}
}
