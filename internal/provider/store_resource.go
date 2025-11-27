package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &StoreResource{}
var _ resource.ResourceWithImportState = &StoreResource{}

func NewStoreResource() resource.Resource {
	return &StoreResource{}
}

// StoreResource defines the resource implementation.
type StoreResource struct {
	client *goztl.Client
}

func (r *StoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_store"
}

func makeHTTPHeadersSchema(desc string) schema.SetNestedAttribute {
	desc = fmt.Sprintf("A set of additional HTTP headers to add to the %s requests.", desc)
	return schema.SetNestedAttribute{
		Description:         desc,
		MarkdownDescription: desc,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Description:         "Name of the HTTP header.",
					MarkdownDescription: "Name of the HTTP header.",
					Required:            true,
				},
				"value": schema.StringAttribute{
					Description:         "Value of the HTTP header.",
					MarkdownDescription: "Value of the HTTP header.",
					Required:            true,
					Sensitive:           true,
				},
			},
		},
		Optional: true,
		Computed: true,
		Default:  setdefault.StaticValue(types.SetValueMust(types.ObjectType{AttrTypes: headerAttrTypes}, []attr.Value{})),
	}
}

var httpBackendSchema schema.SingleNestedAttribute = schema.SingleNestedAttribute{
	Description:         "HTTP backend parameters.",
	MarkdownDescription: "HTTP backend parameters.",
	Attributes: map[string]schema.Attribute{
		"endpoint_url": schema.StringAttribute{
			Description:         "HTTP endpoint URL.",
			MarkdownDescription: "HTTP endpoint URL.",
			Required:            true,
		},
		"username": schema.StringAttribute{
			Description:         "Username for basic authentication.",
			MarkdownDescription: "Username for basic authentication.",
			Optional:            true,
		},
		"password": schema.StringAttribute{
			Description:         "Password for basic authentication.",
			MarkdownDescription: "Password for basic authentication.",
			Sensitive:           true,
			Optional:            true,
		},
		"headers": makeHTTPHeadersSchema("POST"),
		"concurrency": schema.Int64Attribute{
			Description: fmt.Sprintf(
				"Number of threads used to post to the endpoints. Defaults to %d.",
				tfStoreHTTPBackendDefaultConcurrency,
			),
			MarkdownDescription: fmt.Sprintf(
				"Number of threads used to post to the endpoints. Defaults to `%d`.",
				tfStoreHTTPBackendDefaultConcurrency,
			),
			Optional: true,
			Computed: true,
			Default:  int64default.StaticInt64(tfStoreHTTPBackendDefaultConcurrency),
			Validators: []validator.Int64{
				int64validator.AtLeast(tfStoreHTTPBackendMinConcurrency),
			},
		},
		"request_timeout": schema.Int64Attribute{
			Description: fmt.Sprintf(
				"Request timeout in seconds. Defaults to %d seconds.",
				tfStoreHTTPBackendDefaultRequestTimeout,
			),
			MarkdownDescription: fmt.Sprintf(
				"Request timeout in seconds. Defaults to `%d` seconds.",
				tfStoreHTTPBackendDefaultRequestTimeout,
			),
			Optional: true,
			Computed: true,
			Default:  int64default.StaticInt64(tfStoreHTTPBackendDefaultRequestTimeout),
			Validators: []validator.Int64{
				int64validator.Between(tfStoreHTTPBackendMinRequestTimeout, tfStoreHTTPBackendMaxRequestTimeout),
			},
		},
		"max_retries": schema.Int64Attribute{
			Description: fmt.Sprintf(
				"Number of retries after a failed request. Defaults to %d.",
				tfStoreHTTPBackendDefaultMaxRetries,
			),
			MarkdownDescription: fmt.Sprintf(
				"Number of retries after a failed request. Defaults to `%d`.",
				tfStoreHTTPBackendDefaultMaxRetries,
			),
			Optional: true,
			Computed: true,
			Default:  int64default.StaticInt64(tfStoreHTTPBackendDefaultMaxRetries),
			Validators: []validator.Int64{
				int64validator.Between(tfStoreHTTPBackendMinMaxRetries, tfStoreHTTPBackendMaxMaxRetries),
			},
		},
		"verify_tls": schema.BoolAttribute{
			Description:         "Controls whether the TLS certificates will be verified. Defaults to true.",
			MarkdownDescription: "Controls whether the TLS certificates will be verified. Defaults to `true`.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
		},
	},
	Optional: true,
}

