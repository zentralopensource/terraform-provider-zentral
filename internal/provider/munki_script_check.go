package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type munkiScriptCheck struct {
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Type              types.String `tfsdk:"type"`
	Source            types.String `tfsdk:"source"`
	ExpectedResult    types.String `tfsdk:"expected_result"`
	ArchAMD64         types.Bool   `tfsdk:"arch_amd64"`
	ArchARM64         types.Bool   `tfsdk:"arch_arm64"`
	MinOSVersion      types.String `tfsdk:"min_os_version"`
	MaxOSVersion      types.String `tfsdk:"max_os_version"`
	TagIDs            types.Set    `tfsdk:"tag_ids"`
	ExcludedTagIDs    types.Set    `tfsdk:"excluded_tag_ids"`
	Version           types.Int64  `tfsdk:"version"`
	ComplianceCheckID types.Int64  `tfsdk:"compliance_check_id"`
}

func munkiScriptCheckForState(msc *goztl.MunkiScriptCheck) munkiScriptCheck {
	return munkiScriptCheck{
		ID:                types.Int64Value(int64(msc.ID)),
		Name:              types.StringValue(msc.Name),
		Description:       types.StringValue(msc.Description),
		Type:              types.StringValue(msc.Type),
		Source:            types.StringValue(msc.Source),
		ExpectedResult:    types.StringValue(msc.ExpectedResult),
		ArchAMD64:         types.BoolValue(msc.ArchAMD64),
		ArchARM64:         types.BoolValue(msc.ArchARM64),
		MinOSVersion:      types.StringValue(msc.MinOSVersion),
		MaxOSVersion:      types.StringValue(msc.MaxOSVersion),
		TagIDs:            int64SetForState(msc.TagIDs),
		ExcludedTagIDs:    int64SetForState(msc.ExcludedTagIDs),
		Version:           types.Int64Value(int64(msc.Version)),
		ComplianceCheckID: types.Int64Value(int64(msc.ComplianceCheckID)),
	}
}

func munkiScriptCheckRequestWithState(data munkiScriptCheck) *goztl.MunkiScriptCheckRequest {
	return &goztl.MunkiScriptCheckRequest{
		Name:           data.Name.ValueString(),
		Description:    data.Description.ValueString(),
		Type:           data.Type.ValueString(),
		Source:         data.Source.ValueString(),
		ExpectedResult: data.ExpectedResult.ValueString(),
		ArchAMD64:      data.ArchAMD64.ValueBool(),
		ArchARM64:      data.ArchARM64.ValueBool(),
		MinOSVersion:   data.MinOSVersion.ValueString(),
		MaxOSVersion:   data.MaxOSVersion.ValueString(),
		TagIDs:         intListWithState(data.TagIDs),
		ExcludedTagIDs: intListWithState(data.ExcludedTagIDs),
	}
}
