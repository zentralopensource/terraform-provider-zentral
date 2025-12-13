package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const (
	tfStoreHTTPBackend                          string = "HTTP"
	tfStoreKinesisBackend                              = "KINESIS"
	tfStorePantherBackend                              = "PANTHER"
	tfStoreSplunkBackend                               = "SPLUNK"
	tfStoreHTTPBackendDefaultConcurrency        int64  = 1
	tfStoreHTTPBackendMinConcurrency                   = 1
	tfStoreHTTPBackendDefaultRequestTimeout            = 120
	tfStoreHTTPBackendMinRequestTimeout                = 1
	tfStoreHTTPBackendMaxRequestTimeout                = 600
	tfStoreHTTPBackendDefaultMaxRetries                = 3
	tfStoreHTTPBackendMinMaxRetries                    = 1
	tfStoreHTTPBackendMaxMaxRetries                    = 5
	tfStoreKinesisBackendDefaultBatchSize              = 1
	tfStoreKinesisBackendMinBatchSize                  = 1
	tfStoreKinesisBackendMaxBatchSize                  = 500
	tfStoreKinesisSerializationFormatZentral           = "zentral"
	tfStoreKinesisSerializationFormatFirehoseV1        = "firehose_v1"
	tfStorePantherBackendDefaultBatchSize              = 1
	tfStorePantherBackendMinBatchSize                  = 1
	tfStorePantherBackendMaxBatchSize                  = 100
	tfStoreSplunkBackendDefaultTimeout                 = 300
	tfStoreSplunkBackendDefaultBatchSize               = 1
	tfStoreSplunkBackendMinBatchSize                   = 1
	tfStoreSplunkBackendMaxBatchSize                   = 100
)

type store struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Description                types.String `tfsdk:"description"`
	AdminConsole               types.Bool   `tfsdk:"admin_console"`
	EventsURLAuthorizedRoleIDs types.Set    `tfsdk:"events_url_authorized_role_ids"`
	EventFilters               types.Object `tfsdk:"event_filters"`
	Backend                    types.String `tfsdk:"backend"`
	HTTP                       types.Object `tfsdk:"http"`
	Kinesis                    types.Object `tfsdk:"kinesis"`
	Panther                    types.Object `tfsdk:"panther"`
	Splunk                     types.Object `tfsdk:"splunk"`
}

var storeHTTPAttrTypes = map[string]attr.Type{
	"endpoint_url":    types.StringType,
	"verify_tls":      types.BoolType,
	"username":        types.StringType,
	"password":        types.StringType,
	"headers":         types.SetType{ElemType: types.ObjectType{AttrTypes: headerAttrTypes}},
	"concurrency":     types.Int64Type,
	"request_timeout": types.Int64Type,
	"max_retries":     types.Int64Type,
}

var storeKinesisAttrTypes = map[string]attr.Type{
	"region_name":           types.StringType,
	"aws_access_key_id":     types.StringType,
	"aws_secret_access_key": types.StringType,
	"assume_role_arn":       types.StringType,
	"stream":                types.StringType,
	"batch_size":            types.Int64Type,
	"serialization_format":  types.StringType,
}

var storePantherAttrTypes = map[string]attr.Type{
	"endpoint_url": types.StringType,
	"bearer_token": types.StringType,
	"batch_size":   types.Int64Type,
}

var storeSplunkAttrTypes = map[string]attr.Type{
	// HEC
	"hec_url":                       types.StringType,
	"hec_token":                     types.StringType,
	"hec_extra_headers":             types.SetType{ElemType: types.ObjectType{AttrTypes: headerAttrTypes}},
	"hec_request_timeout":           types.Int64Type,
	"hec_index":                     types.StringType,
	"hec_source":                    types.StringType,
	"computer_name_as_host_sources": types.ListType{ElemType: types.StringType},
	"custom_host_field":             types.StringType,
	"serial_number_field":           types.StringType,
	"batch_size":                    types.Int64Type,
	// Events URLs
	"search_app_url": types.StringType,
	// Events search
	"search_url":             types.StringType,
	"search_token":           types.StringType,
	"search_extra_headers":   types.SetType{ElemType: types.ObjectType{AttrTypes: headerAttrTypes}},
	"search_request_timeout": types.Int64Type,
	"search_index":           types.StringType,
	"search_source":          types.StringType,
	// Common
	"verify_tls": types.BoolType,
}

