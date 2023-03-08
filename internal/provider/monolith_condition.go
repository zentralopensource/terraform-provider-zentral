package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type monolithCondition struct {
	ID        types.Int64  `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Predicate types.String `tfsdk:"predicate"`
}

func monolithConditionForState(mc *goztl.MonolithCondition) monolithCondition {
	return monolithCondition{
		ID:        types.Int64Value(int64(mc.ID)),
		Name:      types.StringValue(mc.Name),
		Predicate: types.StringValue(mc.Predicate),
	}
}

func monolithConditionRequestWithState(data monolithCondition) *goztl.MonolithConditionRequest {
	return &goztl.MonolithConditionRequest{
		Name:      data.Name.ValueString(),
		Predicate: data.Predicate.ValueString(),
	}
}