var kinesisBackendSchema schema.SingleNestedAttribute = schema.SingleNestedAttribute{
	Description:         "Kinesis backend parameters.",
	MarkdownDescription: "Kinesis backend parameters.",
	Attributes: map[string]schema.Attribute{
		"region_name": schema.StringAttribute{
			Description:         "AWS region.",
			MarkdownDescription: "AWS region.",
			Required:            true,
		},
		"aws_access_key_id": schema.StringAttribute{
			Description:         "AWS access key ID.",
			MarkdownDescription: "AWS access key ID.",
			Optional:            true,
		},
		"aws_secret_access_key": schema.StringAttribute{
			Description:         "AWS secret access key.",
			MarkdownDescription: "AWS secret access key.",
			Optional:            true,
		},
		"assume_role_arn": schema.StringAttribute{
			Description:         "ARN of the AWS role to assume.",
			MarkdownDescription: "`ARN` of the AWS role to assume.",
			Optional:            true,
		},
		"stream": schema.StringAttribute{
			Description:         "Name of the Kinesis stream.",
			MarkdownDescription: "Name of the Kinesis stream.",
			Required:            true,
		},
		"batch_size": schema.Int64Attribute{
			Description: fmt.Sprintf(
				"Number of events sent in a single request. Defaults to %d. Must be between %d and %d.",
				tfStoreKinesisBackendDefaultBatchSize,
				tfStoreKinesisBackendMinBatchSize,
				tfStoreKinesisBackendMaxBatchSize,
			),
			MarkdownDescription: fmt.Sprintf(
				"Number of events sent in a single request. Defaults to `%d`. Must be between `%d` and `%d`.",
				tfStoreKinesisBackendDefaultBatchSize,
				tfStoreKinesisBackendMinBatchSize,
				tfStoreKinesisBackendMaxBatchSize,
			),
			Optional: true,
			Computed: true,
			Default:  int64default.StaticInt64(tfStoreKinesisBackendDefaultBatchSize),
			Validators: []validator.Int64{
				int64validator.Between(tfStoreKinesisBackendMinBatchSize, tfStoreKinesisBackendMaxBatchSize),
			},
		},
		"serialization_format": schema.StringAttribute{
			Description: fmt.Sprintf(
				"Zentral event serialization format. Either %s or %s.",
				tfStoreKinesisSerializationFormatZentral,
				tfStoreKinesisSerializationFormatFirehoseV1,
			),
			MarkdownDescription: fmt.Sprintf(
				"Zentral event serialization format. Either `%s` or `%s`.",
				tfStoreKinesisSerializationFormatZentral,
				tfStoreKinesisSerializationFormatFirehoseV1,
			),
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf([]string{tfStoreKinesisSerializationFormatZentral, tfStoreKinesisSerializationFormatFirehoseV1}...),
			},
		},
	},
	Optional: true,
}

