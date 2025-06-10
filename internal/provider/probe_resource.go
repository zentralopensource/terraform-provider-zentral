package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ProbeResource{}
var _ resource.ResourceWithImportState = &ProbeResource{}

func NewProbeResource() resource.Resource {
	return &ProbeResource{}
}

// ProbeResource defines the resource implementation.
type ProbeResource struct {
	client *goztl.Client
}

func (r *ProbeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_probe"
}

func (r *ProbeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	nestedInventoryFilter := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"meta_business_unit_ids": schema.SetAttribute{
				Description:         "The IDs of the meta business units.",
				MarkdownDescription: "The `ID`s of the meta business units.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.Int64Type, []attr.Value{})),
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the machine tags.",
				MarkdownDescription: "The `ID`s of the machine tags.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.Int64Type, []attr.Value{})),
			},
			"platforms": schema.SetAttribute{
				Description:         "The machine platforms (LINUX, MACOS, WINDOWS, ANDROID, IOS, IPADOS, TVOS).",
				MarkdownDescription: "The machine platforms (`LINUX`, `MACOS`, `WINDOWS`, `ANDROID`, `IOS`, `IPADOS`, `TVOS`).",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"LINUX", "MACOS", "WINDOWS", "ANDROID", "IOS", "IPADOS", "TVOS"}...)),
				},
			},
			"types": schema.SetAttribute{
				Description:         "The machine types (DESKTOP, EC2, LAPTOP, MOBILE, SERVER, TABLET, TV, VM).",
				MarkdownDescription: "The machine platforms (`DESKTOP`, `EC2`, `LAPTOP`, `MOBILE`, `SERVER`, `TABLET`, `TV`, `VM`).",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"DESKTOP", "EC2", "LAPTOP", "MOBILE", "SERVER", "TABLET", "TV", "VM"}...)),
				},
			},
		},
	}

	nestedMetadataFilter := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"event_types": schema.SetAttribute{
				Description:         "The event types.",
				MarkdownDescription: "The event types.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"event_tags": schema.SetAttribute{
				Description:         "The event tags.",
				MarkdownDescription: "The event tags.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"event_routing_keys": schema.SetAttribute{
				Description:         "The event routing keys.",
				MarkdownDescription: "The event routing keys.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
		},
	}

	resp.Schema = schema.Schema{
		Description:         "Manages probes.",
		MarkdownDescription: "The resource `zentral_probe` manages probes.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the probe.",
				MarkdownDescription: "`ID` of the probe.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the probe.",
				MarkdownDescription: "Name of the probe.",
				Required:            true,
			},
			"slug": schema.StringAttribute{
				Description:         "Slug of the probe.",
				MarkdownDescription: "Slug of the probe.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the probe.",
				MarkdownDescription: "Description of the probe.",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
			},
			"inventory_filters": schema.SetNestedAttribute{
				Description:         "A set of filters used to filter the events based on inventory attributes.",
				MarkdownDescription: "A set of filters used to filter the events based on inventory attributes.",
				NestedObject:        nestedInventoryFilter,
				Computed:            true,
				Optional:            true,
			},
			"metadata_filters": schema.SetNestedAttribute{
				Description:         "A set of filters used to filter the events based on their metadata.",
				MarkdownDescription: "A set of filters used to filter the events based on their metadata.",
				NestedObject:        nestedMetadataFilter,
				Computed:            true,
				Optional:            true,
			},
			"payload_filters": schema.SetAttribute{
				Description:         "A set of filters used to filter the events based on their payload.",
				MarkdownDescription: "A set of filters used to filter the events based on their payload.",
				ElementType: types.SetType{
					ElemType: types.ObjectType{AttrTypes: probePayloadFilterItemAttrTypes},
				},
				Computed: true,
				Optional: true,
			},
			"active": schema.BoolAttribute{
				Description:         "If true, the probe is active.",
				MarkdownDescription: "If true, the probe is active.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"action_ids": schema.SetAttribute{
				Description:         "The IDs of the actions triggered by the probe when an event matches.",
				MarkdownDescription: "The `ID`s of the actions triggered by the probe when an event matches.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"incident_severity": schema.Int64Attribute{
				Description:         "Severity of the incidents triggered by this probe (100, 200, 300).",
				MarkdownDescription: "Severity of the incidents triggered by this probe (`100`, `200`, `300`).",
				Required:            false,
				Validators: []validator.Int64{
					int64validator.OneOf([]int64{100, 200, 300}...),
				},
				Optional: true,
			},
		},
	}
}

func (r *ProbeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProbeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data probe

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlP, _, err := r.client.Probes.Create(ctx, probeRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create probe action, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a probe action")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, probeForState(ztlP))...)
}

func (r *ProbeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data probe

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlP, _, err := r.client.Probes.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read probe action %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a probe action")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, probeForState(ztlP))...)
}

func (r *ProbeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data probe

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlP, _, err := r.client.Probes.Update(ctx, int(data.ID.ValueInt64()), probeRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update probe action %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a probe action")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, probeForState(ztlP))...)
}

func (r *ProbeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data probe

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Probes.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete probe action %d, got error: %s", int(data.ID.ValueInt64()), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a probe action")
}

func (r *ProbeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "probe action", req, resp)
}
