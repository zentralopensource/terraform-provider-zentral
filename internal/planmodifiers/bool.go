package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func DefaultTrue() planmodifier.Bool {
	return defaultTrueModifier{}
}

// defaultTrueModifier is a plan modifier that sets a default True
// types.BoolType attribute when it is not configured. The attribute must be
// marked as Optional and Computed. When setting the state during the resource
// Create, Read, or Update methods, this default value must also be included or
// the Terraform CLI will generate an error.
type defaultTrueModifier struct{}

func (m defaultTrueModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to true")
}

func (m defaultTrueModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `true`")
}

// Set the plan value to true if missing in the configuration
func (m defaultTrueModifier) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.ConfigValue.IsNull() {
		resp.PlanValue = types.BoolValue(true)
	}
}
