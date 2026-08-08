package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmSoftwareUpdateEnforcement struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	DetailsURL    types.String `tfsdk:"details_url"`
	Platforms     types.Set    `tfsdk:"platforms"`
	TagIDs        types.Set    `tfsdk:"tag_ids"`
	OSVersion     types.String `tfsdk:"os_version"`
	BuildVersion  types.String `tfsdk:"build_version"`
	LocalDateTime types.String `tfsdk:"local_datetime"`
	MaxOSVersion  types.String `tfsdk:"max_os_version"`
	DelayDays     types.Int64  `tfsdk:"delay_days"`
	LocalTime     types.String `tfsdk:"local_time"`
}

func mdmSoftwareUpdateEnforcementForState(msue *goztl.MDMSoftwareUpdateEnforcement) mdmSoftwareUpdateEnforcement {
	return mdmSoftwareUpdateEnforcement{
		ID:            types.Int64Value(int64(msue.ID)),
		Name:          types.StringValue(msue.Name),
		DetailsURL:    types.StringValue(msue.DetailsURL),
		Platforms:     stringSetForState(msue.Platforms),
		TagIDs:        int64SetForState(msue.TagIDs),
		OSVersion:     types.StringValue(msue.OSVersion),
		BuildVersion:  types.StringValue(msue.BuildVersion),
		LocalDateTime: optionalStringForState(msue.LocalDateTime),
		MaxOSVersion:  types.StringValue(msue.MaxOSVersion),
		DelayDays:     optionalInt64ForState(msue.DelayDays),
		LocalTime:     optionalStringForState(msue.LocalTime),
	}
}

func mdmSoftwareUpdateEnforcementRequestWithState(data mdmSoftwareUpdateEnforcement) *goztl.MDMSoftwareUpdateEnforcementRequest {
	var dd *int
	if !data.DelayDays.IsUnknown() {
		if !data.DelayDays.IsNull() {
			dd = goztl.Int(int(data.DelayDays.ValueInt64()))
		}
	} else if data.MaxOSVersion.ValueString() != "" {
		dd = goztl.Int(14)
	}

	var lt *string
	if !data.LocalTime.IsUnknown() {
		if !data.LocalTime.IsNull() {
			lt = goztl.String(data.LocalTime.ValueString())
		}
	} else if data.MaxOSVersion.ValueString() != "" {
		lt = goztl.String("09:30:00")
	}

	return &goztl.MDMSoftwareUpdateEnforcementRequest{
		Name:          data.Name.ValueString(),
		DetailsURL:    data.DetailsURL.ValueString(),
		Platforms:     stringListWithStateSet(data.Platforms),
		TagIDs:        intListWithState(data.TagIDs),
		OSVersion:     data.OSVersion.ValueString(),
		BuildVersion:  data.BuildVersion.ValueString(),
		LocalDateTime: optionalStringWithState(data.LocalDateTime),
		MaxOSVersion:  data.MaxOSVersion.ValueString(),
		DelayDays:     dd,
		LocalTime:     lt,
	}
}
