package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type osqueryPack struct {
	ID               types.Int64  `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Slug             types.String `tfsdk:"slug"`
	Description      types.String `tfsdk:"description"`
	DiscoveryQueries types.List   `tfsdk:"discovery_queries"`
	Shard            types.Int64  `tfsdk:"shard"`
	EventRoutingKey  types.String `tfsdk:"event_routing_key"`
}

func osqueryPackForState(op *goztl.OsqueryPack) osqueryPack {
	return osqueryPack{
		ID:               types.Int64Value(int64(op.ID)),
		Name:             types.StringValue(op.Name),
		Slug:             types.StringValue(op.Slug),
		Description:      types.StringValue(op.Description),
		DiscoveryQueries: stringListForState(op.DiscoveryQueries),
		Shard:            optionalInt64ForState(op.Shard),
		EventRoutingKey:  types.StringValue(op.EventRoutingKey),
	}
}

func osqueryPackRequestWithState(data osqueryPack) *goztl.OsqueryPackRequest {
	return &goztl.OsqueryPackRequest{
		Name:             data.Name.ValueString(),
		Description:      data.Description.ValueString(),
		DiscoveryQueries: stringListWithStateList(data.DiscoveryQueries),
		Shard:            optionalIntWithState(data.Shard),
		EventRoutingKey:  data.EventRoutingKey.ValueString(),
	}
}
