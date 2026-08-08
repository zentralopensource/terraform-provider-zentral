package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type turboOneTimeJob struct {
	ID                    types.String `tfsdk:"id"`
	ConfigurationID       types.String `tfsdk:"configuration_id"`
	JobID                 types.String `tfsdk:"job_id"`
	NotBefore             types.String `tfsdk:"not_before"`
	NotAfter              types.String `tfsdk:"not_after"`
	TagIDs                types.Set    `tfsdk:"tag_ids"`
	ExcludedTagIDs        types.Set    `tfsdk:"excluded_tag_ids"`
	SerialNumbers         types.Set    `tfsdk:"serial_numbers"`
	ExcludedSerialNumbers types.Set    `tfsdk:"excluded_serial_numbers"`
}

func turboOneTimeJobForState(totj *goztl.TurboOneTimeJob) turboOneTimeJob {
	return turboOneTimeJob{
		ID:                    types.StringValue(totj.ID),
		ConfigurationID:       types.StringValue(totj.ConfigurationID),
		JobID:                 types.StringValue(totj.JobID),
		NotBefore:             optionalStringForState(totj.NotBefore),
		NotAfter:              optionalStringForState(totj.NotAfter),
		TagIDs:                int64SetForState(totj.TagIDs),
		ExcludedTagIDs:        int64SetForState(totj.ExcludedTagIDs),
		SerialNumbers:         stringSetForState(totj.SerialNumbers),
		ExcludedSerialNumbers: stringSetForState(totj.ExcludedSerialNumbers),
	}
}

func turboOneTimeJobRequestWithState(data turboOneTimeJob) *goztl.TurboOneTimeJobRequest {
	return &goztl.TurboOneTimeJobRequest{
		ConfigurationID:       data.ConfigurationID.ValueString(),
		JobID:                 data.JobID.ValueString(),
		NotBefore:             optionalStringWithState(data.NotBefore),
		NotAfter:              optionalStringWithState(data.NotAfter),
		TagIDs:                intListWithState(data.TagIDs),
		ExcludedTagIDs:        intListWithState(data.ExcludedTagIDs),
		SerialNumbers:         stringListWithStateSet(data.SerialNumbers),
		ExcludedSerialNumbers: stringListWithStateSet(data.ExcludedSerialNumbers),
	}
}
