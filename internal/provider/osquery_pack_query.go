package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type osqueryPackQuery struct {
	ID                types.Int64  `tfsdk:"id"`
	PackID            types.Int64  `tfsdk:"pack_id"`
	QueryID           types.Int64  `tfsdk:"query_id"`
	Slug              types.String `tfsdk:"slug"`
	Interval          types.Int64  `tfsdk:"interval"`
	LogRemovedActions types.Bool   `tfsdk:"log_removed_actions"`
	SnapshotMode      types.Bool   `tfsdk:"snapshot_mode"`
	Shard             types.Int64  `tfsdk:"shard"`
	CanBeDenyListed   types.Bool   `tfsdk:"can_be_denylisted"`
}

func osqueryPackQueryForState(opq *goztl.OsqueryPackQuery) osqueryPackQuery {
	var shard types.Int64
	if opq.Shard != nil {
		shard = types.Int64Value(int64(*opq.Shard))
	} else {
		shard = types.Int64Null()
	}

	return osqueryPackQuery{
		ID:                types.Int64Value(int64(opq.ID)),
		PackID:            types.Int64Value(int64(opq.PackID)),
		QueryID:           types.Int64Value(int64(opq.QueryID)),
		Slug:              types.StringValue(opq.Slug),
		Interval:          types.Int64Value(int64(opq.Interval)),
		LogRemovedActions: types.BoolValue(opq.LogRemovedActions),
		SnapshotMode:      types.BoolValue(opq.SnapshotMode),
		Shard:             shard,
		CanBeDenyListed:   types.BoolValue(opq.CanBeDenyListed),
	}
}

func osqueryPackQueryRequestWithState(data osqueryPackQuery) *goztl.OsqueryPackQueryRequest {
	var shard *int
	if !data.Shard.IsNull() {
		shard = goztl.Int(int(data.Shard.ValueInt64()))
	}

	return &goztl.OsqueryPackQueryRequest{
		PackID:            int(data.PackID.ValueInt64()),
		QueryID:           int(data.QueryID.ValueInt64()),
		Interval:          int(data.Interval.ValueInt64()),
		LogRemovedActions: data.LogRemovedActions.ValueBool(),
		SnapshotMode:      data.SnapshotMode.ValueBool(),
		Shard:             shard,
		CanBeDenyListed:   data.CanBeDenyListed.ValueBool(),
	}
}
