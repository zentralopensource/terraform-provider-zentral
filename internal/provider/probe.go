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
		inventoryFilters = append(inventoryFilters, types.ObjectValueMust(
			probeInventoryFilterAttrTypes,
			map[string]attr.Value{
				"meta_business_unit_ids": int64SetForState(inventoryFilter.MetaBusinessUnitIDs),
				"tag_ids":                int64SetForState(inventoryFilter.TagIDs),
				"platforms":              stringSetForState(inventoryFilter.Platforms),
				"types":                  stringSetForState(inventoryFilter.Types),
			},
		))
	}

	metadataFilters := make([]attr.Value, 0)
	for _, metadataFilter := range p.MetadataFilters {
		metadataFilters = append(metadataFilters, types.ObjectValueMust(
			probeMetadataFilterAttrTypes,
			map[string]attr.Value{
				"event_types":        stringSetForState(metadataFilter.EventTypes),
				"event_tags":         stringSetForState(metadataFilter.EventTags),
				"event_routing_keys": stringSetForState(metadataFilter.EventRoutingKeys),
			},
		))
	}

	payloadFilters := make([]attr.Value, 0)
	for _, payloadFilterItems := range p.PayloadFilters {
		pfis := make([]attr.Value, 0)
		for _, payloadFilterItem := range payloadFilterItems {
			pfis = append(pfis, types.ObjectValueMust(
				probePayloadFilterItemAttrTypes,
				map[string]attr.Value{
					"attribute": types.StringValue(payloadFilterItem.Attribute),
					"operator":  types.StringValue(payloadFilterItem.Operator),
					"values":    stringSetForState(payloadFilterItem.Values),
				},
			))
		}
		payloadFilters = append(payloadFilters, types.SetValueMust(types.ObjectType{AttrTypes: probePayloadFilterItemAttrTypes}, pfis))
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
		ActionIDs:        stringSetForState(p.ActionIDs),
		IncidentSeverity: optionalInt64ForState(p.IncidentSeverity),
	}
}

func probeRequestWithState(data probe) *goztl.ProbeRequest {
	inventoryFilters := make([]goztl.InventoryFilter, 0)
	for _, inventoryFilter := range data.InventoryFilters.Elements() {
		inventoryFilterAttrs := inventoryFilter.(types.Object).Attributes()
		inventoryFilters = append(inventoryFilters, goztl.InventoryFilter{
			MetaBusinessUnitIDs: intListWithState(inventoryFilterAttrs["meta_business_unit_ids"].(types.Set)),
			TagIDs:              intListWithState(inventoryFilterAttrs["tag_ids"].(types.Set)),
			Platforms:           stringListWithStateSet(inventoryFilterAttrs["platforms"].(types.Set)),
			Types:               stringListWithStateSet(inventoryFilterAttrs["types"].(types.Set)),
		})
	}

	metadataFilters := make([]goztl.MetadataFilter, 0)
	for _, metadataFilter := range data.MetadataFilters.Elements() {
		metadataFilterAttrs := metadataFilter.(types.Object).Attributes()
		metadataFilters = append(metadataFilters, goztl.MetadataFilter{
			EventTypes:       stringListWithStateSet(metadataFilterAttrs["event_types"].(types.Set)),
			EventTags:        stringListWithStateSet(metadataFilterAttrs["event_tags"].(types.Set)),
			EventRoutingKeys: stringListWithStateSet(metadataFilterAttrs["event_routing_keys"].(types.Set)),
		})
	}

	payloadFilters := make([][]goztl.PayloadFilterItem, 0)
	for _, payloadFilter := range data.PayloadFilters.Elements() {
		pfis := make([]goztl.PayloadFilterItem, 0)
		for _, payloadFilterItem := range payloadFilter.(types.Set).Elements() {
			payloadFilterItemAttrs := payloadFilterItem.(types.Object).Attributes()
			pfis = append(pfis, goztl.PayloadFilterItem{
				Attribute: payloadFilterItemAttrs["attribute"].(types.String).ValueString(),
				Operator:  payloadFilterItemAttrs["operator"].(types.String).ValueString(),
				Values:    stringListWithStateSet(payloadFilterItemAttrs["values"].(types.Set)),
			})
		}
		payloadFilters = append(payloadFilters, pfis)
	}

	req := &goztl.ProbeRequest{
		Name:             data.Name.ValueString(),
		Description:      data.Description.ValueString(),
		InventoryFilters: inventoryFilters,
		MetadataFilters:  metadataFilters,
		PayloadFilters:   payloadFilters,
		Active:           data.Active.ValueBool(),
		ActionIDs:        stringListWithStateSet(data.ActionIDs),
		IncidentSeverity: optionalIntWithState(data.IncidentSeverity),
	}

	return req
}
