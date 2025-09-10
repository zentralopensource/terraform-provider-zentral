package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

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
