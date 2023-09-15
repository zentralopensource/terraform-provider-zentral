package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

func exTagIDsForState(mav goztl.MDMArtifactVersion) []attr.Value {
	exTagIDs := make([]attr.Value, 0)
	for _, exTagID := range mav.ExcludedTagIDs {
		exTagIDs = append(exTagIDs, types.Int64Value(int64(exTagID)))
	}
	return exTagIDs
}

func tagShardsForState(mav goztl.MDMArtifactVersion) []attr.Value {
	tagShards := make([]attr.Value, 0)
	for _, tagShard := range mav.TagShards {
		tagShards = append(
			tagShards,
			types.ObjectValueMust(
				tagShardAttrTypes,
				map[string]attr.Value{
					"tag_id": types.Int64Value(int64(tagShard.TagID)),
					"shard":  types.Int64Value(int64(tagShard.Shard)),
				},
			),
		)
	}
	return tagShards
}
