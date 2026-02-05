package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
var _ resource.Resource = &ProbeActionResource{}
var _ resource.ResourceWithImportState = &ProbeActionResource{}

func NewProbeActionResource() resource.Resource {
	return &ProbeActionResource{}
}

// ProbeActionResource defines the resource implementation.
type ProbeActionResource struct {
	client *goztl.Client
}

func (r *ProbeActionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_probe_action"
}

func (r *ProbeActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages probe actions.",
		MarkdownDescription: "The resource `zentral_probe_action` manages probe actions.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the action.",
				MarkdownDescription: "`ID` of the action.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the action.",
				MarkdownDescription: "Name of the action.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the action.",
				MarkdownDescription: "Description of the action.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"backend": schema.StringAttribute{
				Description:         "Action backend.",
				MarkdownDescription: "Action backend.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{tfProbeActionHTTPPostBackend, tfProbeActionSlackIncomingWebhookBackend}...),
				},
			},
			"http_post": schema.SingleNestedAttribute{
				Description:         "HTTP Post backend parameters.",
				MarkdownDescription: "HTTP Post backend parameters.",
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Description:         "URL.",
						MarkdownDescription: "`URL`.",
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
					"headers": schema.SetNestedAttribute{
						Description:         "A set of additional HTTP headers to add to the requests.",
						MarkdownDescription: "A set of additional HTTP headers to add to the requests.",
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
						Default:  setdefault.StaticValue(types.SetValueMust(types.ObjectType{AttrTypes: probeActionHTTPPostHeaderAttrTypes}, []attr.Value{})),
					},
					"cel_transformation": schema.StringAttribute{
						Description:         "CEL expression that is used to transform the event data. The input to the expression is a Map with two keys: metadata for the event metadata and payload for the event payload.",
						MarkdownDescription: "CEL expression that is used to transform the event data. The input to the expression is a `Map` with two keys: `metadata` for the event metadata and `payload` for the event payload.",
						Optional:            true,
					},
				},
				Optional: true,
			},
			"slack_incoming_webhook": schema.SingleNestedAttribute{
				Description:         "Slack incoming webhook backend parameters.",
				MarkdownDescription: "Slack incoming webhook backend parameters.",
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Description:         "URL.",
						MarkdownDescription: "`URL`.",
						Required:            true,
						Sensitive:           true,
					},
				},
				Optional: true,
			},
		},
	}
}

func (r *ProbeActionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProbeActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data probeAction

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlPA, _, err := r.client.ProbesActions.Create(ctx, probeActionRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create probe action, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a probe action")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, probeActionForState(ztlPA))...)
}

func (r *ProbeActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data probeAction

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlPA, _, err := r.client.ProbesActions.GetByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read probe action %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a probe action")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, probeActionForState(ztlPA))...)
}

func (r *ProbeActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data probeAction

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlPA, _, err := r.client.ProbesActions.Update(ctx, data.ID.ValueString(), probeActionRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update probe action %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a probe action")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, probeActionForState(ztlPA))...)
}

func (r *ProbeActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data probeAction

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.ProbesActions.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete probe action %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a probe action")
}

func (r *ProbeActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralUUID(ctx, "probe action", req, resp)
}
