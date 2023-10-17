package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type osqueryQuery struct {
	ID                     types.Int64  `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	SQL                    types.String `tfsdk:"sql"`
	Platforms              types.Set    `tfsdk:"platforms"`
	MinOsqueryVersion      types.String `tfsdk:"minimum_osquery_version"`
	Description            types.String `tfsdk:"description"`
	Value                  types.String `tfsdk:"value"`
	Version                types.Int64  `tfsdk:"version"`
	ComplianceCheckEnabled types.Bool   `tfsdk:"compliance_check_enabled"`
	TagID                  types.Int64  `tfsdk:"tag_id"`
	Scheduling             types.Object `tfsdk:"scheduling"`
}

var schedulingAttrTypes = map[string]attr.Type{
	"can_be_denylisted":   types.BoolType,
	"interval":            types.Int64Type,
	"log_removed_actions": types.BoolType,
	"pack_id":             types.Int64Type,
	"shard":               types.Int64Type,
	"snapshot_mode":       types.BoolType,
}

func osqueryQueryForState(oq *goztl.OsqueryQuery) osqueryQuery {
	platforms := make([]attr.Value, 0)
	for _, platform := range oq.Platforms {
		platforms = append(platforms, types.StringValue(platform))
	}

	var minOsqueryVersion types.String
	if oq.MinOsqueryVersion != nil {
		minOsqueryVersion = types.StringValue(*oq.MinOsqueryVersion)
	} else {
		minOsqueryVersion = types.StringNull()
	}

	var tagID types.Int64
	if oq.TagID != nil {
		tagID = types.Int64Value(int64(*oq.TagID))
	} else {
		tagID = types.Int64Null()
	}

	var scheduling types.Object
	if oq.Scheduling != nil {
		var shard types.Int64
		if oq.Scheduling.Shard != nil {
			shard = types.Int64Value(int64(*oq.Scheduling.Shard))
		} else {
			shard = types.Int64Null()
		}
		scheduling = types.ObjectValueMust(
			schedulingAttrTypes,
			map[string]attr.Value{
				"can_be_denylisted":   types.BoolValue(oq.Scheduling.CanBeDenyListed),
				"log_removed_actions": types.BoolValue(oq.Scheduling.LogRemovedActions),
				"interval":            types.Int64Value(int64(oq.Scheduling.Interval)),
				"pack_id":             types.Int64Value(int64(oq.Scheduling.PackID)),
				"shard":               shard,
				"snapshot_mode":       types.BoolValue(oq.Scheduling.SnapshotMode),
			},
		)
	} else {
		scheduling = types.ObjectNull(schedulingAttrTypes)
	}

	return osqueryQuery{
		ID:                     types.Int64Value(int64(oq.ID)),
		Name:                   types.StringValue(oq.Name),
		SQL:                    types.StringValue(oq.SQL),
		Platforms:              types.SetValueMust(types.StringType, platforms),
		MinOsqueryVersion:      minOsqueryVersion,
		Description:            types.StringValue(oq.Description),
		Value:                  types.StringValue(oq.Value),
		Version:                types.Int64Value(int64(oq.Version)),
		ComplianceCheckEnabled: types.BoolValue(oq.ComplianceCheckEnabled),
		TagID:                  tagID,
		Scheduling:             scheduling,
	}
}

func osqueryQueryRequestWithState(data osqueryQuery) *goztl.OsqueryQueryRequest {
	platforms := make([]string, 0)
	for _, platform := range data.Platforms.Elements() { // nil if null or unknown → no iterations
		platforms = append(platforms, platform.(types.String).ValueString())
	}

	var minOsqueryVersion *string
	if !data.MinOsqueryVersion.IsNull() {
		minOsqueryVersion = goztl.String(data.MinOsqueryVersion.ValueString())
	}

	var tagID *int
	if !data.TagID.IsNull() {
		tagID = goztl.Int(int(data.TagID.ValueInt64()))
	}

	req := &goztl.OsqueryQueryRequest{
		Name:                   data.Name.ValueString(),
		SQL:                    data.SQL.ValueString(),
		Platforms:              platforms,
		MinOsqueryVersion:      minOsqueryVersion,
		Description:            data.Description.ValueString(),
		Value:                  data.Value.ValueString(),
		ComplianceCheckEnabled: data.ComplianceCheckEnabled.ValueBool(),
		TagID:                  tagID,
	}

	if !data.Scheduling.IsNull() {
		schedulingMap := data.Scheduling.Attributes()
		if schedulingMap != nil {
			var shard *int
			if !schedulingMap["shard"].(types.Int64).IsNull() {
				shard = goztl.Int(int(schedulingMap["shard"].(types.Int64).ValueInt64()))
			}
			schReq := &goztl.OsqueryQuerySchedulingRequest{
				CanBeDenyListed:   schedulingMap["can_be_denylisted"].(types.Bool).ValueBool(),
				LogRemovedActions: schedulingMap["log_removed_actions"].(types.Bool).ValueBool(),
				Interval:          int(schedulingMap["interval"].(types.Int64).ValueInt64()),
				PackID:            int(schedulingMap["pack_id"].(types.Int64).ValueInt64()),
				Shard:             shard,
				SnapshotMode:      schedulingMap["snapshot_mode"].(types.Bool).ValueBool(),
			}
			req.Scheduling = schReq
		}
	}

	return req
}