var splunkBackendSchema schema.SingleNestedAttribute = schema.SingleNestedAttribute{
	Description:         "Splunk backend parameters.",
	MarkdownDescription: "Splunk backend parameters.",
	Attributes: map[string]schema.Attribute{
		// HEC
		"hec_url": schema.StringAttribute{
			Description:         "HEC endpoint URL.",
			MarkdownDescription: "HEC endpoint URL.",
			Required:            true,
		},
		"hec_token": schema.StringAttribute{
			Description:         "HEC token.",
			MarkdownDescription: "HEC token.",
			Required:            true,
			Sensitive:           true,
		},
		"hec_extra_headers": makeHTTPHeadersSchema("HEC"),
		"hec_request_timeout": schema.Int64Attribute{
			Description:         fmt.Sprintf("HEC request timeout in seconds. Defaults to %d seconds.", tfStoreSplunkBackendDefaultTimeout),
			MarkdownDescription: fmt.Sprintf("HEC request timeout in seconds. Defaults to `%d` seconds.", tfStoreSplunkBackendDefaultTimeout),
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(tfStoreSplunkBackendDefaultTimeout),
		},
		"hec_index": schema.StringAttribute{
			Description:         "HEC index. Usually enforced in the HEC configuration.",
			MarkdownDescription: "HEC index. Usually enforced in the HEC configuration.",
			Optional:            true,
		},
		"hec_source": schema.StringAttribute{
			Description:         "HEC source. Usually enforced in the HEC configuration.",
			MarkdownDescription: "HEC source. Usually enforced in the HEC configuration.",
			Optional:            true,
		},
		"computer_name_as_host_sources": schema.ListAttribute{
			Description:         "List of the preferred inventory sources for the events computer_name field.",
			MarkdownDescription: "List of the preferred inventory sources for the events `computer_name` field.",
			ElementType:         types.StringType,
			Optional:            true,
			Computed:            true,
			Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
		},
		"custom_host_field": schema.StringAttribute{
			Description:         "Name of an optional field to copy the host field value to.",
			MarkdownDescription: "Name of an optional field to copy the `host` field value to.",
			Optional:            true,
		},
		"serial_number_field": schema.StringAttribute{
			Description:         "Name of the field to use for the events machine serial number. Defaults to machine_serial_number.",
			MarkdownDescription: "Name of the field to use for the events machine serial number. Defaults to `machine_serial_number`.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("machine_serial_number"),
		},
		"batch_size": schema.Int64Attribute{
			Description: fmt.Sprintf(
				"Number of events sent in a single request. Defaults to %d. Must be between %d and %d.",
				tfStoreSplunkBackendDefaultBatchSize,
				tfStoreSplunkBackendMinBatchSize,
				tfStoreSplunkBackendMaxBatchSize,
			),
			MarkdownDescription: fmt.Sprintf(
				"Number of events sent in a single request. Defaults to `%d`. Must be between `%d` and `%d`.",
				tfStoreSplunkBackendDefaultBatchSize,
				tfStoreSplunkBackendMinBatchSize,
				tfStoreSplunkBackendMaxBatchSize,
			),
			Optional: true,
			Computed: true,
			Default:  int64default.StaticInt64(tfStoreSplunkBackendDefaultBatchSize),
			Validators: []validator.Int64{
				int64validator.Between(tfStoreSplunkBackendMinBatchSize, tfStoreSplunkBackendMaxBatchSize),
			},
		},
		// Events URLs
		"search_app_url": schema.StringAttribute{
			Description:         "Base URL of the Splunk search application. Used to build the links to the Splunk instance displayed when browsing the events in the Zentral admin console.",
			MarkdownDescription: "Base URL of the Splunk search application. Used to build the links to the Splunk instance displayed when browsing the events in the Zentral admin console.",
			Optional:            true,
		},
		// Events Search
		"search_url": schema.StringAttribute{
			Description:         "Splunk API base URL. Used in combination with search_token to fetch the events displayed in the Zentral admin console.",
			MarkdownDescription: "Splunk API base URL. Used in combination with `search_token` to fetch the events displayed in the Zentral admin console.",
			Optional:            true,
		},
		"search_token": schema.StringAttribute{
			Description:         "Splunk API token. Used in combination with search_url.",
			MarkdownDescription: "Splunk API token. Used in combination with `search_url`.",
			Optional:            true,
			Sensitive:           true,
		},
		"search_extra_headers": makeHTTPHeadersSchema("search API"),
		"search_request_timeout": schema.Int64Attribute{
			Description:         fmt.Sprintf("Search API request timeout in seconds. Defaults to %d seconds.", tfStoreSplunkBackendDefaultTimeout),
			MarkdownDescription: fmt.Sprintf("Search API request timeout in seconds. Defaults to `%d` seconds.", tfStoreSplunkBackendDefaultTimeout),
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(tfStoreSplunkBackendDefaultTimeout),
		},
		"search_index": schema.StringAttribute{
			Description:         "Index to use with the Search API.",
			MarkdownDescription: "Index to use with the Search API.",
			Optional:            true,
		},
		"search_source": schema.StringAttribute{
			Description:         "Source to use with the Search API.",
			MarkdownDescription: "Source to use with the Search API.",
			Optional:            true,
		},
		// Common
		"verify_tls": schema.BoolAttribute{
			Description:         "Controls whether the TLS certificates will be verified. Defaults to `true`.",
			MarkdownDescription: "Controls whether the TLS certificates will be verified. Defaults to `true`.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
		},
	},
	Optional: true,
}

