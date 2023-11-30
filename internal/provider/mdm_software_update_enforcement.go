package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmSoftwareUpdateEnforcement struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	DetailsURL    types.String `tfsdk:"details_url"`
	TagIDs        types.Set    `tfsdk:"tag_ids"`
	OSVersion     types.String `tfsdk:"os_version"`
	BuildVersion  types.String `tfsdk:"build_version"`
	LocalDateTime types.String `tfsdk:"local_datetime"`
	MaxOSVersion  types.String `tfsdk:"max_os_version"`
	DelayDays     types.Int64  `tfsdk:"delay_days"`
	LocalTime     types.String `tfsdk:"local_time"`
}

func mdmSoftwareUpdateEnforcementForState(msue *goztl.MDMSoftwareUpdateEnforcement) mdmSoftwareUpdateEnforcement {
	tagIDs := make([]attr.Value, 0)
	for _, tagID := range msue.TagIDs {
		tagIDs = append(tagIDs, types.Int64Value(int64(tagID)))
	}

	var ldt types.String
	if msue.LocalDateTime != nil {
		ldt = types.StringValue(*msue.LocalDateTime)
	} else {
		ldt = types.StringNull()
	}

	var dd types.Int64
	if msue.DelayDays != nil {
		dd = types.Int64Value(int64(*msue.DelayDays))
	} else {
		dd = types.Int64Null()
	}

	var lt types.String
	if msue.LocalTime != nil {
		lt = types.StringValue(*msue.LocalTime)
	} else {
		lt = types.StringNull()
	}

	return mdmSoftwareUpdateEnforcement{
		ID:            types.Int64Value(int64(msue.ID)),
		Name:          types.StringValue(msue.Name),
		DetailsURL:    types.StringValue(msue.DetailsURL),
		TagIDs:        types.SetValueMust(types.Int64Type, tagIDs),
		OSVersion:     types.StringValue(msue.OSVersion),
		BuildVersion:  types.StringValue(msue.BuildVersion),
		LocalDateTime: ldt,
		MaxOSVersion:  types.StringValue(msue.MaxOSVersion),
		DelayDays:     dd,
		LocalTime:     lt,
	}
}

func mdmSoftwareUpdateEnforcementRequestWithState(data mdmSoftwareUpdateEnforcement) *goztl.MDMSoftwareUpdateEnforcementRequest {
	tagIDs := make([]int, 0)
	for _, tagID := range data.TagIDs.Elements() { // nil if null or unknown → no iterations
		tagIDs = append(tagIDs, int(tagID.(types.Int64).ValueInt64()))
	}

	var ldt *string
	if !data.LocalDateTime.IsNull() {
		ldt = goztl.String(data.LocalDateTime.ValueString())
	}

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
		TagIDs:        tagIDs,
		OSVersion:     data.OSVersion.ValueString(),
		BuildVersion:  data.BuildVersion.ValueString(),
		LocalDateTime: ldt,
		MaxOSVersion:  data.MaxOSVersion.ValueString(),
		DelayDays:     dd,
		LocalTime:     lt,
	}
}
