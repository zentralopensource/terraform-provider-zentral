package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const (
	ztlSantaAllowlist         int    = 1
	ztlSantaBlocklist                = 2
	ztlSantaSilentBlocklist          = 3
	ztlSantaAllowlistCompiler        = 5
	tfSantaAllowlist          string = "ALLOWLIST"
	tfSantaBlocklist                 = "BLOCKLIST"
	tfSantaSilentBlocklist           = "SILENT_BLOCKLIST"
	tfSantaAllowlistCompiler         = "ALLOWLIST_COMPILER"
)

type santaRule struct {
	ID                    types.Int64  `tfsdk:"id"`
	ConfigurationID       types.Int64  `tfsdk:"configuration_id"`
	Policy                types.String `tfsdk:"policy"`
	TargetType            types.String `tfsdk:"target_type"`
	TargetIdentifier      types.String `tfsdk:"target_identifier"`
	Description           types.String `tfsdk:"description"`
	CustomMessage         types.String `tfsdk:"custom_message"`
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
	default:
		panic("Unknown Santa rule policy")
	}

	var rulesetID types.Int64
	if sr.RulesetID != nil {
		rulesetID = types.Int64Value(int64(*sr.RulesetID))
	} else {
		rulesetID = types.Int64Null()
	}

	primaryUsers := make([]attr.Value, 0)
	for _, primaryUser := range sr.PrimaryUsers {
		primaryUsers = append(primaryUsers, types.StringValue(primaryUser))
	}

	excludedPrimaryUsers := make([]attr.Value, 0)
	for _, primaryUser := range sr.ExcludedPrimaryUsers {
		excludedPrimaryUsers = append(excludedPrimaryUsers, types.StringValue(primaryUser))
	}

	serialNumbers := make([]attr.Value, 0)
	for _, serialNumber := range sr.SerialNumbers {
		serialNumbers = append(serialNumbers, types.StringValue(serialNumber))
	}

	excludedSerialNumbers := make([]attr.Value, 0)
	for _, serialNumber := range sr.ExcludedSerialNumbers {
		excludedSerialNumbers = append(excludedSerialNumbers, types.StringValue(serialNumber))
	}

	tagIDs := make([]attr.Value, 0)
	for _, tagID := range sr.TagIDs {
		tagIDs = append(tagIDs, types.Int64Value(int64(tagID)))
	}

	excludedTagIDs := make([]attr.Value, 0)
	for _, tagID := range sr.ExcludedTagIDs {
		excludedTagIDs = append(excludedTagIDs, types.Int64Value(int64(tagID)))
	}

	return santaRule{
		ID:                    types.Int64Value(int64(sr.ID)),
		ConfigurationID:       types.Int64Value(int64(sr.ConfigurationID)),
		Policy:                types.StringValue(policy),
		TargetType:            types.StringValue(sr.TargetType),
		TargetIdentifier:      types.StringValue(sr.TargetIdentifier),
		Description:           types.StringValue(sr.Description),
		CustomMessage:         types.StringValue(sr.CustomMessage),
		RulesetID:             rulesetID,
		PrimaryUsers:          types.SetValueMust(types.StringType, primaryUsers),
		ExcludedPrimaryUsers:  types.SetValueMust(types.StringType, excludedPrimaryUsers),
		SerialNumbers:         types.SetValueMust(types.StringType, serialNumbers),
		ExcludedSerialNumbers: types.SetValueMust(types.StringType, excludedSerialNumbers),
		TagIDs:                types.SetValueMust(types.Int64Type, tagIDs),
		ExcludedTagIDs:        types.SetValueMust(types.Int64Type, excludedTagIDs),
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
	default:
		panic("Unknown Santa rule policy")
	}

	primaryUsers := make([]string, 0)
	for _, primaryUser := range data.PrimaryUsers.Elements() { // nil if null or unknown → no iterations
		primaryUsers = append(primaryUsers, primaryUser.(types.String).ValueString())
	}

	excludedPrimaryUsers := make([]string, 0)
	for _, primaryUser := range data.ExcludedPrimaryUsers.Elements() { // nil if null or unknown → no iterations
		excludedPrimaryUsers = append(excludedPrimaryUsers, primaryUser.(types.String).ValueString())
	}

	serialNumbers := make([]string, 0)
	for _, serialNumber := range data.SerialNumbers.Elements() { // nil if null or unknown → no iterations
		serialNumbers = append(serialNumbers, serialNumber.(types.String).ValueString())
	}

	excludedSerialNumbers := make([]string, 0)
	for _, serialNumber := range data.ExcludedSerialNumbers.Elements() { // nil if null or unknown → no iterations
		excludedSerialNumbers = append(excludedSerialNumbers, serialNumber.(types.String).ValueString())
	}

	tagIDs := make([]int, 0)
	for _, tagID := range data.TagIDs.Elements() { // nil if null or unknown → no iterations
		tagIDs = append(tagIDs, int(tagID.(types.Int64).ValueInt64()))
	}

	excludedTagIDs := make([]int, 0)
	for _, tagID := range data.ExcludedTagIDs.Elements() { // nil if null or unknown → no iterations
		excludedTagIDs = append(excludedTagIDs, int(tagID.(types.Int64).ValueInt64()))
	}

	return &goztl.SantaRuleRequest{
		ConfigurationID:       int(data.ConfigurationID.ValueInt64()),
		Policy:                policy,
		TargetType:            data.TargetType.ValueString(),
		TargetIdentifier:      data.TargetIdentifier.ValueString(),
		Description:           data.Description.ValueString(),
		CustomMessage:         data.CustomMessage.ValueString(),
		PrimaryUsers:          primaryUsers,
		ExcludedPrimaryUsers:  excludedPrimaryUsers,
		SerialNumbers:         serialNumbers,
		ExcludedSerialNumbers: excludedSerialNumbers,
		TagIDs:                tagIDs,
		ExcludedTagIDs:        excludedTagIDs,
	}
}
