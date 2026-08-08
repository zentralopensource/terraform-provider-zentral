package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type turboMSCPCheck struct {
	ID                types.String `tfsdk:"id"`
	RuleID            types.String `tfsdk:"rule_id"`
	Baseline          types.String `tfsdk:"baseline"`
	ODVInt            types.Int64  `tfsdk:"odv_int"`
	ODVString         types.String `tfsdk:"odv_string"`
	ODVBool           types.Bool   `tfsdk:"odv_bool"`
	Version           types.Int64  `tfsdk:"version"`
	ComplianceCheckID types.Int64  `tfsdk:"compliance_check_id"`
	JobID             types.String `tfsdk:"job_id"`
}

func turboMSCPCheckForState(tmc *goztl.TurboMSCPCheck) turboMSCPCheck {
	return turboMSCPCheck{
		ID:                types.StringValue(tmc.ID),
		RuleID:            types.StringValue(tmc.RuleID),
		Baseline:          types.StringValue(tmc.Baseline),
		ODVInt:            optionalInt64ForState(tmc.ODVInt),
		ODVString:         optionalStringForState(tmc.ODVString),
		ODVBool:           optionalBoolForState(tmc.ODVBool),
		Version:           types.Int64Value(int64(tmc.Version)),
		ComplianceCheckID: types.Int64Value(int64(tmc.ComplianceCheckID)),
		JobID:             types.StringValue(tmc.JobID),
	}
}

func turboMSCPCheckRequestWithState(data turboMSCPCheck) *goztl.TurboMSCPCheckRequest {
	return &goztl.TurboMSCPCheckRequest{
		RuleID:    data.RuleID.ValueString(),
		Baseline:  data.Baseline.ValueString(),
		ODVInt:    optionalIntWithState(data.ODVInt),
		ODVString: optionalStringWithState(data.ODVString),
		ODVBool:   optionalBoolWithState(data.ODVBool),
	}
}
