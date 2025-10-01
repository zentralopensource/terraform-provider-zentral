package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

func tagShardsForState(mav goztl.MDMArtifactVersion) types.Set {
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
	return types.SetValueMust(types.ObjectType{AttrTypes: tagShardAttrTypes}, tagShards)
}

func tagShardsWithState(ts types.Set) []goztl.TagShard {
	tagShards := make([]goztl.TagShard, 0)
	for _, tagShard := range ts.Elements() { // nil if null or unknown → no iterations
		tagShardMap := tagShard.(types.Object).Attributes()
		if tagShardMap != nil {
			tagShards = append(
				tagShards,
				goztl.TagShard{
					TagID: int(tagShardMap["tag_id"].(types.Int64).ValueInt64()),
					Shard: int(tagShardMap["shard"].(types.Int64).ValueInt64()),
				},
			)
		}
	}
	return tagShards
}
