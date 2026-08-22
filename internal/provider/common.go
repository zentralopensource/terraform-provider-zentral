package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

// Server support

// A Zentral that predates an attribute drops it from the requests, and leaves it out of the
// responses, where goztl decodes it as a zero value. For an attribute whose zero value means what
// the absent feature does, configurations that leave it alone keep working. Configurations that set
// it get their value silently dropped, and the zero echoed back differs from the plan, which
// Terraform reports as an inconsistent result — a provider bug, as far as the message goes. Naming
// the cause is the point of this check. State has to be set before it runs, or the resource the
// server did create is left out of the state.
func checkInt64AttributeSupport(diags *diag.Diagnostics, attribute string, minVersion string, planned types.Int64, echoed int) {
	if echoed != 0 || planned.IsNull() || planned.IsUnknown() || planned.ValueInt64() == 0 {
		return
	}
	diags.AddError(
		"Unsupported Zentral version",
		fmt.Sprintf(
			"The Zentral server ignored the %s attribute and returned 0 instead of %d. "+
				"This attribute requires Zentral %s or later.",
			attribute, planned.ValueInt64(), minVersion,
		),
	)
}

// A Zentral that predates the data asset source attribute drops it from the request, and then
// rejects the request for the file_uri and the file_sha256 it still requires. The 400 names two
// attributes the configuration does not set, which reads as a provider bug. Ask the endpoint what
// it accepts instead of guessing from the body of the error: DRF answers an OPTIONS request with
// the fields of the serializer, for the methods the token is allowed to use.
//
// Say nothing when the probe cannot answer. A server that has no OPTIONS on the API, a token
// without the view permission, and a network error all leave the support of the attribute unknown,
// and the error from the API is a better answer than a wrong diagnosis.
func addMDMDataAssetSourceSupportDiagnostic(
	ctx context.Context, client *goztl.Client, diags *diag.Diagnostics, err error, source types.String,
) {
	if source.IsNull() || source.IsUnknown() {
		return
	}
	var errorResponse *goztl.ErrorResponse
	if !errors.As(err, &errorResponse) || errorResponse.Response == nil {
		return
	}
	if status := errorResponse.Response.StatusCode; status < 400 || status > 499 {
		return
	}
	endpointOptions, _, optionsErr := client.MDMDataAssets.Options(ctx)
	if optionsErr != nil {
		return
	}
	supported, known := endpointOptions.SupportsField(http.MethodPost, "source")
	if !known || supported {
		return
	}
	diags.AddError(
		"Unsupported Zentral version",
		fmt.Sprintf(
			"The Zentral server does not accept the source attribute of a MDM data asset. "+
				"This attribute requires Zentral %s or later.",
			minZentralVersionDataAssetSource,
		),
	)
}

// Timestamps

