package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

var probeInventoryFilterAttrTypes = map[string]attr.Type{
	"meta_business_unit_ids": types.SetType{ElemType: types.Int64Type},
	"tag_ids":                types.SetType{ElemType: types.Int64Type},
	"platforms":              types.SetType{ElemType: types.StringType},
	"types":                  types.SetType{ElemType: types.StringType},
}

var probeMetadataFilterAttrTypes = map[string]attr.Type{
	"event_types":        types.SetType{ElemType: types.StringType},
	"event_tags":         types.SetType{ElemType: types.StringType},
	"event_routing_keys": types.SetType{ElemType: types.StringType},
}

var probePayloadFilterItemAttrTypes = map[string]attr.Type{
	"attribute": types.StringType,
	"operator":  types.StringType,
	"values":    types.SetType{ElemType: types.StringType},
}

type probe struct {
	ID               types.Int64  `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Slug             types.String `tfsdk:"slug"`
	Description      types.String `tfsdk:"description"`
	InventoryFilters types.Set    `tfsdk:"inventory_filters"`
	MetadataFilters  types.Set    `tfsdk:"metadata_filters"`
	PayloadFilters   types.Set    `tfsdk:"payload_filters"`
	Active           types.Bool   `tfsdk:"active"`
	ActionIDs        types.Set    `tfsdk:"action_ids"`
	IncidentSeverity types.Int64  `tfsdk:"incident_severity"`
}

func probeForState(p *goztl.Probe) probe {
	inventoryFilters := make([]attr.Value, 0)
	for _, inventoryFilter := range p.InventoryFilters {
		mbuIDs := make([]attr.Value, 0)
		for _, mbuID := range inventoryFilter.MetaBusinessUnitIDs {
			mbuIDs = append(mbuIDs, types.Int64Value(int64(mbuID)))
		}
		tagIDs := make([]attr.Value, 0)
		for _, tagID := range inventoryFilter.TagIDs {
			tagIDs = append(tagIDs, types.Int64Value(int64(tagID)))
		}
		platforms := make([]attr.Value, 0)
		for _, platform := range inventoryFilter.Platforms {
			platforms = append(platforms, types.StringValue(platform))
		}
		ts := make([]attr.Value, 0)
		for _, t := range inventoryFilter.Types {
			ts = append(ts, types.StringValue(t))
		}
		inventoryFilters = append(inventoryFilters, types.ObjectValueMust(
			probeInventoryFilterAttrTypes,
			map[string]attr.Value{
				"meta_business_unit_ids": types.SetValueMust(types.Int64Type, mbuIDs),
				"tag_ids":                types.SetValueMust(types.Int64Type, tagIDs),
				"platforms":              types.SetValueMust(types.StringType, platforms),
				"types":                  types.SetValueMust(types.StringType, ts),
			},
		))
	}

	metadataFilters := make([]attr.Value, 0)
	for _, metadataFilter := range p.MetadataFilters {
		etys := make([]attr.Value, 0)
		for _, eventType := range metadataFilter.EventTypes {
			etys = append(etys, types.StringValue(eventType))
		}
		etas := make([]attr.Value, 0)
		for _, eventTag := range metadataFilter.EventTags {
			etas = append(etas, types.StringValue(eventTag))
		}
		erks := make([]attr.Value, 0)
		for _, eventRoutingKey := range metadataFilter.EventRoutingKeys {
			erks = append(erks, types.StringValue(eventRoutingKey))
		}
		metadataFilters = append(metadataFilters, types.ObjectValueMust(
			probeMetadataFilterAttrTypes,
			map[string]attr.Value{
				"event_types":        types.SetValueMust(types.StringType, etys),
				"event_tags":         types.SetValueMust(types.StringType, etas),
				"event_routing_keys": types.SetValueMust(types.StringType, erks),
			},
		))
	}

	payloadFilters := make([]attr.Value, 0)
	for _, payloadFilterItems := range p.PayloadFilters {
		pfis := make([]attr.Value, 0)
		for _, payloadFilterItem := range payloadFilterItems {
			vals := make([]attr.Value, 0)
			for _, val := range payloadFilterItem.Values {
				vals = append(vals, types.StringValue(val))
			}
			pfis = append(pfis, types.ObjectValueMust(
				probePayloadFilterItemAttrTypes,
				map[string]attr.Value{
					"attribute": types.StringValue(payloadFilterItem.Attribute),
					"operator":  types.StringValue(payloadFilterItem.Operator),
					"values":    types.SetValueMust(types.StringType, vals),
				},
			))
		}
		payloadFilters = append(payloadFilters, types.SetValueMust(types.ObjectType{AttrTypes: probePayloadFilterItemAttrTypes}, pfis))
	}

	actionIDs := make([]attr.Value, 0)
	for _, actionID := range p.ActionIDs {
		actionIDs = append(actionIDs, types.StringValue(actionID))
	}

	var incidentSeverity types.Int64
	if p.IncidentSeverity != nil {
		incidentSeverity = types.Int64Value(int64(*p.IncidentSeverity))
	} else {
		incidentSeverity = types.Int64Null()
	}

	return probe{
		ID:               types.Int64Value(int64(p.ID)),
		Name:             types.StringValue(p.Name),
		Slug:             types.StringValue(p.Slug),
		Description:      types.StringValue(p.Description),
		InventoryFilters: types.SetValueMust(types.ObjectType{AttrTypes: probeInventoryFilterAttrTypes}, inventoryFilters),
		MetadataFilters:  types.SetValueMust(types.ObjectType{AttrTypes: probeMetadataFilterAttrTypes}, metadataFilters),
		PayloadFilters:   types.SetValueMust(types.SetType{ElemType: types.ObjectType{AttrTypes: probePayloadFilterItemAttrTypes}}, payloadFilters),
		Active:           types.BoolValue(p.Active),
		ActionIDs:        types.SetValueMust(types.StringType, actionIDs),
		IncidentSeverity: incidentSeverity,
	}
}

func probeRequestWithState(data probe) *goztl.ProbeRequest {
	inventoryFilters := make([]goztl.InventoryFilter, 0)
	for _, inventoryFilter := range data.InventoryFilters.Elements() {
		inventoryFilterAttrs := inventoryFilter.(types.Object).Attributes()
		mbuIDs := make([]int, 0)
		for _, mbuID := range inventoryFilterAttrs["meta_business_unit_ids"].(types.Set).Elements() {
			mbuIDs = append(mbuIDs, int(mbuID.(types.Int64).ValueInt64()))
		}
		tagIDs := make([]int, 0)
		for _, tagID := range inventoryFilterAttrs["tag_ids"].(types.Set).Elements() {
			tagIDs = append(tagIDs, int(tagID.(types.Int64).ValueInt64()))
		}
		platforms := make([]string, 0)
		for _, platform := range inventoryFilterAttrs["platforms"].(types.Set).Elements() {
			platforms = append(platforms, platform.(types.String).ValueString())
		}
		ts := make([]string, 0)
		for _, t := range inventoryFilterAttrs["types"].(types.Set).Elements() {
			ts = append(ts, t.(types.String).ValueString())
		}
		inventoryFilters = append(inventoryFilters, goztl.InventoryFilter{
			MetaBusinessUnitIDs: mbuIDs,
			TagIDs:              tagIDs,
			Platforms:           platforms,
			Types:               ts,
		})
	}

	metadataFilters := make([]goztl.MetadataFilter, 0)
	for _, metadataFilter := range data.MetadataFilters.Elements() {
		metadataFilterAttrs := metadataFilter.(types.Object).Attributes()
		etys := make([]string, 0)
		for _, eventType := range metadataFilterAttrs["event_types"].(types.Set).Elements() {
			etys = append(etys, eventType.(types.String).ValueString())
		}
		etas := make([]string, 0)
		for _, eventTag := range metadataFilterAttrs["event_tags"].(types.Set).Elements() {
			etas = append(etas, eventTag.(types.String).ValueString())
		}
		erks := make([]string, 0)
		for _, eventRoutingKey := range metadataFilterAttrs["event_routing_keys"].(types.Set).Elements() {
			erks = append(erks, eventRoutingKey.(types.String).ValueString())
		}
		metadataFilters = append(metadataFilters, goztl.MetadataFilter{
			EventTypes:       etys,
			EventTags:        etas,
			EventRoutingKeys: erks,
		})
	}

	payloadFilters := make([][]goztl.PayloadFilterItem, 0)
	for _, payloadFilter := range data.PayloadFilters.Elements() {
		pfis := make([]goztl.PayloadFilterItem, 0)
		for _, payloadFilterItem := range payloadFilter.(types.Set).Elements() {
			payloadFilterItemAttrs := payloadFilterItem.(types.Object).Attributes()
			vals := make([]string, 0)
			for _, val := range payloadFilterItemAttrs["values"].(types.Set).Elements() {
				vals = append(vals, val.(types.String).ValueString())
			}
			pfis = append(pfis, goztl.PayloadFilterItem{
				Attribute: payloadFilterItemAttrs["attribute"].(types.String).ValueString(),
				Operator:  payloadFilterItemAttrs["operator"].(types.String).ValueString(),
				Values:    vals,
			})
		}
		payloadFilters = append(payloadFilters, pfis)
	}

	actionIDs := make([]string, 0)
	for _, actionID := range data.ActionIDs.Elements() {
		actionIDs = append(actionIDs, actionID.(types.String).ValueString())
	}

	var incidentSeverity *int
	if !data.IncidentSeverity.IsNull() {
		incidentSeverity = goztl.Int(int(data.IncidentSeverity.ValueInt64()))
	}

	req := &goztl.ProbeRequest{
		Name:             data.Name.ValueString(),
		Description:      data.Description.ValueString(),
		InventoryFilters: inventoryFilters,
		MetadataFilters:  metadataFilters,
		PayloadFilters:   payloadFilters,
		Active:           data.Active.ValueBool(),
		ActionIDs:        actionIDs,
		IncidentSeverity: incidentSeverity,
	}

	return req
}
