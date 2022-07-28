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
		platforms = append(platforms, types.String{Value: pv})
	}

	tagIDs := make([]attr.Value, 0)
	for _, tv := range j.TagIDs {
		tagIDs = append(tagIDs, types.Int64{Value: int64(tv)})
	}

	return jmespathCheck{
		ID:                 types.Int64{Value: int64(j.ID)},
		Name:               types.String{Value: j.Name},
		Description:        types.String{Value: j.Description},
		SourceName:         types.String{Value: j.SourceName},
		Platforms:          types.Set{ElemType: types.StringType, Elems: platforms},
		TagIDs:             types.Set{ElemType: types.Int64Type, Elems: tagIDs},
		JMESPathExpression: types.String{Value: j.JMESPathExpression},
		Version:            types.Int64{Value: int64(j.Version)},
	}
}
