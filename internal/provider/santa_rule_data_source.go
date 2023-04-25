package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &SantaRuleDataSource{}

func NewSantaRuleDataSource() datasource.DataSource {
	return &SantaRuleDataSource{}
}

// SantaRuleDataSource defines the data source implementation.
type SantaRuleDataSource struct {
	client *goztl.Client
}

func (d *SantaRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_santa_rule"
}

func (d *SantaRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Santa rule to be retrieved by its ID.",
		MarkdownDescription: "The data source `zentral_santa_rule` allows details of a Santa rule to be retrieved by its `ID`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Santa rule.",
				MarkdownDescription: "`ID` of the Santa rule.",
				Optional:            true,
			},
			"configuration_id": schema.Int64Attribute{
				Description:         "ID of the Santa configuration.",
				MarkdownDescription: "`ID` of the Santa configuration.",
				Computed:            true,
			},
			"policy": schema.StringAttribute{
				Description:         "Policy. Valid values are ALLOWLIST, BLOCKLIST, SILENT_BLOCKLIST and ALLOWLIST_COMPILER.",
				MarkdownDescription: "Policy. Valid values are `ALLOWLIST`, `BLOCKLIST`, `SILENT_BLOCKLIST` and `ALLOWLIST_COMPILER`.",
				Computed:            true,
			},
			"target_type": schema.StringAttribute{
				Description:         "Target type. Valid values are BINARY, BUNDLE, CERTIFICATE and TEAMID.",
				MarkdownDescription: "Target type. Valid values are `BINARY`, `BUNDLE`, `CERTIFICATE` and `TEAMID`.",
				Computed:            true,
			},
			"target_identifier": schema.StringAttribute{
				Description:         "Target identifier: binary, bundle, certificate sha256 or team ID.",
				MarkdownDescription: "Target identifier: binary, bundle, certificate sha256 or team ID.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the rule. Only displayed in the Zentral GUI.",
				MarkdownDescription: "Description of the rule. Only displayed in the Zentral GUI.",
				Computed:            true,
			},
			"custom_message": schema.StringAttribute{
				Description:         "Custom message displayed in the popover when a binary is blocked.",
				MarkdownDescription: "Custom message displayed in the popover when a binary is blocked.",
				Computed:            true,
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
				Computed:            true,
			},
			"excluded_primary_users": schema.SetAttribute{
				Description:         "The excluded primary users used to scope the rule.",
				MarkdownDescription: "The excluded primary users used to scope the rule.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"serial_numbers": schema.SetAttribute{
				Description:         "The serial numbers used to scope the rule.",
				MarkdownDescription: "The serial numbers used to scope the rule.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"excluded_serial_numbers": schema.SetAttribute{
				Description:         "The excluded serial numbers used to scope the rule.",
				MarkdownDescription: "The excluded serial numbers used to scope the rule.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags used to scope the rule.",
				MarkdownDescription: "The `ID`s of the tags used to scope the rule.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"excluded_tag_ids": schema.SetAttribute{
				Description:         "The IDs of the excluded tags used to scope the rule.",
				MarkdownDescription: "The `ID`s of the excluded tags used to scope the rule.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				Description:         "Rule version.",
				MarkdownDescription: "Rule version.",
				Computed:            true,
			},
		},
	}
}

func (d *SantaRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*goztl.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *goztl.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *SantaRuleDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data santaRule
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid `zentral_santa_rule` data source",
			"`id` missing",
		)
	}
}

func (d *SantaRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data santaRule

	// Read Terraform rule data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlSR *goztl.SantaRule
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlSR, _, err = d.client.SantaRules.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Santa rule '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	}

	if ztlSR != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, santaRuleForState(ztlSR))...)
	}
}
