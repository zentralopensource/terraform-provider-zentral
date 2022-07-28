package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func resourceImportStatePassthroughZentralID(ctx context.Context, name string, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	ztlID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("Zentral %s ID must be an integer", name),
		)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.Int64{Value: ztlID})...)
	}
}