func httpBackendForState(s *goztl.Store) types.Object {
	var b types.Object
	if s.HTTP != nil {
		b = types.ObjectValueMust(
			storeHTTPAttrTypes,
			map[string]attr.Value{
				"endpoint_url":    types.StringValue(s.HTTP.EndpointURL),
				"verify_tls":      types.BoolValue(s.HTTP.VerifyTLS),
				"username":        optionalStringForState(s.HTTP.Username),
				"password":        optionalStringForState(s.HTTP.Password),
				"headers":         headersForState(s.HTTP.Headers),
				"concurrency":     types.Int64Value(int64(s.HTTP.Concurrency)),
				"request_timeout": types.Int64Value(int64(s.HTTP.RequestTimeout)),
				"max_retries":     types.Int64Value(int64(s.HTTP.MaxRetries)),
			},
		)
	} else {
		b = types.ObjectNull(storeHTTPAttrTypes)
	}
	return b
}

func kinesisBackendForState(s *goztl.Store) types.Object {
	var b types.Object
	if s.Kinesis != nil {
		b = types.ObjectValueMust(
			storeKinesisAttrTypes,
			map[string]attr.Value{
				"region_name":           types.StringValue(s.Kinesis.RegionName),
				"aws_access_key_id":     optionalStringForState(s.Kinesis.AWSAccessKeyID),
				"aws_secret_access_key": optionalStringForState(s.Kinesis.AWSSecretAccessKey),
				"assume_role_arn":       optionalStringForState(s.Kinesis.AssumeRoleARN),
				"stream":                types.StringValue(s.Kinesis.Stream),
				"batch_size":            types.Int64Value(int64(s.Kinesis.BatchSize)),
				"serialization_format":  types.StringValue(s.Kinesis.SerializationFormat),
			},
		)
	} else {
		b = types.ObjectNull(storeKinesisAttrTypes)
	}
	return b
}

func pantherBackendForState(s *goztl.Store) types.Object {
	var b types.Object
	if s.Panther != nil {
		b = types.ObjectValueMust(
			storePantherAttrTypes,
			map[string]attr.Value{
				"endpoint_url": types.StringValue(s.Panther.EndpointURL),
				"bearer_token": types.StringValue(s.Panther.BearerToken),
				"batch_size":   types.Int64Value(int64(s.Panther.BatchSize)),
			},
		)
	} else {
		b = types.ObjectNull(storePantherAttrTypes)
	}
	return b
}

func splunkBackendForState(s *goztl.Store) types.Object {
	var b types.Object
	if s.Splunk != nil {
		b = types.ObjectValueMust(
			storeSplunkAttrTypes,
			map[string]attr.Value{
				// HEC
				"hec_url":                       types.StringValue(s.Splunk.HECURL),
				"hec_token":                     types.StringValue(s.Splunk.HECToken),
				"hec_extra_headers":             headersForState(s.Splunk.HECExtraHeaders),
				"hec_request_timeout":           types.Int64Value(int64(s.Splunk.HECRequestTimeout)),
				"hec_index":                     optionalStringForState(s.Splunk.HECIndex),
				"hec_source":                    optionalStringForState(s.Splunk.HECSource),
				"computer_name_as_host_sources": stringListForState(s.Splunk.ComputerNameAsHostSources),
				"custom_host_field":             optionalStringForState(s.Splunk.CustomHostField),
				"serial_number_field":           types.StringValue(s.Splunk.SerialNumberField),
				"batch_size":                    types.Int64Value(int64(s.Splunk.BatchSize)),
				// Events URLs
				"search_app_url": optionalStringForState(s.Splunk.SearchAppURL),
				// Events search
				"search_url":             optionalStringForState(s.Splunk.SearchURL),
				"search_token":           optionalStringForState(s.Splunk.SearchToken),
				"search_extra_headers":   headersForState(s.Splunk.SearchExtraHeaders),
				"search_request_timeout": types.Int64Value(int64(s.Splunk.SearchRequestTimeout)),
				"search_index":           optionalStringForState(s.Splunk.SearchIndex),
				"search_source":          optionalStringForState(s.Splunk.SearchSource),
				// Common
				"verify_tls": types.BoolValue(s.Splunk.VerifyTLS),
			},
		)
	} else {
		b = types.ObjectNull(storeSplunkAttrTypes)
	}
	return b
}

func storeForState(s *goztl.Store) store {
	return store{
		ID:                         types.StringValue(s.ID),
		Name:                       types.StringValue(s.Name),
		Description:                types.StringValue(s.Description),
		AdminConsole:               types.BoolValue(s.AdminConsole),
		EventsURLAuthorizedRoleIDs: int64SetForState(s.EventsURLAuthorizedRoleIDs),
		EventFilters:               eventFilterSetForState(s.EventFilters),
		Backend:                    types.StringValue(s.Backend),
		HTTP:                       httpBackendForState(s),
		Kinesis:                    kinesisBackendForState(s),
		Panther:                    pantherBackendForState(s),
		Splunk:                     splunkBackendForState(s),
	}
}

