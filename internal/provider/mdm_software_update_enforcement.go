package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const minZentralVersionSUEDeviceOSVersion = "v2026.6"

type mdmSoftwareUpdateEnforcement struct {
	ID                 types.Int64  `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	DetailsURL         types.String `tfsdk:"details_url"`
	Platforms          types.Set    `tfsdk:"platforms"`
	TagIDs             types.Set    `tfsdk:"tag_ids"`
	MinDeviceOSVersion types.String `tfsdk:"min_device_os_version"`
	MaxDeviceOSVersion types.String `tfsdk:"max_device_os_version"`
	OSVersion          types.String `tfsdk:"os_version"`
	BuildVersion       types.String `tfsdk:"build_version"`
	LocalDateTime      types.String `tfsdk:"local_datetime"`
	MaxOSVersion       types.String `tfsdk:"max_os_version"`
	DelayDays          types.Int64  `tfsdk:"delay_days"`
	LocalTime          types.String `tfsdk:"local_time"`
}

func mdmSoftwareUpdateEnforcementForState(msue *goztl.MDMSoftwareUpdateEnforcement) mdmSoftwareUpdateEnforcement {
	return mdmSoftwareUpdateEnforcement{
		ID:                 types.Int64Value(int64(msue.ID)),
		Name:               types.StringValue(msue.Name),
		DetailsURL:         types.StringValue(msue.DetailsURL),
		Platforms:          stringSetForState(msue.Platforms),
		TagIDs:             int64SetForState(msue.TagIDs),
		MinDeviceOSVersion: types.StringValue(msue.MinDeviceOSVersion),
		MaxDeviceOSVersion: types.StringValue(msue.MaxDeviceOSVersion),
		OSVersion:          types.StringValue(msue.OSVersion),
		BuildVersion:       types.StringValue(msue.BuildVersion),
		LocalDateTime:      optionalStringForState(msue.LocalDateTime),
		MaxOSVersion:       types.StringValue(msue.MaxOSVersion),
		DelayDays:          optionalInt64ForState(msue.DelayDays),
		LocalTime:          optionalStringForState(msue.LocalTime),
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
		Name:               data.Name.ValueString(),
		DetailsURL:         data.DetailsURL.ValueString(),
		Platforms:          stringListWithStateSet(data.Platforms),
		TagIDs:             intListWithState(data.TagIDs),
		MinDeviceOSVersion: data.MinDeviceOSVersion.ValueString(),
		MaxDeviceOSVersion: data.MaxDeviceOSVersion.ValueString(),
		OSVersion:          data.OSVersion.ValueString(),
		BuildVersion:       data.BuildVersion.ValueString(),
		LocalDateTime:      optionalStringWithState(data.LocalDateTime),
		MaxOSVersion:       data.MaxOSVersion.ValueString(),
		DelayDays:          dd,
		LocalTime:          lt,
	}
}

func checkMDMSoftwareUpdateEnforcementDeviceOSVersionSupport(diags *diag.Diagnostics, planned mdmSoftwareUpdateEnforcement, echoed *goztl.MDMSoftwareUpdateEnforcement) {
	checkStringAttributeSupport(
		diags, "min_device_os_version", minZentralVersionSUEDeviceOSVersion,
		planned.MinDeviceOSVersion, echoed.MinDeviceOSVersion, "",
	)
	checkStringAttributeSupport(
		diags, "max_device_os_version", minZentralVersionSUEDeviceOSVersion,
		planned.MaxDeviceOSVersion, echoed.MaxDeviceOSVersion, "",
	)
}

// The rule of the Zentral API for a bound of the device OS version range: one to three integer
// components, not all zero, no supplemental extra. The components are padded to three, as Zentral
// compares them.
func parseDeviceOSVersionBound(bound string) ([3]int, error) {
	var parsed [3]int
	components := strings.Split(bound, ".")
	if len(components) > 3 {
		return parsed, errors.New("must be a version with at most three components, such as 26 or 26.3.1")
	}
	allZero := true
	for i, component := range components {
		if component == "" || strings.Trim(component, "0123456789") != "" {
			return parsed, errors.New("must be a version with at most three components, such as 26 or 26.3.1")
		}
		n, err := strconv.Atoi(component)
		if err != nil {
			return parsed, errors.New("must be a version with at most three components, such as 26 or 26.3.1")
		}
		parsed[i] = n
		if n != 0 {
			allZero = false
		}
	}
	if allZero {
		return parsed, errors.New("must not be 0")
	}
	return parsed, nil
}

type deviceOSVersionBoundValidator struct{}

var _ validator.String = deviceOSVersionBoundValidator{}

func (v deviceOSVersionBoundValidator) Description(ctx context.Context) string {
	return "a version with at most three components, such as 26 or 26.3.1, and not 0"
}

func (v deviceOSVersionBoundValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v deviceOSVersionBoundValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.ValueString() == "" {
		return
	}
	if _, err := parseDeviceOSVersionBound(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Attribute Value",
			fmt.Sprintf("Attribute %s %s, got: %q.", req.Path, err, req.ConfigValue.ValueString()),
		)
	}
}

// The API refuses a maximum that is not greater than the minimum with a 400 at apply time. A bound
// that does not parse is reported by its own validator, so it is left out here.
func validateDeviceOSVersionRange(diags *diag.Diagnostics, minBound, maxBound types.String) {
	isSet := func(s types.String) bool {
		return !s.IsNull() && !s.IsUnknown() && s.ValueString() != ""
	}
	if !isSet(minBound) || !isSet(maxBound) {
		return
	}
	parsedMin, err := parseDeviceOSVersionBound(minBound.ValueString())
	if err != nil {
		return
	}
	parsedMax, err := parseDeviceOSVersionBound(maxBound.ValueString())
	if err != nil {
		return
	}
	for i := range parsedMin {
		if parsedMax[i] > parsedMin[i] {
			return
		}
		if parsedMax[i] < parsedMin[i] {
			break
		}
	}
	diags.AddAttributeError(
		path.Root("max_device_os_version"),
		"Invalid Attribute Value",
		fmt.Sprintf(
			"max_device_os_version must be greater than min_device_os_version, got: %q and %q.",
			maxBound.ValueString(), minBound.ValueString(),
		),
	)
}