func makeEventFiltersSchema(desc string) schema.SetNestedAttribute {
	lowerDesc := strings.ToLower(desc)
	return schema.SetNestedAttribute{
		Description:         fmt.Sprintf("%sd event filters.", desc),
		MarkdownDescription: fmt.Sprintf("%sd event filters.", desc),
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"tags": schema.SetAttribute{
					Description:         fmt.Sprintf("Set of the tags of the events to %s.", lowerDesc),
					MarkdownDescription: fmt.Sprintf("Set of the tags of the events to %s.", lowerDesc),
					ElementType:         types.StringType,
					Optional:            true,
				},
				"event_type": schema.SetAttribute{
					Description:         fmt.Sprintf("Set of the event types of the events to %s.", lowerDesc),
					MarkdownDescription: fmt.Sprintf("Set of the event types of the events to %s.", lowerDesc),
					ElementType:         types.StringType,
					Optional:            true,
				},
				"routing_key": schema.SetAttribute{
					Description:         fmt.Sprintf("Set of the routing keys of the events to %s.", lowerDesc),
					MarkdownDescription: fmt.Sprintf("Set of the routing keys of the events to %s.", lowerDesc),
					ElementType:         types.StringType,
					Optional:            true,
				},
			},
		},
		Optional: true,
	}
}

func (r *StoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages stores.",
		MarkdownDescription: "The resource `zentral_store` manages stores.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the store (UUID).",
				MarkdownDescription: "`ID` of the store (UUID).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the store.",
				MarkdownDescription: "Name of the store.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the store.",
				MarkdownDescription: "Description of the store.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"admin_console": schema.BoolAttribute{
				Description:         "Controls wether the store is used in the Zentral admin console. Only one store can be used in the admin console. Defaults to false",
				MarkdownDescription: "Controls wether the store is used in the Zentral admin console. Only one store can be used in the admin console. Defaults to `false`",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"events_url_authorized_role_ids": schema.SetAttribute{
				Description:         "The IDs of the roles authorized to see the links to the external event store.",
				MarkdownDescription: "The `ID`s of the roles authorized to see the links to the external event store.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.Int64Type, []attr.Value{})),
			},
			"event_filters": schema.SingleNestedAttribute{
				Description:         "Used to filter the events sent to the store. By default, all the events are sent to the store.",
				MarkdownDescription: "Used to filter the events sent to the store. By default, all the events are sent to the store.",
				Attributes: map[string]schema.Attribute{
					"included_event_filters": makeEventFiltersSchema("Include"),
					"excluded_event_filters": makeEventFiltersSchema("Exclude"),
				},
				Optional: true,
				Computed: true,
				Default:  objectdefault.StaticValue(defaultEventFilterSet()),
			},
			"backend": schema.StringAttribute{
				Description:         "Store backend.",
				MarkdownDescription: "Store backend.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{tfStoreHTTPBackend, tfStoreKinesisBackend, tfStoreSplunkBackend}...),
				},
			},
			"http":    httpBackendSchema,
			"kinesis": kinesisBackendSchema,
			"splunk":  splunkBackendSchema,
		},
	}
}

func (r *StoreResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*goztl.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *goztl.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *StoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data store

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlS, _, err := r.client.Stores.Create(ctx, storeRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create store, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a store")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, storeForState(ztlS))...)
}

func (r *StoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data store

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlS, _, err := r.client.Stores.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read store %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a store")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, storeForState(ztlS))...)
}

func (r *StoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data store

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlS, _, err := r.client.Stores.Update(ctx, data.ID.ValueString(), storeRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update store %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a store")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, storeForState(ztlS))...)
}

func (r *StoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data store

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Stores.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete store %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a store")
}

func (r *StoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "store", req, resp)
}