func httpBackendWithState(data store) *goztl.StoreHTTP {
	var b *goztl.StoreHTTP
	if !data.HTTP.IsNull() {
		bMap := data.HTTP.Attributes()
		b = &goztl.StoreHTTP{
			EndpointURL:    bMap["endpoint_url"].(types.String).ValueString(),
			VerifyTLS:      bMap["verify_tls"].(types.Bool).ValueBool(),
			Username:       optionalStringWithState(bMap["username"].(types.String)),
			Password:       optionalStringWithState(bMap["password"].(types.String)),
			Headers:        headersWithState(bMap["headers"].(types.Set)),
			Concurrency:    int(bMap["concurrency"].(types.Int64).ValueInt64()),
			RequestTimeout: int(bMap["request_timeout"].(types.Int64).ValueInt64()),
			MaxRetries:     int(bMap["max_retries"].(types.Int64).ValueInt64()),
		}
	}
	return b
}

func kinesisBackendWithState(data store) *goztl.StoreKinesis {
	var b *goztl.StoreKinesis
	if !data.Kinesis.IsNull() {
		bMap := data.Kinesis.Attributes()
		b = &goztl.StoreKinesis{
			RegionName:          bMap["region_name"].(types.String).ValueString(),
			AWSAccessKeyID:      optionalStringWithState(bMap["aws_access_key_id"].(types.String)),
			AWSSecretAccessKey:  optionalStringWithState(bMap["aws_secret_access_key"].(types.String)),
			AssumeRoleARN:       optionalStringWithState(bMap["assume_role_arn"].(types.String)),
			Stream:              bMap["stream"].(types.String).ValueString(),
			BatchSize:           int(bMap["batch_size"].(types.Int64).ValueInt64()),
			SerializationFormat: bMap["serialization_format"].(types.String).ValueString(),
		}
	}
	return b
}

func pantherBackendWithState(data store) *goztl.StorePanther {
	var b *goztl.StorePanther
	if !data.Panther.IsNull() {
		bMap := data.Panther.Attributes()
		b = &goztl.StorePanther{
			EndpointURL: bMap["endpoint_url"].(types.String).ValueString(),
			BearerToken: bMap["bearer_token"].(types.String).ValueString(),
			BatchSize:   int(bMap["batch_size"].(types.Int64).ValueInt64()),
		}
	}
	return b
}

func splunkBackendWithState(data store) *goztl.StoreSplunk {
	var b *goztl.StoreSplunk
	if !data.Splunk.IsNull() {
		bMap := data.Splunk.Attributes()
		b = &goztl.StoreSplunk{
			// HEC
			HECURL:                    bMap["hec_url"].(types.String).ValueString(),
			HECToken:                  bMap["hec_token"].(types.String).ValueString(),
			HECExtraHeaders:           headersWithState(bMap["hec_extra_headers"].(types.Set)),
			HECRequestTimeout:         int(bMap["hec_request_timeout"].(types.Int64).ValueInt64()),
			HECIndex:                  optionalStringWithState(bMap["hec_index"].(types.String)),
			HECSource:                 optionalStringWithState(bMap["hec_source"].(types.String)),
			ComputerNameAsHostSources: stringListWithStateList(bMap["computer_name_as_host_sources"].(types.List)),
			CustomHostField:           optionalStringWithState(bMap["custom_host_field"].(types.String)),
			SerialNumberField:         bMap["serial_number_field"].(types.String).ValueString(),
			BatchSize:                 int(bMap["batch_size"].(types.Int64).ValueInt64()),
			// Events URLs
			SearchAppURL: optionalStringWithState(bMap["search_app_url"].(types.String)),
			// Events search
			SearchURL:            optionalStringWithState(bMap["search_url"].(types.String)),
			SearchToken:          optionalStringWithState(bMap["search_token"].(types.String)),
			SearchExtraHeaders:   headersWithState(bMap["search_extra_headers"].(types.Set)),
			SearchRequestTimeout: int(bMap["search_request_timeout"].(types.Int64).ValueInt64()),
			SearchIndex:          optionalStringWithState(bMap["search_index"].(types.String)),
			SearchSource:         optionalStringWithState(bMap["search_source"].(types.String)),
			// Common
			VerifyTLS: bMap["verify_tls"].(types.Bool).ValueBool(),
		}
	}
	return b
}

func storeRequestWithState(data store) *goztl.StoreRequest {
	return &goztl.StoreRequest{
		Name:                       data.Name.ValueString(),
		Description:                data.Description.ValueString(),
		AdminConsole:               data.AdminConsole.ValueBool(),
		EventsURLAuthorizedRoleIDs: intListWithState(data.EventsURLAuthorizedRoleIDs),
		EventFilters:               eventFilterSetWithState(data.EventFilters),
		Backend:                    data.Backend.ValueString(),
		HTTP:                       httpBackendWithState(data),
		Kinesis:                    kinesisBackendWithState(data),
		Panther:                    pantherBackendWithState(data),
		Splunk:                     splunkBackendWithState(data),
	}
}
