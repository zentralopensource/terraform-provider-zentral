package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &SantaScopedClientModeResource{}
var _ resource.ResourceWithImportState = &SantaScopedClientModeResource{}
var _ resource.ResourceWithValidateConfig = &SantaScopedClientModeResource{}

func NewSantaScopedClientModeResource() resource.Resource {
	return &SantaScopedClientModeResource{}
}

// SantaScopedClientModeResource defines the resource implementation.
type SantaScopedClientModeResource struct {
	client *goztl.Client
}

func (r *SantaScopedClientModeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_santa_scoped_client_mode"
}

func (r *SantaScopedClientModeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := santaScopedConfigurationItemAttributes("client mode")
	attributes["client_mode"] = schema.StringAttribute{
		Description:         "Client mode of the machines in scope. Valid values are MONITOR and LOCKDOWN. A configuration has at most one entry per client mode.",
		MarkdownDescription: "Client mode of the machines in scope. Valid values are `MONITOR` and `LOCKDOWN`. A configuration has at most one entry per client mode.",
		Required:            true,
		Validators: []validator.String{
			stringvalidator.OneOf(tfSantaMonitor, tfSantaLockdown),
		},
	}
	attributes["event_detail_source"] = schema.StringAttribute{
		Description: "Where the block notification button of the machines in scope comes from. Valid values are INHERIT, VOTING_PORTAL, CUSTOM and NONE. " +
			"INHERIT takes the button of the configuration, VOTING_PORTAL sends a link to the user portal of the voting realm of the configuration, " +
			"CUSTOM sends event_detail_url and event_detail_text, and NONE removes the button. Defaults to INHERIT.",
		MarkdownDescription: "Where the block notification button of the machines in scope comes from. Valid values are `INHERIT`, `VOTING_PORTAL`, `CUSTOM` and `NONE`. " +
			"`INHERIT` takes the button of the configuration, `VOTING_PORTAL` sends a link to the user portal of the voting realm of the configuration, " +
			"`CUSTOM` sends `event_detail_url` and `event_detail_text`, and `NONE` removes the button. Defaults to `INHERIT`.",
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString(tfSantaEventDetailInherit),
		Validators: []validator.String{
			stringvalidator.OneOf(
				tfSantaEventDetailInherit, tfSantaEventDetailVotingPortal, tfSantaEventDetailCustom, tfSantaEventDetailNone,
			),
		},
	}
	attributes["event_detail_url"] = schema.StringAttribute{
		Description: "URL of the block notification button. Required when event_detail_source is CUSTOM, and only allowed then. " +
			"The following sequences are replaced: %file_identifier%, %bundle_or_file_identifier%, %file_bundle_id%, %team_id%, " +
			"%signing_id%, %cdhash%, %username%, %machine_id%, %hostname%, %uuid%, %serial%.",
		MarkdownDescription: "URL of the block notification button. Required when `event_detail_source` is `CUSTOM`, and only allowed then. " +
			"The following sequences are replaced: `%file_identifier%`, `%bundle_or_file_identifier%`, `%file_bundle_id%`, `%team_id%`, " +
			"`%signing_id%`, `%cdhash%`, `%username%`, `%machine_id%`, `%hostname%`, `%uuid%`, `%serial%`.",
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString(""),
		Validators: []validator.String{
			stringvalidator.LengthAtMost(1024),
			trimmedStringValidator,
		},
	}
	attributes["event_detail_text"] = schema.StringAttribute{
		Description:         "Label of the block notification button. Only allowed when event_detail_source is CUSTOM or VOTING_PORTAL.",
		MarkdownDescription: "Label of the block notification button. Only allowed when `event_detail_source` is `CUSTOM` or `VOTING_PORTAL`.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
		Validators: []validator.String{
			trimmedStringValidator,
		},
	}

	resp.Schema = schema.Schema{
		Description: "Manages Santa scoped client modes.",
		MarkdownDescription: "The resource `zentral_santa_scoped_client_mode` manages Santa scoped client modes: a client mode, and a block notification button, " +
			"for some of the machines of a Santa configuration. The first scope level that matches a machine decides: serial numbers, then primary users, then tags. " +
			"An entry without any scope matches every machine its exclusions do not.",
		Attributes: attributes,
	}
}

func (r *SantaScopedClientModeResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data santaScopedClientMode

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	validateSantaEventDetail(&resp.Diagnostics, data.EventDetailSource, data.EventDetailURL, data.EventDetailText, tfSantaEventDetailInherit)
}

func (r *SantaScopedClientModeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SantaScopedClientModeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data santaScopedClientMode

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlSSCM, _, err := r.client.SantaScopedClientModes.Create(ctx, santaScopedClientModeRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Santa scoped client mode, got error: %s", err),
		)
		addEndpointSupportDiagnostic(&resp.Diagnostics, err, "Santa scoped client mode", minZentralVersionSantaScopedConfigurationItems)
		return
	}

	tflog.Trace(ctx, "created a Santa scoped client mode")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, santaScopedClientModeForState(ztlSSCM))...)
}

func (r *SantaScopedClientModeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data santaScopedClientMode

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlSSCM, _, err := r.client.SantaScopedClientModes.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Santa scoped client mode %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Santa scoped client mode")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, santaScopedClientModeForState(ztlSSCM))...)
}

func (r *SantaScopedClientModeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data santaScopedClientMode

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlSSCM, _, err := r.client.SantaScopedClientModes.Update(ctx, int(data.ID.ValueInt64()), santaScopedClientModeRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Santa scoped client mode %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Santa scoped client mode")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, santaScopedClientModeForState(ztlSSCM))...)
}

func (r *SantaScopedClientModeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data santaScopedClientMode

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.SantaScopedClientModes.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Santa scoped client mode %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Santa scoped client mode")
}

func (r *SantaScopedClientModeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Santa scoped client mode", req, resp)
}
