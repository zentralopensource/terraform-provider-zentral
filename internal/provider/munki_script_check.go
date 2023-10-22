package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type munkiScriptCheck struct {
	ID             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Type           types.String `tfsdk:"type"`
	Source         types.String `tfsdk:"source"`
	ExpectedResult types.String `tfsdk:"expected_result"`
	ArchAMD64      types.Bool   `tfsdk:"arch_amd64"`
	ArchARM64      types.Bool   `tfsdk:"arch_arm64"`
	MinOSVersion   types.String `tfsdk:"min_os_version"`
	MaxOSVersion   types.String `tfsdk:"max_os_version"`
	TagIDs         types.Set    `tfsdk:"tag_ids"`
	Version        types.Int64  `tfsdk:"version"`
}

func munkiScriptCheckForState(msc *goztl.MunkiScriptCheck) munkiScriptCheck {
	tagIDs := make([]attr.Value, 0)
	for _, tv := range msc.TagIDs {
		tagIDs = append(tagIDs, types.Int64Value(int64(tv)))
	}

	return munkiScriptCheck{
		ID:             types.Int64Value(int64(msc.ID)),
		Name:           types.StringValue(msc.Name),
		Description:    types.StringValue(msc.Description),
		Type:           types.StringValue(msc.Type),
		Source:         types.StringValue(msc.Source),
		ExpectedResult: types.StringValue(msc.ExpectedResult),
		ArchAMD64:      types.BoolValue(msc.ArchAMD64),
		ArchARM64:      types.BoolValue(msc.ArchARM64),
		MinOSVersion:   types.StringValue(msc.MinOSVersion),
		MaxOSVersion:   types.StringValue(msc.MaxOSVersion),
		TagIDs:         types.SetValueMust(types.Int64Type, tagIDs),
		Version:        types.Int64Value(int64(msc.Version)),
	}
}

func munkiScriptCheckRequestWithState(data munkiScriptCheck) *goztl.MunkiScriptCheckRequest {
	tagIDs := make([]int, 0)
	for _, tagID := range data.TagIDs.Elements() { // nil if null or unknown → no iterations
		tagIDs = append(tagIDs, int(tagID.(types.Int64).ValueInt64()))
	}

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
		TagIDs:         tagIDs,
	}
}
