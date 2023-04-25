package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
var _ resource.Resource = &SantaRuleResource{}
var _ resource.ResourceWithImportState = &SantaRuleResource{}

func NewSantaRuleResource() resource.Resource {
	return &SantaRuleResource{}
}

// SantaRuleResource defines the resource implementation.
type SantaRuleResource struct {
	client *goztl.Client
}

func (r *SantaRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_santa_rule"
}

func (r *SantaRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Santa rules.",
		MarkdownDescription: "The resource `zentral_santa_rule` manages Santa rules.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Santa rule.",
				MarkdownDescription: "`ID` of the Santa rule.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"configuration_id": schema.Int64Attribute{
				Description:         "ID of the Santa configuration.",
				MarkdownDescription: "`ID` of the Santa configuration.",
				Required:            true,
			},
			"policy": schema.StringAttribute{
				Description:         "Policy. Valid values are ALLOWLIST, BLOCKLIST, SILENT_BLOCKLIST and ALLOWLIST_COMPILER.",
				MarkdownDescription: "Policy. Valid values are `ALLOWLIST`, `BLOCKLIST`, `SILENT_BLOCKLIST` and `ALLOWLIST_COMPILER`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{tfSantaAllowlist, tfSantaBlocklist, tfSantaSilentBlocklist, tfSantaAllowlistCompiler}...),
				},
			},
			"target_type": schema.StringAttribute{
				Description:         "Target type. Valid values are BINARY, BUNDLE, CERTIFICATE and TEAMID.",
				MarkdownDescription: "Target type. Valid values are `BINARY`, `BUNDLE`, `CERTIFICATE` and `TEAMID`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"BINARY", "BUNDLE", "CERTIFICATE", "TEAMID"}...),
				},
			},
			"target_identifier": schema.StringAttribute{
				Description:         "Target identifier: binary, bundle, certificate sha256 or team ID.",
				MarkdownDescription: "Target identifier: binary, bundle, certificate sha256 or team ID.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the rule. Only displayed in the Zentral GUI.",
				MarkdownDescription: "Description of the rule. Only displayed in the Zentral GUI.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"custom_message": schema.StringAttribute{
				Description:         "Custom message displayed in the popover when a binary is blocked.",
				MarkdownDescription: "Custom message displayed in the popover when a binary is blocked.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"ruleset_id": schema.Int64Attribute{
				Description:         "ID of the Santa ruleset.",
				MarkdownDescription: "`ID` of the Santa ruleset.",
				Computed:            true,
			},
			"primary_users": schema.SetAttribute{
				Description:         "The primary users used to scope the rule.",
				MarkdownDescription: "The primary users used to scope the rule.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"excluded_primary_users": schema.SetAttribute{
				Description:         "The excluded primary users used to scope the rule.",
				MarkdownDescription: "The excluded primary users used to scope the rule.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"serial_numbers": schema.SetAttribute{
				Description:         "The serial numbers used to scope the rule.",
				MarkdownDescription: "The serial numbers used to scope the rule.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"excluded_serial_numbers": schema.SetAttribute{
				Description:         "The excluded serial numbers used to scope the rule.",
				MarkdownDescription: "The excluded serial numbers used to scope the rule.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags used to scope the rule.",
				MarkdownDescription: "The `ID`s of the tags used to scope the rule.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.Int64Type, []attr.Value{})),
			},
			"excluded_tag_ids": schema.SetAttribute{
				Description:         "The IDs of the excluded tags used to scope the rule.",
				MarkdownDescription: "The `ID`s of the excluded tags used to scope the rule.",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.Int64Type, []attr.Value{})),
			},
			"version": schema.Int64Attribute{
				Description:         "Rule version.",
				MarkdownDescription: "Rule version.",
				Computed:            true,
			},
		},
	}
}

func (r *SantaRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SantaRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data santaRule

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlSR, _, err := r.client.SantaRules.Create(ctx, santaRuleRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Santa rule, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Santa rule")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, santaRuleForState(ztlSR))...)
}

func (r *SantaRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data santaRule

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlSR, _, err := r.client.SantaRules.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Santa rule %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Santa rule")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, santaRuleForState(ztlSR))...)
}

func (r *SantaRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data santaRule

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlSR, _, err := r.client.SantaRules.Update(ctx, int(data.ID.ValueInt64()), santaRuleRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Santa rule %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Santa rule")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, santaRuleForState(ztlSR))...)
}

func (r *SantaRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data santaRule

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.SantaRules.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Santa rule %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Santa rule")
}

func (r *SantaRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Santa rule", req, resp)
}
