package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	dqs := make([]attr.Value, 0)
	for _, dq := range op.DiscoveryQueries {
		dqs = append(dqs, types.StringValue(dq))
	}

	var shard types.Int64
	if op.Shard != nil {
		shard = types.Int64Value(int64(*op.Shard))
	} else {
		shard = types.Int64Null()
	}

	return osqueryPack{
		ID:               types.Int64Value(int64(op.ID)),
		Name:             types.StringValue(op.Name),
		Slug:             types.StringValue(op.Slug),
		Description:      types.StringValue(op.Description),
		DiscoveryQueries: types.ListValueMust(types.StringType, dqs),
		Shard:            shard,
		EventRoutingKey:  types.StringValue(op.EventRoutingKey),
	}
}

func osqueryPackRequestWithState(data osqueryPack) *goztl.OsqueryPackRequest {
	dqs := make([]string, 0)
	for _, dq := range data.DiscoveryQueries.Elements() { // nil if null or unknown → no iterations
		dqs = append(dqs, dq.(types.String).ValueString())
	}

	var shard *int
	if !data.Shard.IsNull() {
		shard = goztl.Int(int(data.Shard.ValueInt64()))
	}

	return &goztl.OsqueryPackRequest{
		Name:             data.Name.ValueString(),
		Description:      data.Description.ValueString(),
		DiscoveryQueries: dqs,
		Shard:            shard,
		EventRoutingKey:  data.EventRoutingKey.ValueString(),
	}
}
