package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const (
	ztlSantaAllowlist         int    = 1
	ztlSantaBlocklist                = 2
	ztlSantaSilentBlocklist          = 3
	ztlSantaAllowlistCompiler        = 5
	ztlSantaCEL                      = 9
	tfSantaAllowlist          string = "ALLOWLIST"
	tfSantaBlocklist                 = "BLOCKLIST"
	tfSantaSilentBlocklist           = "SILENT_BLOCKLIST"
	tfSantaAllowlistCompiler         = "ALLOWLIST_COMPILER"
	tfSantaCEL                       = "CEL"
)

type santaRule struct {
	ID                    types.Int64  `tfsdk:"id"`
	ConfigurationID       types.Int64  `tfsdk:"configuration_id"`
	Policy                types.String `tfsdk:"policy"`
	CELExpr               types.String `tfsdk:"cel_expr"`
	TargetType            types.String `tfsdk:"target_type"`
	TargetIdentifier      types.String `tfsdk:"target_identifier"`
	Description           types.String `tfsdk:"description"`
	CustomMessage         types.String `tfsdk:"custom_message"`
	CustomURL             types.String `tfsdk:"custom_url"`
	RulesetID             types.Int64  `tfsdk:"ruleset_id"`
	PrimaryUsers          types.Set    `tfsdk:"primary_users"`
	ExcludedPrimaryUsers  types.Set    `tfsdk:"excluded_primary_users"`
	SerialNumbers         types.Set    `tfsdk:"serial_numbers"`
	ExcludedSerialNumbers types.Set    `tfsdk:"excluded_serial_numbers"`
	TagIDs                types.Set    `tfsdk:"tag_ids"`
	ExcludedTagIDs        types.Set    `tfsdk:"excluded_tag_ids"`
	Version               types.Int64  `tfsdk:"version"`
}

func santaRuleForState(sr *goztl.SantaRule) santaRule {
	var policy string
	switch sr.Policy {
	case ztlSantaAllowlist:
		policy = tfSantaAllowlist
	case ztlSantaBlocklist:
		policy = tfSantaBlocklist
	case ztlSantaSilentBlocklist:
		policy = tfSantaSilentBlocklist
	case ztlSantaAllowlistCompiler:
		policy = tfSantaAllowlistCompiler
	case ztlSantaCEL:
		policy = tfSantaCEL
	default:
		panic("Unknown Santa rule policy")
	}

	return santaRule{
		ID:                    types.Int64Value(int64(sr.ID)),
		ConfigurationID:       types.Int64Value(int64(sr.ConfigurationID)),
		Policy:                types.StringValue(policy),
		CELExpr:               types.StringValue(sr.CELExpr),
		TargetType:            types.StringValue(sr.TargetType),
		TargetIdentifier:      types.StringValue(sr.TargetIdentifier),
		Description:           types.StringValue(sr.Description),
		CustomMessage:         types.StringValue(sr.CustomMessage),
		CustomURL:             types.StringValue(sr.CustomURL),
		RulesetID:             optionalInt64ForState(sr.RulesetID),
		PrimaryUsers:          stringSetForState(sr.PrimaryUsers),
		ExcludedPrimaryUsers:  stringSetForState(sr.ExcludedPrimaryUsers),
		SerialNumbers:         stringSetForState(sr.SerialNumbers),
		ExcludedSerialNumbers: stringSetForState(sr.ExcludedSerialNumbers),
		TagIDs:                int64SetForState(sr.TagIDs),
		ExcludedTagIDs:        int64SetForState(sr.ExcludedTagIDs),
		Version:               types.Int64Value(int64(sr.Version)),
	}
}

func santaRuleRequestWithState(data santaRule) *goztl.SantaRuleRequest {
	var policy int
	switch data.Policy {
	case types.StringValue(tfSantaAllowlist):
		policy = ztlSantaAllowlist
	case types.StringValue(tfSantaBlocklist):
		policy = ztlSantaBlocklist
	case types.StringValue(tfSantaSilentBlocklist):
		policy = ztlSantaSilentBlocklist
	case types.StringValue(tfSantaAllowlistCompiler):
		policy = ztlSantaAllowlistCompiler
	case types.StringValue(tfSantaCEL):
		policy = ztlSantaCEL
	default:
		panic("Unknown Santa rule policy")
	}

	return &goztl.SantaRuleRequest{
		ConfigurationID:       int(data.ConfigurationID.ValueInt64()),
		Policy:                policy,
		CELExpr:               data.CELExpr.ValueString(),
		TargetType:            data.TargetType.ValueString(),
		TargetIdentifier:      data.TargetIdentifier.ValueString(),
		Description:           data.Description.ValueString(),
		CustomMessage:         data.CustomMessage.ValueString(),
		CustomURL:             data.CustomURL.ValueString(),
		PrimaryUsers:          stringListWithStateSet(data.PrimaryUsers),
		ExcludedPrimaryUsers:  stringListWithStateSet(data.ExcludedPrimaryUsers),
		SerialNumbers:         stringListWithStateSet(data.SerialNumbers),
		ExcludedSerialNumbers: stringListWithStateSet(data.ExcludedSerialNumbers),
		TagIDs:                intListWithState(data.TagIDs),
		ExcludedTagIDs:        intListWithState(data.ExcludedTagIDs),
	}
}