// Zentral runs with USE_TZ = False. The API parses any ISO 8601 timestamp, but it stores and echoes
// back a naive UTC datetime with no timezone suffix and no fractional seconds. Only that form comes
// back unchanged, and an echo that differs from the configured value fails the apply, so it is the
// only form the timestamp attributes accept.
var naiveTimestampRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}$`)

var naiveTimestampValidator = stringvalidator.RegexMatches(
	naiveTimestampRe,
	"must be a UTC timestamp without a timezone suffix, e.g. 2026-09-01T09:00:00",
)

// Bool (optional)

func optionalBoolForState(b *bool) types.Bool {
	var bfs types.Bool
	if b != nil {
		bfs = types.BoolValue(*b)
	} else {
		bfs = types.BoolNull()
	}
	return bfs
}

func optionalBoolWithState(b types.Bool) *bool {
	var bws *bool
	if !b.IsNull() {
		bws = goztl.Bool(b.ValueBool())
	}
	return bws
}

// Int64 (optional)

func optionalInt64ForState(i *int) types.Int64 {
	var ifs types.Int64
	if i != nil {
		ifs = types.Int64Value(int64(*i))
	} else {
		ifs = types.Int64Null()
	}
	return ifs
}

func optionalIntWithState(i types.Int64) *int {
	var iws *int
	if !i.IsNull() {
		iws = goztl.Int(int(i.ValueInt64()))
	}
	return iws
}

// String (optional)

func optionalStringForState(s *string) types.String {
	var sfs types.String
	if s != nil {
		sfs = types.StringValue(*s)
	} else {
		sfs = types.StringNull()
	}
	return sfs
}

func optionalStringWithState(s types.String) *string {
	var sws *string
	if !s.IsNull() {
		sws = goztl.String(s.ValueString())
	}
	return sws
}

// List of strings

func stringListForState(lstr []string) types.List {
	sstr := make([]attr.Value, 0)
	for _, s := range lstr {
		sstr = append(sstr, types.StringValue(s))
	}
	return types.ListValueMust(types.StringType, sstr)
}

func stringListWithStateList(sstr types.List) []string {
	lstr := make([]string, 0)
	for _, s := range sstr.Elements() {
		lstr = append(lstr, s.(types.String).ValueString())
	}
	return lstr
}

func stringListWithStateSet(sstr types.Set) []string {
	lstr := make([]string, 0)
	for _, s := range sstr.Elements() {
		lstr = append(lstr, s.(types.String).ValueString())
	}
	return lstr
}

// Set of integers

func int64SetForState(li []int) types.Set {
	si := make([]attr.Value, 0)
	for _, i := range li {
		si = append(si, types.Int64Value(int64(i)))
	}
	return types.SetValueMust(types.Int64Type, si)
}

func intListWithState(si types.Set) []int {
	li := make([]int, 0)
	for _, i := range si.Elements() {
		li = append(li, int(i.(types.Int64).ValueInt64()))
	}
	return li
}

// Set of strings

func stringSetForState(lstr []string) types.Set {
	sstr := make([]attr.Value, 0)
	for _, s := range lstr {
		sstr = append(sstr, types.StringValue(s))
	}
	return types.SetValueMust(types.StringType, sstr)
}

func nullableStringSetForState(lstr []string) types.Set {
	if len(lstr) == 0 {
		return types.SetNull(types.StringType)
	} else {
		return stringSetForState(lstr)
	}
}

// EventFilterSet

var eventFilterAttrTypes = map[string]attr.Type{
	"tags":        types.SetType{ElemType: types.StringType},
	"event_type":  types.SetType{ElemType: types.StringType},
	"routing_key": types.SetType{ElemType: types.StringType},
}

var eventFilterSetAttrTypes = map[string]attr.Type{
	"excluded_event_filters": types.SetType{ElemType: types.ObjectType{AttrTypes: eventFilterAttrTypes}},
	"included_event_filters": types.SetType{ElemType: types.ObjectType{AttrTypes: eventFilterAttrTypes}},
}

func eventFiltersForState(zfs []goztl.EventFilter) types.Set {
	effs := make([]attr.Value, 0)
	for _, zef := range zfs {
		effs = append(effs, types.ObjectValueMust(
			eventFilterAttrTypes,
			map[string]attr.Value{
				"tags":        nullableStringSetForState(zef.Tags),
				"event_type":  nullableStringSetForState(zef.EventType),
				"routing_key": nullableStringSetForState(zef.RoutingKey),
			},
		))
	}
	if len(effs) > 0 {
		return types.SetValueMust(types.ObjectType{AttrTypes: eventFilterAttrTypes}, effs)
	} else {
		// empty filter list converted to a Null Set for the state
		return types.SetNull(types.ObjectType{AttrTypes: eventFilterAttrTypes})
	}
}

func eventFilterSetForState(zefs *goztl.EventFilterSet) types.Object {
	if zefs != nil {
		return types.ObjectValueMust(
			eventFilterSetAttrTypes,
			map[string]attr.Value{
				"excluded_event_filters": eventFiltersForState(zefs.ExcludedEventFilters),
				"included_event_filters": eventFiltersForState(zefs.IncludedEventFilters),
			},
		)
	} else {
		return types.ObjectNull(eventFilterSetAttrTypes)
	}
}

func eventFiltersWithState(efs types.Set) []goztl.EventFilter {
	efws := make([]goztl.EventFilter, 0)
	for _, mapF := range efs.Elements() {
		mapFAttrs := mapF.(types.Object).Attributes()
		efws = append(efws, goztl.EventFilter{
			Tags:       stringListWithStateSet(mapFAttrs["tags"].(types.Set)),
			EventType:  stringListWithStateSet(mapFAttrs["event_type"].(types.Set)),
			RoutingKey: stringListWithStateSet(mapFAttrs["routing_key"].(types.Set)),
		})
	}
	return efws
}

func eventFilterSetWithState(tfs types.Object) *goztl.EventFilterSet {
	var efs *goztl.EventFilterSet
	if !tfs.IsNull() {
		mapFSAttrs := tfs.Attributes()
		efs = &goztl.EventFilterSet{
			IncludedEventFilters: eventFiltersWithState(mapFSAttrs["included_event_filters"].(types.Set)),
			ExcludedEventFilters: eventFiltersWithState(mapFSAttrs["excluded_event_filters"].(types.Set)),
		}
	}
	return efs
}

func defaultEventFilterSet() types.Object {
	return types.ObjectValueMust(
		eventFilterSetAttrTypes,
		map[string]attr.Value{
			// empty filters are Null Sets for the state
			"included_event_filters": types.SetNull(types.ObjectType{AttrTypes: eventFilterAttrTypes}),
			"excluded_event_filters": types.SetNull(types.ObjectType{AttrTypes: eventFilterAttrTypes}),
		},
	)
}

// HTTPHeader

var headerAttrTypes = map[string]attr.Type{
	"name":  types.StringType,
	"value": types.StringType,
}

func headersForState(zhs []goztl.HTTPHeader) types.Set {
	hfs := make([]attr.Value, 0)
	for _, zh := range zhs {
		hfs = append(hfs, types.ObjectValueMust(
			headerAttrTypes,
			map[string]attr.Value{
				"name":  types.StringValue(zh.Name),
				"value": types.StringValue(zh.Value),
			},
		))
	}
	return types.SetValueMust(types.ObjectType{AttrTypes: headerAttrTypes}, hfs)
}

func headersWithState(ths types.Set) []goztl.HTTPHeader {
	hws := make([]goztl.HTTPHeader, 0)
	for _, mapH := range ths.Elements() {
		mapHAttrs := mapH.(types.Object).Attributes()
		hws = append(hws, goztl.HTTPHeader{
			Name:  mapHAttrs["name"].(types.String).ValueString(),
			Value: mapHAttrs["value"].(types.String).ValueString(),
		})
	}
	return hws
}
