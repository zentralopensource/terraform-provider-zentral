package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dataschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const (
	tfCertIssuerDigicertBackend        string = "DIGICERT"
	tfCertIssuerIDentBackend                  = "IDENT"
	tfCertIssuerMicrosoftCABackend            = "MICROSOFT_CA"
	tfCertIssuerOktaCABackend                 = "OKTA_CA"
	tfCertIssuerStaticChallengeBackend        = "STATIC_CHALLENGE"
	tfDigicertDefaultAPIBaseURL               = "https://one.digicert.com/mpki/api/"
	tfDigicertDeviceSeatType                  = "DEVICE_SEAT"
	tfDigicertCommonName                      = "common_name"
	tfDigicertEmail                           = "email"
	tfDigicertSerialNumber                    = "serial_number"
	tfDigicertUniqueIdentifier                = "unique_identifier"
	tfDigicertUserIdentifier                  = "user_identifier"
	tfDigicertPseudonym                       = "pseudonym"
	tfDigicertDNQualifier                     = "dn_qualifier"
	tfDigicertRFC882Name                      = "rfc822Name"
	tfDigicertDNSName                         = "dNSName"
	tfDigicertUserSeatType                    = "USER_SEAT"
	tfIDentDefaultRequestTimeout       int64  = 30
	tfIDentDefaultMaxRetries                  = 3
)

// Digicert

var digicertAttrTypes = map[string]attr.Type{
	"api_base_url":       types.StringType,
	"api_token":          types.StringType,
	"profile_guid":       types.StringType,
	"business_unit_guid": types.StringType,
	"seat_type":          types.StringType,
	"seat_id_mapping":    types.StringType,
	"default_seat_email": types.StringType,
}

func digicertBackendForState(digi *goztl.Digicert) types.Object {
	var b types.Object
	if digi != nil {
		b = types.ObjectValueMust(
			digicertAttrTypes,
			map[string]attr.Value{
				"api_base_url":       types.StringValue(digi.APIBaseURL),
				"api_token":          types.StringValue(digi.APIToken),
				"profile_guid":       types.StringValue(digi.ProfileGUID),
				"business_unit_guid": types.StringValue(digi.BusinessUnitGUID),
				"seat_type":          types.StringValue(digi.SeatType),
				"seat_id_mapping":    types.StringValue(digi.SeatIDMapping),
				"default_seat_email": types.StringValue(digi.DefaultSeatEmail),
			},
		)
	} else {
		b = types.ObjectNull(digicertAttrTypes)
	}
	return b
}

func digicertBackendWithState(backend types.Object) *goztl.Digicert {
	var b *goztl.Digicert
	if !backend.IsNull() {
		bMap := backend.Attributes()
		b = &goztl.Digicert{
			APIBaseURL:       bMap["api_base_url"].(types.String).ValueString(),
			APIToken:         bMap["api_token"].(types.String).ValueString(),
			ProfileGUID:      bMap["profile_guid"].(types.String).ValueString(),
			BusinessUnitGUID: bMap["business_unit_guid"].(types.String).ValueString(),
			SeatType:         bMap["seat_type"].(types.String).ValueString(),
			SeatIDMapping:    bMap["seat_id_mapping"].(types.String).ValueString(),
			DefaultSeatEmail: bMap["default_seat_email"].(types.String).ValueString(),
		}
	}
	return b
}

func makeDigicertBackendDataSourceAttribute() dataschema.SingleNestedAttribute {
	desc := "Digicert backend parameters."
	return dataschema.SingleNestedAttribute{
		Description:         desc,
		MarkdownDescription: desc,
		Attributes: map[string]dataschema.Attribute{
			"api_base_url": dataschema.StringAttribute{
				Description:         "API base URL.",
				MarkdownDescription: "API base URL.",
				Computed:            true,
			},
			"api_token": dataschema.StringAttribute{
				Description:         "API token.",
				MarkdownDescription: "API token.",
				Sensitive:           true,
				Computed:            true,
			},
			"profile_guid": dataschema.StringAttribute{
				Description:         "Profile GUID.",
				MarkdownDescription: "Profile GUID.",
				Computed:            true,
			},
			"business_unit_guid": dataschema.StringAttribute{
				Description:         "Business unit GUID.",
				MarkdownDescription: "Business unit GUID.",
				Computed:            true,
			},
			"seat_type": dataschema.StringAttribute{
				Description:         "Seat type.",
				MarkdownDescription: "Seat type.",
				Computed:            true,
			},
			"seat_id_mapping": dataschema.StringAttribute{
				Description:         "Seat ID mapping.",
				MarkdownDescription: "Seat ID mapping.",
				Computed:            true,
			},
			"default_seat_email": dataschema.StringAttribute{
				Description:         "Default seat email.",
				MarkdownDescription: "Default seat email.",
				Computed:            true,
			},
		},
		Computed: true,
	}
}

func makeDigicertBackendResourceAttribute() schema.SingleNestedAttribute {
	desc := "IDent backend parameters."
	return schema.SingleNestedAttribute{
		Description:         desc,
		MarkdownDescription: desc,
		Attributes: map[string]schema.Attribute{
			"api_base_url": schema.StringAttribute{
				Description:         "API base URL. Defaults to https://one.digicert.com/mpki/api/.",
				MarkdownDescription: "API base URL. Defaults to `https://one.digicert.com/mpki/api/`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(tfDigicertDefaultAPIBaseURL),
			},
			"api_token": schema.StringAttribute{
				Description:         "API token.",
				MarkdownDescription: "API token.",
				Sensitive:           true,
				Required:            true,
			},
			"profile_guid": schema.StringAttribute{
				Description:         "Profile GUID.",
				MarkdownDescription: "Profile GUID.",
				Required:            true,
			},
			"business_unit_guid": schema.StringAttribute{
				Description:         "Business unit GUID.",
				MarkdownDescription: "Business unit GUID.",
				Required:            true,
			},
			"seat_type": schema.StringAttribute{
				Description:         "Seat type. DEVICE_SEAT or USER_SEAT. Defaults to DEVICE_SEAT.",
				MarkdownDescription: "Seat type. `DEVICE_SEAT` or `USER_SEAT`. Defaults to `DEVICE_SEAT`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(tfDigicertDeviceSeatType),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{tfDigicertDeviceSeatType, tfDigicertUserSeatType}...),
				},
			},
			"seat_id_mapping": schema.StringAttribute{
				Description:         "Seat ID mapping. Possible values: common_name, email, serial_number, unique_identifier, user_identifier, pseudonym, dn_qualifier, rfc822Name, dNSName. Defaults to common_name.",
				MarkdownDescription: "Seat ID mapping. Possible values: `common_name`, `email`, `serial_number`, `unique_identifier`, `user_identifier`, `pseudonym`, `dn_qualifier`, `rfc822Name`, `dNSName`. Defaults to `common_name`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(tfDigicertCommonName),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						tfDigicertCommonName,
						tfDigicertEmail,
						tfDigicertSerialNumber,
						tfDigicertUniqueIdentifier,
						tfDigicertUserIdentifier,
						tfDigicertPseudonym,
						tfDigicertDNQualifier,
						tfDigicertRFC882Name,
						tfDigicertDNSName,
					}...),
				},
			},
			"default_seat_email": schema.StringAttribute{
				Description:         "Default seat email.",
				MarkdownDescription: "Default seat email.",
				Required:            true,
			},
		},
		Optional: true,
	}
}

// IDent

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
