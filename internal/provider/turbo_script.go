package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type turboScript struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Description            types.String `tfsdk:"description"`
	Source                 types.String `tfsdk:"source"`
	TagID                  types.Int64  `tfsdk:"tag_id"`
	ArchAMD64              types.Bool   `tfsdk:"arch_amd64"`
	ArchARM64              types.Bool   `tfsdk:"arch_arm64"`
	MinOSVersion           types.String `tfsdk:"min_os_version"`
	MaxOSVersion           types.String `tfsdk:"max_os_version"`
	ComplianceCheckEnabled types.Bool   `tfsdk:"compliance_check_enabled"`
	ComplianceCheckID      types.Int64  `tfsdk:"compliance_check_id"`
	Version                types.Int64  `tfsdk:"version"`
	JobID                  types.String `tfsdk:"job_id"`
}

func turboScriptForState(ts *goztl.TurboScript) turboScript {
	return turboScript{
		ID:                     types.StringValue(ts.ID),
		Name:                   types.StringValue(ts.Name),
		Description:            types.StringValue(ts.Description),
		Source:                 types.StringValue(ts.Source),
		TagID:                  optionalInt64ForState(ts.TagID),
		ArchAMD64:              types.BoolValue(ts.ArchAMD64),
		ArchARM64:              types.BoolValue(ts.ArchARM64),
		MinOSVersion:           types.StringValue(ts.MinOSVersion),
		MaxOSVersion:           types.StringValue(ts.MaxOSVersion),
		ComplianceCheckEnabled: types.BoolValue(ts.ComplianceCheckEnabled),
		ComplianceCheckID:      optionalInt64ForState(ts.ComplianceCheckID),
		Version:                types.Int64Value(int64(ts.Version)),
		JobID:                  types.StringValue(ts.JobID),
	}
}

func turboScriptRequestWithState(data turboScript) *goztl.TurboScriptRequest {
	return &goztl.TurboScriptRequest{
		Name:                   data.Name.ValueString(),
		Description:            data.Description.ValueString(),
		Source:                 data.Source.ValueString(),
		TagID:                  optionalIntWithState(data.TagID),
		ArchAMD64:              data.ArchAMD64.ValueBool(),
		ArchARM64:              data.ArchARM64.ValueBool(),
		MinOSVersion:           data.MinOSVersion.ValueString(),
		MaxOSVersion:           data.MaxOSVersion.ValueString(),
		ComplianceCheckEnabled: data.ComplianceCheckEnabled.ValueBool(),
	}
}
