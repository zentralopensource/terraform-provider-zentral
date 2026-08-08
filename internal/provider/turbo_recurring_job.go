package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type turboRecurringJob struct {
	ID                    types.String `tfsdk:"id"`
	ConfigurationID       types.String `tfsdk:"configuration_id"`
	JobID                 types.String `tfsdk:"job_id"`
	Interval              types.Int64  `tfsdk:"interval"`
	TagIDs                types.Set    `tfsdk:"tag_ids"`
	ExcludedTagIDs        types.Set    `tfsdk:"excluded_tag_ids"`
	SerialNumbers         types.Set    `tfsdk:"serial_numbers"`
	ExcludedSerialNumbers types.Set    `tfsdk:"excluded_serial_numbers"`
}

func turboRecurringJobForState(trj *goztl.TurboRecurringJob) turboRecurringJob {
	return turboRecurringJob{
		ID:                    types.StringValue(trj.ID),
		ConfigurationID:       types.StringValue(trj.ConfigurationID),
		JobID:                 types.StringValue(trj.JobID),
		Interval:              optionalInt64ForState(trj.Interval),
		TagIDs:                int64SetForState(trj.TagIDs),
		ExcludedTagIDs:        int64SetForState(trj.ExcludedTagIDs),
		SerialNumbers:         stringSetForState(trj.SerialNumbers),
		ExcludedSerialNumbers: stringSetForState(trj.ExcludedSerialNumbers),
	}
}

func turboRecurringJobRequestWithState(data turboRecurringJob) *goztl.TurboRecurringJobRequest {
	return &goztl.TurboRecurringJobRequest{
		ConfigurationID:       data.ConfigurationID.ValueString(),
		JobID:                 data.JobID.ValueString(),
		Interval:              optionalIntWithState(data.Interval),
		TagIDs:                intListWithState(data.TagIDs),
		ExcludedTagIDs:        intListWithState(data.ExcludedTagIDs),
		SerialNumbers:         stringListWithStateSet(data.SerialNumbers),
		ExcludedSerialNumbers: stringListWithStateSet(data.ExcludedSerialNumbers),
	}
}
