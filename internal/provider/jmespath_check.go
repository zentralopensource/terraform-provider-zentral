package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type jmespathCheck struct {
	ID                 types.Int64  `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	SourceName         types.String `tfsdk:"source_name"`
	Platforms          types.Set    `tfsdk:"platforms"`
	TagIDs             types.Set    `tfsdk:"tag_ids"`
	JMESPathExpression types.String `tfsdk:"jmespath_expression"`
	Version            types.Int64  `tfsdk:"version"`
}

func jmespathCheckForState(j *goztl.JMESPathCheck) jmespathCheck {
	platforms := make([]attr.Value, 0)
	for _, pv := range j.Platforms {
		platforms = append(platforms, types.StringValue(pv))
	}

	tagIDs := make([]attr.Value, 0)
	for _, tv := range j.TagIDs {
		tagIDs = append(tagIDs, types.Int64Value(int64(tv)))
	}

	return jmespathCheck{
		ID:                 types.Int64Value(int64(j.ID)),
		Name:               types.StringValue(j.Name),
		Description:        types.StringValue(j.Description),
		SourceName:         types.StringValue(j.SourceName),
		Platforms:          types.SetValueMust(types.StringType, platforms),
		TagIDs:             types.SetValueMust(types.Int64Type, tagIDs),
		JMESPathExpression: types.StringValue(j.JMESPathExpression),
		Version:            types.Int64Value(int64(j.Version)),
	}
}
